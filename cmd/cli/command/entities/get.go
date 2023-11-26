package entities

import (
	"encoding/json"
	"fmt"

	"github.com/common-fate/ciem/service/authz"
	"github.com/urfave/cli/v2"
)

var getCommand = cli.Command{
	Name: "get",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := authz.NewClient(newInsecureClient(), "http://127.0.0.1:5050")

		entities, err := client.FilterEntities(ctx, authz.FilterEntitiesInput{})
		if err != nil {
			return err
		}

		out, err := json.MarshalIndent(entities.Entities, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(out))

		return nil
	},
}
