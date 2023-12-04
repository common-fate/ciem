package entities

import (
	"fmt"

	"github.com/common-fate/sdk/service/entity"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var getCommand = cli.Command{
	Name: "get",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "type", Usage: "entity type to load", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		client := entity.NewClient(entity.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		var all entity.ListOutput

		call := client.ListRequest(entity.ListInput{
			Type: c.String("type"),
		})

		call.Pages(ctx, func(lo *entity.ListOutput) error {
			all.Entities = append(all.Entities, lo.Entities...)
			return nil
		})

		out, err := protojson.Marshal(&all)
		if err != nil {
			return err
		}

		fmt.Println(string(out))

		return nil
	},
}
