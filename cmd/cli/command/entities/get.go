package entities

import (
	"encoding/json"
	"fmt"

	"github.com/common-fate/sdk/service/authz"
	"github.com/urfave/cli/v2"
)

var getCommand = cli.Command{
	Name: "get",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := authz.NewClient(authz.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		entities, err := client.Query(ctx, authz.QueryInput{})
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
