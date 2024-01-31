package access

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/briandowns/spinner"
	"github.com/common-fate/cli/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var ensureCommand = cli.Command{
	Name:  "ensure",
	Usage: "Ensure access to some entitlements (will request, active, or extend access as necessary)",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "target", Required: true},
		&cli.StringSliceFlag{Name: "role", Required: true},
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text' or 'json')"},
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "skip the confirmation prompt"},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context

		outputFormat := c.String("output")

		if outputFormat != "text" && outputFormat != "json" {
			return errors.New("--output flag must be either 'text' or 'json'")
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

		apiURL, err := url.Parse(cfg.APIURL)
		if err != nil {
			return err
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

		if !c.Bool("confirm") {
			jsonOutput := c.String("output") == "json"

			// run the dry-run first
			hasChanges, err := DryRun(ctx, apiURL, client, &req, jsonOutput)
			if err != nil {
				return err
			}
			if !hasChanges {
				fmt.Println("no access changes")
				return nil
			}
		}

		// if we get here, dry-run has passed the user has confirmed they want to proceed.
		req.DryRun = false

		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " ensuring access..."
		si.Writer = os.Stderr
		si.Start()

		res, err := client.BatchEnsure(ctx, connect.NewRequest(&req))
		if err != nil {
			si.Stop()
			return err
		}

		si.Stop()

		clio.Debugw("BatchEnsure response", "response", res)

		if outputFormat == "text" {

			// tree := treeprint.New()

			names := map[eid.EID]string{}

			for _, g := range res.Msg.Grants {
				names[eid.New("Access::Grant", g.Grant.Id)] = g.Grant.Name

				exp := "<invalid expiry>"

				if g.Grant.ExpiresAt != nil {
					exp = ShortDur(time.Until(g.Grant.ExpiresAt.AsTime()))
				}

				switch g.Change {
				case accessv1alpha1.GrantChange_GRANT_CHANGE_ACTIVATED:
					color.New(color.BgHiGreen).Printf("[ACTIVATED]")
					color.New(color.FgGreen).Printf(" %s was activated for %s: %s\n", g.Grant.Name, exp, RequestURL(apiURL, g.Grant))
					continue

				case accessv1alpha1.GrantChange_GRANT_CHANGE_EXTENDED:
					color.New(color.BgBlue).Printf("[EXTENDED]")
					color.New(color.FgBlue).Printf(" %s was extended for another %s: %s\n", g.Grant.Name, exp, RequestURL(apiURL, g.Grant))
					continue

				case accessv1alpha1.GrantChange_GRANT_CHANGE_REQUESTED:
					color.New(color.BgHiYellow, color.FgBlack).Printf("[REQUESTED]")
					color.New(color.FgYellow).Printf(" %s requires approval: %s\n", g.Grant.Name, RequestURL(apiURL, g.Grant))
					continue

				case accessv1alpha1.GrantChange_GRANT_CHANGE_PROVISIONING_FAILED:
					// shouldn't happen in the dry-run request but handle anyway
					color.New(color.FgRed).Printf("[ERROR] %s failed provisioning: %s\n", g.Grant.Name, RequestURL(apiURL, g.Grant))
					continue
				}

				switch g.Grant.Status {
				case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
					color.New(color.FgGreen).Printf("[ACTIVE] %s is already active for the next %s: %s\n", g.Grant.Name, exp, RequestURL(apiURL, g.Grant))
					continue
				case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING:
					color.New(color.FgWhite).Printf("[PENDING] %s is already pending: %s\n", g.Grant.Name, RequestURL(apiURL, g.Grant))
					continue
				case accessv1alpha1.GrantStatus_GRANT_STATUS_CLOSED:
					color.New(color.FgWhite).Printf("[CLOSED] %s is closed but was still returned: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, RequestURL(apiURL, g.Grant))
					continue
				}

				color.New(color.FgWhite).Printf("[UNSPECIFIED] %s is in an unspecified status: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, RequestURL(apiURL, g.Grant))
			}

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

func ShortDur(d time.Duration) string {
	if d > time.Minute {
		d = d.Round(time.Minute)
	} else {
		d = d.Round(time.Second)
	}

	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
