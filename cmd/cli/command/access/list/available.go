package list

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/table"
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
		&cli.StringFlag{Name: "selector", Usage: "filter for a particular resource selector"},
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
		availabilities, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*accessv1alpha1.Availability, *string, error) {
			res, err := client.QueryAvailabilities(ctx, connect.NewRequest(&accessv1alpha1.QueryAvailabilitiesRequest{
				PageToken: grab.Value(nextToken),
			}))
			if err != nil {
				return nil, nil, err
			}
			return res.Msg.Availabilities, &res.Msg.NextPageToken, nil
		})
		if err != nil {
			return err
		}

		selector := c.String("selector")
		for _, av := range availabilities {
			if selector != "" && av.TargetSelector.Id != selector {
				continue
			}
			all.Availabilities = append(all.Availabilities, av)
		}

		output := c.String("output")
		switch output {
		case "table":
			w := table.New(os.Stdout)
			w.Columns("TARGET", "NAME", "ROLE", "DURATION")

			for _, e := range all.Availabilities {
				w.Row(e.Target.Eid.Display(), e.Target.Name, e.Role.Name, e.Duration.AsDuration().String())
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "wide":
			w := table.New(os.Stdout)
			w.Columns("TARGET", "NAME", "ROLE", "DURATION", "SELECTOR", "PRIORITY")

			for _, e := range all.Availabilities {
				w.Row(e.Target.Eid.Display(), e.Target.Name, e.Role.Name, e.Duration.AsDuration().String(), e.TargetSelector.Id, strconv.FormatUint(uint64(e.Priority), 10))
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
