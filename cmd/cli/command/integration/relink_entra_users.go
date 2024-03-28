package integration

import (
	"connectrpc.com/connect"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	resetv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/reset/v1alpha1"
	"github.com/common-fate/sdk/service/control/integration/reset"
	"github.com/urfave/cli/v2"
)

var relinkEntraUsers = cli.Command{
	Name:  "relink-entra-users",
	Usage: "links together Common Fate users with their Entra user counterpart. Used to remediate any issues which may have come up during initial sync of Entra users.",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := reset.NewFromConfig(cfg)

		_, err = client.RelinkEntraUsers(ctx, connect.NewRequest(&resetv1alpha1.RelinkEntraUsersRequest{}))
		if err != nil {
			return err
		}

		clio.Success("relinked Entra::User entities to CF::User entities")
		return nil
	},
}
