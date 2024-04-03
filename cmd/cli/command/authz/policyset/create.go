package policyset

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/service/authz/policyset"
	"github.com/urfave/cli/v2"
)

var createCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "file", Usage: "Policy Set file to create", Required: true},
		&cli.StringFlag{Name: "id", Usage: "Policy Set ID", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := policyset.NewFromConfig(cfg)

		f, err := os.ReadFile(c.String("file"))
		if err != nil {
			return err
		}

		_, err = client.Create(ctx, policyset.CreateInput{
			PolicySet: policyset.Input{
				ID:   c.String("id"),
				Text: string(f),
			},
		})
		if err != nil {
			return err
		}

		clio.Success("created policyset")
		return nil
	},
}
