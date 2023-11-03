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
			Resources: []*cf.Resource{
				{
					Resource: &cf.Resource_GcpProject{
						GcpProject: &cf.GCPProject{
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

		for _, e := range res.Msg.AccessRequest.Entitlements {
			gcp := e.Resource.GetGcpProject()
			if gcp == nil {
				continue
			}

			if e.Status == cf.EntitlementStatus_ENTITLEMENT_STATUS_ACTIVE {
				clio.Successf("access to %s with role %s is now active", gcp.Project, gcp.Role)
			}
		}

		return nil
	},
}
