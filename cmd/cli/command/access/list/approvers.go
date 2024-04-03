package list

import (
	"errors"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var approversCommand = cli.Command{
	Name:    "approvers",
	Usage:   "List approvers for an entitlement",
	Aliases: []string{"ap"},
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table',  or 'json')"},
		&cli.StringFlag{Name: "target"},
		&cli.StringFlag{Name: "role"},
		&cli.StringFlag{Name: "grant-id"},
	},
	Action: func(c *cli.Context) error {
		var query = &accessv1alpha1.QueryApproversRequest{}
		target := c.String("target")
		role := c.String("role")
		grantID := c.String("grant-id")

		if grantID != "" {
			if target != "" || role != "" {
				return errors.New("either --target and --role or --grant-id flags must be specified but not both")
			}
			id, err := eid.Parse(grantID)
			if err != nil {
				return err
			}
			query.Query = &accessv1alpha1.QueryApproversRequest_Grant{
				Grant: id.ToAPI(),
			}
		} else {
			if target == "" || role == "" {
				return errors.New("either --target and --role or --grant-id flags must be specified but not both")
			}

			query.Query = &accessv1alpha1.QueryApproversRequest_TargetRole{
				TargetRole: &accessv1alpha1.TargetRole{
					Target: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: c.String("target"),
						},
					},
					Role: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: c.String("role"),
						},
					},
				},
			}
		}

		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		res, err := client.QueryApprovers(ctx, connect.NewRequest(query))
		if err != nil {
			return err
		}

		output := c.String("output")
		switch output {
		case "table":
			w := table.New(os.Stdout)
			w.Columns("ID", "NAME", "EMAIL")

			for _, approver := range res.Msg.Approvers {
				w.Row(approver.Eid.Display(), approver.Name, approver.Email)
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "json":
			resJSON, err := protojson.Marshal(res.Msg)
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
