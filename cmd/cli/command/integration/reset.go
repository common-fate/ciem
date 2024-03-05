package integration

import (
	"connectrpc.com/connect"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	resetv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/reset/v1alpha1"
	"github.com/common-fate/sdk/service/control/integration/reset"
	"github.com/urfave/cli/v2"
)

var resetCommand = cli.Command{
	Name:        "reset",
	Usage:       "used to remove all entities for a specified integration",
	Subcommands: []*cli.Command{&resetEntraCommand},
}

var resetEntraCommand = cli.Command{
	Name:  "entra",
	Usage: "Removes all Entra user entities in a deployment",
	Flags: []cli.Flag{&cli.BoolFlag{Name: "dry", Usage: "using dry run will not remove any entities and return back what will be removed from executing the command", Required: true}},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		isDryRun := c.Bool("dry")
		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := reset.NewFromConfig(cfg)

		out, err := client.ResetEntraUsers(ctx, &connect.Request[resetv1alpha1.ResetEntraUsersRequest]{
			Msg: &resetv1alpha1.ResetEntraUsersRequest{
				DryRun: isDryRun,
			},
		})
		if err != nil {
			return err
		}

		if isDryRun {
			clio.Infof("The following entities will be deleted (Dry run: %t)", isDryRun)
		} else {
			clio.Info("The following entities have been deleted:")
		}

		for _, ent := range out.Msg.DeletedEntities {
			clio.Infof("- %s::%s", ent.Type, ent.Id)
		}
		if !isDryRun {
			clio.Success("Complete")
		}

		return nil
	},
}
