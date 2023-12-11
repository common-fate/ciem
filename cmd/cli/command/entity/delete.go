package entity

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	"github.com/common-fate/sdk/service/entity"
	"github.com/urfave/cli/v2"
)

var deleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "eid", Usage: "Entity ID to delete"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := entity.NewFromConfig(cfg)

		var uids []eid.EID

		for _, e := range c.StringSlice("eid") {
			id, err := eid.Parse(e)
			if err != nil {
				return err
			}
			uids = append(uids, id)
		}

		_, err = client.BatchUpdate(ctx, entity.BatchUpdateInput{
			Delete: uids,
		})
		if err != nil {
			return err
		}

		clio.Successf("deleted %v entities", len(uids))
		return nil
	},
}
