package access

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/treeprint"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var Command = cli.Command{
	Name:  "access",
	Usage: "Request access to entitlements",
	Subcommands: []*cli.Command{
		&gcpRequest,
	},
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "target"},
		&cli.StringSliceFlag{Name: "role"},
		&cli.StringFlag{Name: "output", Value: "tree", Usage: "output format ('tree' or 'json')"},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context

		outputFormat := c.String("output")

		if outputFormat != "tree" && outputFormat != "json" {
			return errors.New("--output flag must be either 'tree' or 'json'")
		}

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		targets := c.StringSlice("target")
		roles := c.StringSlice("role")

		if len(targets) != len(roles) {
			return errors.New("you need to provide --role flag for each --target flag. For example:\n'cf jit request access --target AWS::Account::123456789012 --role AdministratorAccess --target OtherAccount --role Developer")
		}

		client := access.NewFromConfig(cfg)

		req := accessv1alpha1.BatchEnsureRequest{}

		for i, target := range targets {
			req.Entitlements = append(req.Entitlements, &accessv1alpha1.EntitlementInput{
				Target: &accessv1alpha1.Specifier{
					Specify: &accessv1alpha1.Specifier_Lookup{
						Lookup: target,
					},
				},
				Role: &accessv1alpha1.Specifier{
					Specify: &accessv1alpha1.Specifier_Lookup{
						Lookup: roles[i],
					},
				},
			})
		}

		res, err := client.BatchEnsure(ctx, connect.NewRequest(&req))
		if err != nil {
			return err
		}

		clio.Debugw("BatchEnsure response", "response", res)

		if outputFormat == "tree" {
			tree := treeprint.New()

			for i, req := range res.Msg.AccessRequests {
				var reqNode treeprint.Tree
				if !req.Existing {
					reqNode = tree.AddMetaBranch("CREATED ACCESS REQUEST", req.Id)
				} else {
					reqNode = tree.AddMetaBranch("EXISTING ACCESS REQUEST", req.Id)
				}

				for gi, g := range req.Grants {
					titleColor := color.New(color.FgWhite).SprintfFunc()

					switch g.Status {
					case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING_APPROVAL:
						titleColor = color.New(color.FgYellow).SprintfFunc()
					case accessv1alpha1.GrantStatus_GRANT_STATUS_APPROVED:
						titleColor = color.New(color.FgBlue).SprintfFunc()
					case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
						titleColor = color.New(color.FgGreen).SprintfFunc()
					case accessv1alpha1.GrantStatus_GRANT_STATUS_INACTIVE:
						titleColor = color.New(color.FgRed).SprintfFunc()
					}

					status := displayGrantStatus(g)

					grantLabel := titleColor("%s to %s", g.Role.Display(), g.Target.Display())

					grantNode := reqNode.AddMetaBranch(titleColor(status), grantLabel)
					// grantNode.AddGap()

					// targetNode := grantNode.AddBranch("Target")
					// targetNode.AddNode(g.Target.Uid.Display())
					// grantNode.AddGap()
					// roleNode := grantNode.AddBranch("Role")
					// roleNode.AddNode(g.Role.Uid.Display())

					// if g.Status == accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE && g.ExpiresAt != nil {
					// 	timingNode := grantNode.AddBranch("Timing")
					// 	timingNode.AddBranch(fmt.Sprintf("Expires In: %s", time.Until(g.ExpiresAt.AsTime()).Round(time.Second)))
					// }

					if g.ReviewersPreview != nil {
						if len(g.ReviewersPreview.Reviewers) == 0 {
							grantNode.AddNode("No Reviewers")
						} else {
							reviewers := grantNode.AddBranch("Reviewers")
							for _, r := range g.ReviewersPreview.Reviewers {

								reviewers.AddNode(r.Display())
							}
							// TODO: we can display a prompt like '5 more not shown' etc if the count of reviewers is greater than what the server sent
						}
					}

					if g.RecentActivity != nil && len(g.RecentActivity.Logs) > 0 {
						diags := grantNode.AddBranch("Activity")
						for _, d := range g.RecentActivity.Logs {
							color := colorAction(d.Action)
							diags.AddMetaNode(color(d.Action), fmt.Sprintf("%s: %s", d.OccurredAt.AsTime().Format(time.RFC3339), d.Message))
						}
					}

					// print a gap if there are more grants
					if gi < len(req.Grants)-1 {
						reqNode.AddGap()
					}
				}

				// print a gap if there are more requests
				if i < len(res.Msg.AccessRequests)-1 {
					tree.AddGap()
				}
			}

			fmt.Println(tree.String())

			if res.Msg.Warnings != nil {
				for _, w := range res.Msg.Warnings.Errors {
					clio.Error(w)
				}
			}

		}

		if outputFormat == "json" {
			resJSON, err := protojson.Marshal(res.Msg)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		}

		return nil
	},
}

