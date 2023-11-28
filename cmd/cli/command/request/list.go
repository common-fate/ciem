package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/table"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/request"
	"github.com/urfave/cli/v2"
)

type RequestsResponse struct {
	Requests []*accessv1alpha1.AccessRequest `json:"requests"`
}

var list = cli.Command{
	Name:  "list",
	Usage: "List Access Requests",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table', 'wide', or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		all := RequestsResponse{
			Requests: []*accessv1alpha1.AccessRequest{},
		}

		client := request.NewFromConfig(cfg)

		done := false
		var pageToken string

		for !done {
			res, err := client.QueryAccessRequests(ctx, connect.NewRequest(&accessv1alpha1.QueryAccessRequestsRequest{
				PageToken: pageToken,
			}))
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			all.Requests = append(all.Requests, res.Msg.AccessRequests...)

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
			w.Columns("ID", "PRINCIPAL", "ROLE", "TARGET", "STATUS")

			for _, r := range all.Requests {
				for _, g := range r.Grants {
					w.Row(r.Id, g.Principal.Display(), g.Role.Display(), g.Target.Display(), g.Status.String())
				}
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "wide":
			w := table.New(os.Stdout)
			w.Columns("ID", "GRANT", "PRINCIPAL", "ROLE", "TARGET", "STATUS")

			for _, r := range all.Requests {
				for _, g := range r.Grants {
					w.Row(r.Id, g.Id, g.Principal.Display(), g.Role.Display(), g.Target.Display(), g.Status.String())
				}
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
