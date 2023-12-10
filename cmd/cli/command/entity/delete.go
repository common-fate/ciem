package entity

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/service/entity"
	"github.com/common-fate/sdk/uid"
	"github.com/urfave/cli/v2"
)

var deleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "eid", Usage: "Entity ID to delete"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := entity.NewClient(entity.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		var uids []uid.UID

		for _, e := range c.StringSlice("eid") {
			id, err := uid.Parse(e)
			if err != nil {
				return err
			}
			uids = append(uids, id)
		}

		_, err := client.BatchUpdate(ctx, entity.BatchUpdateInput{
			Delete: uids,
		})
		if err != nil {
			return err
		}

		clio.Successf("deleted %v entities", len(uids))
		return nil
	},
}