func colorAction(action string) func(format string, a ...interface{}) string {
	if action == "grant.approved" {
		return color.New(color.FgGreen).SprintfFunc()
	}
	if action == "grant.self_approved" {
		return color.New(color.FgGreen).SprintfFunc()
	}
	if action == "grant.extended" {
		return color.New(color.FgGreen).SprintfFunc()
	}
	if action == "grant.activated" {
		return color.New(color.FgBlue).SprintfFunc()
	}
	if action == "grant.error" {
		return color.New(color.FgRed).SprintfFunc()
	}
	if action == "grant.cancelled" {
		return color.New(color.FgRed).SprintfFunc()
	}
	if action == "grant.revoked" {
		return color.New(color.FgRed).SprintfFunc()
	}

	return color.New(color.FgWhite).SprintfFunc()
}

func displayGrantStatus(g *accessv1alpha1.Grant) string {
	if g.Status == accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE && g.ExpiresAt != nil && g.ExpiresAt.AsTime().After(time.Now()) {
		exp := time.Until(g.ExpiresAt.AsTime()).Round(time.Minute)
		return fmt.Sprintf("Active for next %s", shortDur(exp))
	}

	switch g.Status {
	case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
		return "Active"
	case accessv1alpha1.GrantStatus_GRANT_STATUS_APPROVED:
		return "Approved"
	case accessv1alpha1.GrantStatus_GRANT_STATUS_INACTIVE:
		return "Inactive"
	case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING_APPROVAL:
		return "Pending"
	}

	return "<UNSPECIFIED STATUS>"
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
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

		res, err := client.Ensure(ctx, connect.NewRequest(&accessv1alpha1.EnsureRequest{
			Entitlement: &accessv1alpha1.EntitlementInput{
				Role: &accessv1alpha1.Specifier{
					Specify: &accessv1alpha1.Specifier_Uid{
						Uid: &authzv1alpha1.UID{
							Type: "GCP::Role",
							Id:   c.String("role"),
						},
					},
				},
				Target: &accessv1alpha1.Specifier{
					Specify: &accessv1alpha1.Specifier_Uid{
						Uid: &authzv1alpha1.UID{
							Type: "GCP::Project",
							Id:   c.String("id"),
						},
					},
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
					if g.ReviewersPreview == nil {
						continue
					}

					var reviewers []string

					for _, r := range g.ReviewersPreview.Reviewers {
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
			expiresAt := res.Msg.Grant.ExpiresAt
			if expiresAt != nil {

				if expiresAt.AsTime().Before(time.Now()) {
					clio.Warnf("access is active but was due to expire at %s. There may be a temporary problem with Common Fate. You should report this issue to your Common Fate administrator. Common Fate may remove your entitlement automatically.\nTo check whether the entitlement has been removed, you can run 'cf jit grant status --id %s'", expiresAt.AsTime(), res.Msg.Grant.Id)
				} else {
					expiresIn := time.Until(expiresAt.AsTime())
					clio.Successf("access is active and expires in %s", expiresIn)
				}
			} else {
				clio.Successf("access is active")
			}

		}
		return nil
	},
}
