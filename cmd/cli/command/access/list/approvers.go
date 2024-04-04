package list

import (
	"errors"
	"fmt"
	"os"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

var approversCommand = cli.Command{
	Name:    "approvers",
	Usage:   "List user principals which are able to approve the given entitlement.",
	Aliases: []string{"ap"},
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table',  or 'json')"},
		&cli.StringFlag{Name: "target", Usage: "Use this in conjunction with the --role flag. Can be either an ID, EID, name"},
		&cli.StringFlag{Name: "role", Usage: "Use this in conjunction with the --target flag. Can be either an ID, EID, name"},
		&cli.StringFlag{Name: "grant", Usage: `Use a grant id instead of a --target and --role. Can be either an ID or EID in the form Access::Grant::"<ID>"`},
	},
	Action: func(c *cli.Context) error {
		var query = &accessv1alpha1.QueryApproversRequest{}
		target := c.String("target")
		role := c.String("role")
		grantID := c.String("grant")

		if grantID != "" {
			if target != "" || role != "" {
				return errors.New("either --target and --role or --grant flags must be specified but not both")
			}
			id, err := eid.Parse(grantID)
			if err != nil {
				id = eid.New("Access::Grant", grantID)
				clio.Warnw("could not parse --grant as an EID", zap.Error(err))
				clio.Debugf("treating --grant as an ID because it could not be parsed as an EID, id has been formatted as: %s", id)
			} else if id.Type != "Access::Grant" {
				return errors.New(`--grant flags value should be eitehr an ID or and EID in the format Access::Grant::"<ID>"`)
			}
			query.Query = &accessv1alpha1.QueryApproversRequest_Grant{
				Grant: id.ToAPI(),
			}
		} else {
			if target == "" || role == "" {
				return errors.New("either --target and --role or --grant flags must be specified but not both")
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
