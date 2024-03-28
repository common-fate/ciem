package list

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
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
		&cli.StringFlag{Name: "selector", Usage: "filter for a particular resource selector"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		output := c.String("output")

		if output == "table" {
			allEntitlements := accessv1alpha1.QueryEntitlementsResponse{
				Entitlements: []*accessv1alpha1.Entitlement{},
			}
			done := false
			var pageToken string

			selector := c.String("selector")

			for !done {
				res, err := client.QueryEntitlements(ctx, connect.NewRequest(&accessv1alpha1.QueryEntitlementsRequest{
					PageToken: pageToken,
				}))
				if err != nil {
					return err
				}

				for _, av := range res.Msg.Entitlements {
					if selector != "" {
						continue
					}
					allEntitlements.Entitlements = append(allEntitlements.Entitlements, av)

				}

				if res.Msg.NextPageToken == "" {
					done = true
				} else {
					pageToken = res.Msg.NextPageToken
				}
			}

			w := table.New(os.Stdout)
			w.Columns("TARGET", "NAME", "ROLE")

			for _, e := range allEntitlements.Entitlements {
				w.Row(e.Target.Eid.Display(), e.Target.Name, e.Role.Name)
			}

			err = w.Flush()
			if err != nil {
				return err
			}
		} else {
			allAvailabilities := accessv1alpha1.QueryAvailabilitiesResponse{
				Availabilities: []*accessv1alpha1.Availability{},
			}
			done := false
			var pageToken string

			selector := c.String("selector")

			for !done {
				res, err := client.QueryAvailabilities(ctx, connect.NewRequest(&accessv1alpha1.QueryAvailabilitiesRequest{
					PageToken: pageToken,
				}))
				if err != nil {
					return err
				}

				for _, av := range res.Msg.Availabilities {
					if selector != "" {
						continue
					}
					allAvailabilities.Availabilities = append(allAvailabilities.Availabilities, av)

				}

				if res.Msg.NextPageToken == "" {
					done = true
				} else {
					pageToken = res.Msg.NextPageToken
				}
			}

			switch output {
			case "wide":
				w := table.New(os.Stdout)
				w.Columns("TARGET", "NAME", "ROLE", "DURATION", "SELECTOR", "PRIORITY")

				for _, e := range allAvailabilities.Availabilities {
					w.Row(e.Target.Eid.Display(), e.Target.Name, e.Role.Name, e.Duration.AsDuration().String(), e.TargetSelector.Id, strconv.FormatUint(uint64(e.Priority), 10))
				}

				err = w.Flush()
				if err != nil {
					return err
				}
			case "json":
				resJSON, err := protojson.Marshal(&allAvailabilities)
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
