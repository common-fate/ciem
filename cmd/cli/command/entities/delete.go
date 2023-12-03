package entities

import (
	"encoding/json"
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/service/entity"
	"github.com/common-fate/sdk/uid"
	"github.com/urfave/cli/v2"
)

var deleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := entity.NewClient(entity.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		f, err := os.ReadFile(c.Path("file"))
		if err != nil {
			return err
		}

		var entities []entity.EntityJSON

		err = json.Unmarshal(f, &entities)
		if err != nil {
			return err

		}

		var uids []uid.UID

		for _, e := range entities {
			uids = append(uids, e.UID)
		}

		_, err = client.BatchUpdate(ctx, entity.BatchUpdateInput{
			Delete: uids,
		})
		if err != nil {
			return err
		}

		clio.Successf("deleted %v entities", len(entities))
		return nil
	},
}
