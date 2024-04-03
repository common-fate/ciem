package preview

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var userAccess = cli.Command{
	Name:    "user-access",
	Aliases: []string{"ua"},
	Usage:   "Preview which entitlements a given user principal can request access to, and whether they require approval.",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table',  or 'json')"},
		&cli.StringFlag{Name: "principal", Required: true, Usage: "Can be either an ID, EID, email, full name"},
		&cli.StringFlag{Name: "target-type", Usage: "E.g AWS::Account will filter for the aws account availabilities"},
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

		targetType := c.String("target-type")

		res, err := client.PreviewEntitlements(ctx, connect.NewRequest(&accessv1alpha1.PreviewEntitlementsRequest{
			Principal: &accessv1alpha1.Specifier{
				Specify: &accessv1alpha1.Specifier_Lookup{
					Lookup: c.String("principal"),
				},
			},
			TargetType: grab.If(targetType == "", nil, &targetType),
		}))
		if err != nil {
			return err
		}

		output := c.String("output")
		switch output {
		case "table":
			w := table.New(os.Stdout)
			w.Columns("TARGET", "NAME", "ROLE", "REQUIRES APPROVAL")

			for _, entitlement := range res.Msg.Entitlements {
				w.Row(entitlement.Target.Eid.Display(), entitlement.Target.Name, entitlement.Role.Name, strconv.FormatBool(!entitlement.AutoApproved))
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
