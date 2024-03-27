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
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		all := accessv1alpha1.QueryAvailabilitiesResponse{
			Availabilities: []*accessv1alpha1.Availability{},
		}

		target, err := eid.Parse(c.String("target"))
		if err != nil {
			return err
		}
		role, err := eid.Parse(c.String("role"))
		if err != nil {
			return err
		}
		res, err := client.QueryApprovers(ctx, connect.NewRequest(&accessv1alpha1.QueryApproversRequest{
			Target: target.ToAPI(),
			Role:   role.ToAPI(),
		}))
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
			resJSON, err := protojson.Marshal(&all)
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
