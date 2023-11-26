package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	accessv1alpha1 "github.com/common-fate/ciem/gen/commonfate/access/v1alpha1"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/ciem/service/access"
	"github.com/common-fate/ciem/table"
	"github.com/urfave/cli/v2"
)

var List = cli.Command{
	Name:  "list",
	Usage: "List available entitlements",
	Subcommands: []*cli.Command{
		&gcpList,
		&allCommand,
	},
}

type EntitlementsResponse struct {
	Entitlements []*accessv1alpha1.Entitlement `json:"entitlements"`
}

var allCommand = cli.Command{
	Name:  "all",
	Usage: "List available entitlements for all providers",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table', 'wide', or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		all := EntitlementsResponse{
			Entitlements: []*accessv1alpha1.Entitlement{},
		}
		done := false
		var pageToken string

		for !done {
			res, err := client.QueryEntitlements(ctx, connect.NewRequest(&accessv1alpha1.QueryEntitlementsRequest{
				PageToken: pageToken,
			}))
			if err != nil {
				return err
			}

			all.Entitlements = append(all.Entitlements, res.Msg.Entitlements...)

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
			w.Columns("TARGET", "NAME", "ROLE", "STATUS")

			for _, e := range all.Entitlements {
				var status string
				if e.Jit != nil {
					switch e.Jit.Status {
					case accessv1alpha1.JITStatus_JIT_STATUS_ACTIVE:
						status = "active"
					}
				}

				w.Row(formatUID(e.Target), e.TargetName, formatUID(e.Role), status)
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "wide":
			w := table.New(os.Stdout)
			w.Columns("TARGET", "NAME", "ROLE", "STATUS", "DURATION", "EXPIRES")

			for _, e := range all.Entitlements {
				var duration string
				if e.Jit != nil && e.Jit.Duration != nil {
					duration = e.Jit.Duration.AsDuration().String()
				}

				var status string
				if e.Jit != nil {
					switch e.Jit.Status {
					case accessv1alpha1.JITStatus_JIT_STATUS_ACTIVE:
						status = "active"
					}
				}

				var expiresIn string
				if e.Jit != nil && e.Jit.ExpiresIn != nil {
					expiresIn = e.Jit.ExpiresIn.AsDuration().String()
				}

				w.Row(formatUID(e.Target), e.TargetName, e.RoleName, status, duration, expiresIn)
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

func formatUID(uid *authzv1alpha1.UID) string {
	return fmt.Sprintf(`%s::%s`, uid.Type, uid.Id)
}

var gcpList = cli.Command{
	Name:  "gcp",
	Usage: "List available GCP entitlements",
	Action: func(c *cli.Context) error {
		// ctx := c.Context

		// cfg, err := config.LoadDefault(ctx)
		// if err != nil {
		// 	return err
		// }

		// client := access.NewFromConfig(cfg)

		// res, err := client.ListAccessRequests(ctx, connect.NewRequest(&accessv1alpha1.ListAccessRequestsRequest{}))
		// if err != nil {
		// 	return err
		// }

		// w := table.New(os.Stdout)
		// w.Columns("PROJECT", "ROLE")

		// for _, e := range res.Msg.AccessRequests {
		// 	switch v := e.Entitlement.Resource.Resource.(type) {
		// 	case *accessv1alpha1.Resource_AwsAccount:
		// 		w.Row(v.AwsAccount.AccountId, v.AwsAccount.Role)
		// 	case *accessv1alpha1.Resource_GcpProject:
		// 		w.Row(v.GcpProject.Project, v.GcpProject.Role)
		// 	}
		// }

		// w.Flush()

		return nil
	},
}
