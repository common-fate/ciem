package entity

import (
	"encoding/json"
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/service/entity"
	"github.com/urfave/cli/v2"
)

var putCommand = cli.Command{
	Name: "put",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := entity.NewFromConfig(cfg)

		f, err := os.ReadFile(c.Path("file"))
		if err != nil {
			return err
		}

		var entities []entity.EntityJSON

		err = json.Unmarshal(f, &entities)
		if err != nil {
			return err
		}

		_, err = client.BatchPutJSON(ctx, entity.BatchPutJSONInput{
			Entities: entities,
		})
		if err != nil {
			return err
		}

		clio.Successf("put %v entities", len(entities))
		return nil
	},
}
