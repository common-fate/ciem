package command

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	accessv1alpha1 "github.com/common-fate/ciem/gen/commonfate/access/v1alpha1"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/ciem/service/access"
	"github.com/common-fate/ciem/service/request"
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
		if res.Msg.AccessRequest != nil && res.Msg.AccessRequest.Status == accessv1alpha1.RequestStatus_REQUEST_STATUS_ACTIVE {
			clio.Successf("access to %s with role %s is now active", res.Msg.AccessRequest.Resource.Uid.Id, res.Msg.AccessRequest.Action)
		}

		return nil
	},
}

var review = cli.Command{
	Name:  "review",
	Usage: "Review access to a GCP entitlement",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "request-id"},
		&cli.BoolFlag{Name: "approve", Value: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		cfg.APIURL = "http://localhost:8080"
		client := request.NewFromConfig(cfg)

		dec := accessv1alpha1.RequestReviewDecision_REQUEST_REVIEW_DECISION_CLOSE
		if c.Bool("approve") {
			dec = accessv1alpha1.RequestReviewDecision_REQUEST_REVIEW_DECISION_APPROVE
		}
		res, err := client.ReviewAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.ReviewAccessRequestRequest{
			Id:       c.String("request-id"),
			Decision: dec,
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
