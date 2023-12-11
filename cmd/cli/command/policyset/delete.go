package policyset

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/service/authz/policyset"
	"github.com/urfave/cli/v2"
)

var deleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "PolicySet ID to delete", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := policyset.NewFromConfig(cfg)

		_, err = client.Delete(ctx, policyset.DeleteInput{
			ID: c.String("id"),
		})
		if err != nil {
			return err
		}

		clio.Success("deleted policyset")
		return nil
	},
}
