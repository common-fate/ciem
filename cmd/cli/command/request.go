package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	accessv1alpha1 "github.com/common-fate/ciem/gen/commonfate/access/v1alpha1"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/ciem/service/access"
	"github.com/common-fate/ciem/service/request"
	"github.com/common-fate/ciem/table"
	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var Request = cli.Command{
	Name:  "request",
	Usage: "Request access to an entitlement",
	Subcommands: []*cli.Command{
		&gcpRequest,
		&review,
		&revoke,
		&list,
		&cancel,
	},
}

var gcpRequest = cli.Command{
	Name:  "gcp",
	Usage: "Request access to a GCP entitlements",
	Subcommands: []*cli.Command{
		&gcpProjectRequest,
	},
}

var gcpProjectRequest = cli.Command{
	Name:  "project",
	Usage: "Request access to a GCP Project",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "role", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		res, err := client.Grant(ctx, connect.NewRequest(&accessv1alpha1.GrantRequest{
			Role: &authzv1alpha1.UID{
				Type: "GCP::Role",
				Id:   c.String("role"),
			},
			Target: &authzv1alpha1.UID{
				Type: "GCP::Project",
				Id:   c.String("id"),
			},
		}))
		if err != nil {
			return err
		}

		clio.Infow("response", "response", res)
		if res.Msg.Decision == accessv1alpha1.Decision_DECISION_DENIED {
			clio.Warnf("access was denied")
		}
		if res.Msg.Decision == accessv1alpha1.Decision_DECISION_REVIEW_REQUIRED {
			clio.Warnf("access requires review")
		}
		if res.Msg.Decision == accessv1alpha1.Decision_DECISION_ALLOWED {
			expiresIn := time.Until(res.Msg.ExpiresAt.AsTime())
			clio.Successf("access is active and expires in %s", expiresIn)
		}
		return nil
	},
}

type RequestsResponse struct {
	Requests []*accessv1alpha1.AccessRequest `json:"requests"`
}

var list = cli.Command{
	Name:  "list",
	Usage: "List Access Requests",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table', 'wide', or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		all := RequestsResponse{
			Requests: []*accessv1alpha1.AccessRequest{},
		}

		client := request.NewFromConfig(cfg)

		done := false
		var pageToken string

		for !done {
			res, err := client.QueryAccessRequests(ctx, connect.NewRequest(&accessv1alpha1.QueryAccessRequestsRequest{
				PageToken: pageToken,
			}))
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			all.Requests = append(all.Requests, res.Msg.AccessRequests...)

			if res.Msg.NextPageToken == "" {
				done = true
			} else {
				pageToken = res.Msg.NextPageToken
			}
		}

		output := c.String("output")
		switch output {
		case "table":
			w := table.New(os.Stdout)
			w.Columns("ID", "PRINCIPAL", "ROLE", "TARGET", "STATUS")

			for _, r := range all.Requests {
				for _, g := range r.Grants {
					w.Row(r.Id, formatUID(g.Principal), formatUID(g.Role), formatUID(g.Target), g.Status.String())
				}
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "wide":
			w := table.New(os.Stdout)
			w.Columns("ID", "GRANT", "PRINCIPAL", "ROLE", "TARGET", "STATUS")

			for _, r := range all.Requests {
				for _, g := range r.Grants {
					w.Row(r.Id, g.Id, formatUID(g.Principal), formatUID(g.Role), formatUID(g.Target), g.Status.String())
				}
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "json":
			resJSON, err := json.Marshal(all)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		default:
			return errors.New("invalid --output flag, valid values are [json, table, wide]")
		}

		return nil
	},
}

var review = cli.Command{
	Name:  "review",
	Usage: "Review access to a GCP entitlement",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "request-id", Required: true},
		&cli.BoolFlag{Name: "approve", Usage: "Approve the Access Request"},
		&cli.BoolFlag{Name: "close", Usage: "Close the Access Request"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := request.NewFromConfig(cfg)

		dec := accessv1alpha1.Review_REVIEW_UNSPECIFIED
		if c.Bool("approve") {
			dec = accessv1alpha1.Review_REVIEW_APPROVE
		}

		if c.Bool("close") {
			dec = accessv1alpha1.Review_REVIEW_CLOSE
		}

		if c.Bool("approve") && c.Bool("close") {
			return errors.New("you can't provide both --approve and --close: please provide either one only")
		}

		if dec == accessv1alpha1.Review_REVIEW_UNSPECIFIED {
			return errors.New("either --approve or --close must be provided")
		}

		res, err := client.ReviewAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.ReviewAccessRequestRequest{
			Id:     c.String("request-id"),
			Review: dec,
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}
		if c.Bool("approve") {
			clio.Successf("approved request")
		} else {
			clio.Successf("closed request")
		}

		return nil
	},
}

var revoke = cli.Command{
	Name:  "revoke",
	Usage: "Revoke access to a GCP entitlement",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "request-id"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		cfg.APIURL = "http://localhost:8080"
		client := request.NewFromConfig(cfg)

		res, err := client.RevokeAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.RevokeAccessRequestRequest{
			Id: c.String("request-id"),
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}
		clio.Successf("revoked request")

		return nil
	},
}

var cancel = cli.Command{
	Name:  "cancel",
	Usage: "Cancel access to a GCP entitlement",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "request-id"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		cfg.APIURL = "http://localhost:8080"
		client := request.NewFromConfig(cfg)

		res, err := client.CancelAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.CancelAccessRequestRequest{
			Id: c.String("request-id"),
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}
		clio.Successf("cancelled request")

		return nil
	},
}
