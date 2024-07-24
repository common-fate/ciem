package integration

import (
	"os"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/sdk/config"
	integrationv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/v1alpha1"
	"github.com/common-fate/sdk/service/control/integration"
	"github.com/urfave/cli/v2"
)

var list = cli.Command{
	Name: "list",

	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := integration.NewFromConfig(cfg)
		res, err := client.ListIntegrations(ctx, connect.NewRequest(&integrationv1alpha1.ListIntegrationsRequest{}))
		if err != nil {
			return err
		}
		w := table.New(os.Stdout)
		w.Columns("ID", "NAME", "STATUS")

		for _, integration := range res.Msg.Integrations {

			w.Row(integration.Id, integration.Name, integration.Status.String())
		}

		err = w.Flush()
		if err != nil {
			return err
		}
		return nil
	},
}
