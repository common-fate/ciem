package list

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var availableCommand = cli.Command{
	Name:    "available",
	Usage:   "List available entitlements that access can be requested to",
	Aliases: []string{"av"},
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table', 'wide', or 'json')"},
		&cli.StringFlag{Name: "target-type", Usage: "filter for a particular target type"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		output := c.String("output")
		targetType := c.String("target-type")
		entitlements, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*accessv1alpha1.Entitlement, *string, error) {
			res, err := client.QueryEntitlements(ctx, connect.NewRequest(&accessv1alpha1.QueryEntitlementsRequest{
				PageToken:  grab.Value(nextToken),
				TargetType: grab.If(targetType == "", nil, &targetType),
			}))
			if err != nil {
				return nil, nil, err
			}
			return res.Msg.Entitlements, &res.Msg.NextPageToken, nil
		})
		if err != nil {
			return err
		}

		if output == "table" {

			w := table.New(os.Stdout)
			w.Columns("TARGET", "NAME", "ROLE")

			for _, e := range entitlements {
				w.Row(e.Target.Eid.Display(), e.Target.Name, e.Role.Name)
			}

			err = w.Flush()
			if err != nil {
				return err
			}
		} else {

			switch output {
			case "wide":
				w := table.New(os.Stdout)
				w.Columns("TARGET", "NAME", "ROLE", "AUTO-APPROVED", "TARGET PATH")

				for _, e := range entitlements {

					w.Row(e.Target.Eid.Display(), e.Target.Name, e.Role.Name, strconv.FormatBool(e.AutoApproved), strings.Join(grab.Map(e.TargetPath, func(p *accessv1alpha1.NamedEID) string { return p.Display() }), " / "))
				}

				err = w.Flush()
				if err != nil {
					return err
				}
			case "json":
				resJSON, err := protojson.Marshal(&accessv1alpha1.QueryEntitlementsResponse{Entitlements: entitlements})
				if err != nil {
					return err
				}
				fmt.Println(string(resJSON))
			default:
				return errors.New("invalid --output flag, valid values are [json, table, wide]")
			}

		}

		return nil
	},
}
