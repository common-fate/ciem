package access

import (
	"strings"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "access",
	Usage: "Request access to entitlements",
	Subcommands: []*cli.Command{
		&gcpRequest,
	},
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "target-type"},
		&cli.StringSliceFlag{Name: "target-id"},
		&cli.StringSliceFlag{Name: "role-type"},
		&cli.StringSliceFlag{Name: "role-id"},
	},
	// Action: func(c *cli.Context) error {
	// 	ctx := c.Context

	// 	cfg, err := config.LoadDefault(ctx)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	client := access.NewFromConfig(cfg)

	// 	req := &accessv1alpha1.GrantRequest{
	// 		Role: &authzv1alpha1.UID{
	// 			Type: "GCP::Role",
	// 			Id:   c.String("role"),
	// 		},
	// 		Target: &authzv1alpha1.UID{
	// 			Type: "GCP::Project",
	// 			Id:   c.String("id"),
	// 		},
	// 	}

	// 	res, err := client.Grant(ctx, connect.NewRequest(&accessv1alpha1.GrantRequest{
	// 		Role: &authzv1alpha1.UID{
	// 			Type: "GCP::Role",
	// 			Id:   c.String("role"),
	// 		},
	// 		Target: &authzv1alpha1.UID{
	// 			Type: "GCP::Project",
	// 			Id:   c.String("id"),
	// 		},
	// 	}))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	clio.Infow("response", "response", res)
	// 	if res.Msg.Decision == accessv1alpha1.Decision_DECISION_DENIED {
	// 		clio.Warnf("access was denied")
	// 	}
	// 	if res.Msg.Decision == accessv1alpha1.Decision_DECISION_REVIEW_REQUIRED {
	// 		clio.Warnf("access requires review")
	// 		var hasReviewers bool

	// 		if res.Msg.AccessRequest != nil {
	// 			for _, g := range res.Msg.AccessRequest.Grants {
	// 				var reviewers []string

	// 				for _, r := range g.Reviewers {
	// 					reviewers = append(reviewers, r.Display())
	// 				}

	// 				if len(reviewers) > 0 {
	// 					hasReviewers = true
	// 					clio.Warnf("access to %s with role %s will be reviewed by:\n%s", g.Target.Display(), g.Role.Display(), strings.Join(reviewers, "\n"))
	// 				}
	// 			}
	// 		}

	// 		if hasReviewers {
	// 			clio.Infof("reviewers can visit https://example.commonfate.cloud to approve access, or can run the following CLI command:\ncf request approve --id %s", res.Msg.AccessRequest.Id)
	// 		}
	// 	}
	// 	if res.Msg.Decision == accessv1alpha1.Decision_DECISION_ALLOWED {
	// 		expiresIn := time.Until(res.Msg.ExpiresAt.AsTime())
	// 		clio.Successf("access is active and expires in %s", expiresIn)
	// 	}
	// 	return nil
	// },
}

var gcpRequest = cli.Command{
	Name:  "gcp",
	Usage: "Get access to GCP entitlements",
	Subcommands: []*cli.Command{
		&gcpProjectRequest,
	},
}

var gcpProjectRequest = cli.Command{
	Name:  "project",
	Usage: "Get access to a GCP Project",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "role", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)

		res, err := client.Grant(ctx, connect.NewRequest(&accessv1alpha1.GrantRequest{
			Entitlement: &accessv1alpha1.EntitlementInput{
				Role: &authzv1alpha1.UID{
					Type: "GCP::Role",
					Id:   c.String("role"),
				},
				Target: &authzv1alpha1.UID{
					Type: "GCP::Project",
					Id:   c.String("id"),
				},
			},
		}))
		if err != nil {
			return err
		}

		clio.Infow("response", "response", res)
		if res.Msg.Grant.Status == accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING_APPROVAL {
			clio.Warnf("access requires review")
			var hasReviewers bool

			if res.Msg.AccessRequest != nil {
				for _, g := range res.Msg.AccessRequest.Grants {
					var reviewers []string

					for _, r := range g.Reviewers {
						reviewers = append(reviewers, r.Display())
					}

					if len(reviewers) > 0 {
						hasReviewers = true
						clio.Warnf("access to %s with role %s will be reviewed by:\n%s", g.Target.Display(), g.Role.Display(), strings.Join(reviewers, "\n"))
					}
				}
			}

			if hasReviewers {
				clio.Infof("reviewers can visit https://example.commonfate.cloud to approve access, or can run the following CLI command:\ncf request approve --id %s", res.Msg.AccessRequest.Id)
			}
		}
		if res.Msg.Grant.Status == accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE {
			clio.Successf("access is active and expires in <TODO>")
		}
		return nil
	},
}
