package auditlog

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/table"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/audit"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var listCommand = cli.Command{
	Name: "list",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table', 'wide', or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := audit.NewFromConfig(cfg)

		all := accessv1alpha1.QueryAuditLogsResponse{
			AuditLogs: []*accessv1alpha1.AuditLog{},
		}

		done := false
		var pageToken string

		for !done {
			res, err := client.QueryAuditLogs(ctx, connect.NewRequest(&accessv1alpha1.QueryAuditLogsRequest{
				PageToken: pageToken,
			}))
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			all.AuditLogs = append(all.AuditLogs, res.Msg.AuditLogs...)

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
			w.Columns("TIME", "ACTOR", "ACTION", "MESSAGE", "TARGET")

			for _, r := range all.AuditLogs {
				w.Row(r.OccurredAt.AsTime().Format(time.RFC3339), r.Actor.Display(), r.Action, r.Message, r.PrimaryTarget.Display())

			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "wide":
			w := table.New(os.Stdout)
			w.Columns("TIME", "ACTOR", "ACTION", "MESSAGE", "TARGETS")

			for _, r := range all.AuditLogs {
				var targets []string
				for _, t := range r.AllTargets {
					targets = append(targets, t.Display())
				}
				w.Row(r.OccurredAt.AsTime().Format(time.RFC3339), r.Actor.Display(), r.Action, r.Message, strings.Join(targets, ", "))

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
