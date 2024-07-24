package integration

import (
	"errors"

	"connectrpc.com/connect"
	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	integrationv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/v1alpha1"
	"github.com/common-fate/sdk/service/control/integration"
	"github.com/urfave/cli/v2"
)

var delete = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "integration ID to delete", Required: true},
		&cli.BoolFlag{Name: "confirm", Usage: "Confirm whether to delete the integration"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		confirm := c.Bool("confirm")
		if !confirm {
			err = survey.AskOne(&survey.Confirm{
				Message: "Please confirm you wish to delete the integration",
			}, &confirm)
			if err != nil {
				return err
			}
		}

		if !confirm {
			return errors.New("action cancelled by user")
		}
		client := integration.NewFromConfig(cfg)
		_, err = client.DeleteIntegration(ctx, connect.NewRequest(&integrationv1alpha1.DeleteIntegrationRequest{
			Id: c.String("id"),
		}))
		if err != nil {
			return err
		}

		clio.Success("deleted integration")
		return nil
	},
}
