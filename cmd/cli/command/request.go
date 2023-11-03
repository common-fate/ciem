package command

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	cf "github.com/common-fate/ciem/gen/proto/commonfatecloud/v1alpha1"
	"github.com/common-fate/ciem/service/access"
	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var Request = cli.Command{
	Name:  "request",
	Usage: "Request access to an entitlement",
	Subcommands: []*cli.Command{
		&gcpRequest,
	},
}

var gcpRequest = cli.Command{
	Name:  "gcp",
	Usage: "Request access to a GCP entitlement",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "project"},
		&cli.StringFlag{Name: "role"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		res, err := client.CreateAccessRequest(ctx, connect.NewRequest(&cf.CreateAccessRequestRequest{
			Entitlements: []*cf.Entitlement{
				{
					Target: &cf.Entitlement_Gcp{
						Gcp: &cf.GCPEntitlement{
							Project: c.String("project"),
							Role:    c.String("role"),
						},
					},
				},
			},
		}))
		if err != nil {
			return err
		}

		clio.Successf("created access request %s", res.Msg.AccessRequest.Id)

		return nil
	},
}
