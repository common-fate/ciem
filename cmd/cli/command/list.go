package command

import (
	"os"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	accessv1alpha1 "github.com/common-fate/ciem/gen/commonfate/cloud/access/v1alpha1"
	"github.com/common-fate/ciem/service/access"
	"github.com/common-fate/ciem/table"
	"github.com/urfave/cli/v2"
)

var List = cli.Command{
	Name:  "list",
	Usage: "List available entitlements",
	Subcommands: []*cli.Command{
		&gcpList,
	},
}

var gcpList = cli.Command{
	Name:  "gcp",
	Usage: "List available GCP entitlements",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		res, err := client.ListEntitlementsForProvider(ctx, connect.NewRequest(&accessv1alpha1.ListEntitlementsForProviderRequest{
			Provider: accessv1alpha1.EntitlementProvider_ENTITLEMENT_PROVIDER_GCP,
		}))
		if err != nil {
			return err
		}

		w := table.New(os.Stdout)
		w.Columns("PROJECT", "ROLE")

		for _, e := range res.Msg.Entitlements {
			gcp := e.Resource.GetGcpProject()
			if gcp == nil {
				continue
			}
			w.Row(gcp.Project, gcp.Role)
		}

		w.Flush()

		return nil
	},
}
