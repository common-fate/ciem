package command

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	accessv1alpha1 "github.com/common-fate/ciem/gen/commonfate/cloud/access/v1alpha1"
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

		res, err := client.CreateAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.CreateAccessRequestRequest{
			Resources: []*accessv1alpha1.Resource{
				{
					Resource: &accessv1alpha1.Resource_GcpProject{
						GcpProject: &accessv1alpha1.GCPProject{
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

		for _, e := range res.Msg.RequestGroup.Requests {
			gcp := e.Entitlement.Resource.GetGcpProject()
			if gcp == nil {
				continue
			}

			if e.Entitlement.Status == accessv1alpha1.EntitlementStatus_ENTITLEMENT_STATUS_ACTIVE {
				clio.Successf("access to %s with role %s is now active", gcp.Project, gcp.Role)
			}
		}

		return nil
	},
}
