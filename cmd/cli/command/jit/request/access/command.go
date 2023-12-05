package access

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/ciem/treeprint"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/common-fate/sdk/uid"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var Command = cli.Command{
	Name:  "access",
	Usage: "Request access to entitlements",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "target", Required: true},
		&cli.StringSliceFlag{Name: "role", Required: true},
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

			names := map[uid.UID]string{}

			for i, req := range res.Msg.AccessRequests {
				var reqNode treeprint.Tree
				if !req.Existing {
					reqNode = tree.AddMetaBranch("CREATED ACCESS REQUEST", req.Id)
				} else {
					reqNode = tree.AddMetaBranch("EXISTING ACCESS REQUEST", req.Id)
				}

				for gi, g := range req.Grants {
					names[uid.New("JIT::Grant", g.Id)] = g.Name

					titleColor := color.New(color.FgWhite).SprintfFunc()

					switch g.Status {
					case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING:
						titleColor = color.New(color.FgYellow).SprintfFunc()
					case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
						titleColor = color.New(color.FgGreen).SprintfFunc()
					case accessv1alpha1.GrantStatus_GRANT_STATUS_INACTIVE:
						titleColor = color.New(color.FgRed).SprintfFunc()
					}

					status := displayGrantStatus(g)

					grantLabel := titleColor(g.Name)

					reqNode.AddMetaBranch(titleColor(status), grantLabel)

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

			printdiags.Print(res.Msg.Diagnostics, names)

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
	case accessv1alpha1.GrantStatus_GRANT_STATUS_INACTIVE:
		return "Inactive"
	case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING:
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
