package integration

import (
	"connectrpc.com/connect"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	resetv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/reset/v1alpha1"
	"github.com/common-fate/sdk/service/control/integration/reset"
	"github.com/urfave/cli/v2"
)

var deleteOauthTokenCommand = cli.Command{
	Name: "delete-oauth-token",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "token ID to delete", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := reset.NewFromConfig(cfg)

		_, err = client.RemoveOAuthToken(ctx, connect.NewRequest(&resetv1alpha1.RemoveOAuthTokenRequest{
			Id: c.String("id"),
		}))
		if err != nil {
			return err
		}

		clio.Success("deleted oauth token")
		return nil
	},
}
