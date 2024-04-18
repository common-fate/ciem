package access

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/briandowns/spinner"
	"github.com/common-fate/cli/awsconfig"
	"github.com/common-fate/cli/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	awsv1alpha1 "github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1"
	"github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1/awsv1alpha1connect"
	grantedv1alpha1 "github.com/common-fate/sdk/service/granted/registry"

	"github.com/common-fate/sdk/service/access"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
)

var ensureCommand = cli.Command{
	Name:  "ensure",
	Usage: "Ensure access to some entitlements (will request, active, or extend access as necessary)",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "target", Required: true},
		&cli.StringSliceFlag{Name: "role", Required: true},
		&cli.DurationFlag{Name: "duration", Aliases: []string{"d"}, Usage: "Override the default duration with a custom one. Must be less than the max duration available."},
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text' or 'json')"},
		&cli.StringFlag{Name: "reason", Usage: "The reason for requesting access"},
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "skip the confirmation prompt"},
		&cli.BoolFlag{Name: "skip-local-config-update", Usage: "Skip updating any local configuration files, such as '~/aws/config'"},
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
		duration := c.Duration("duration")

		if len(targets) != len(roles) {
			return errors.New("you need to provide --role flag for each --target flag. For example:\n'cf jit request access --target AWS::Account::123456789012 --role AdministratorAccess --target OtherAccount --role Developer")
		}

		apiURL, err := url.Parse(cfg.APIURL)
		if err != nil {
			return err
		}

		client := access.NewFromConfig(cfg)
		reason := c.String("reason")
		req := accessv1alpha1.BatchEnsureRequest{
			Justification: &accessv1alpha1.Justification{
				Reason: grab.If(reason == "", nil, &reason),
			},
		}

		for i, target := range targets {

			ent := &accessv1alpha1.EntitlementInput{
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
			}

			if duration > 0 {
				ent.Duration = durationpb.New(duration)
			}
			req.Entitlements = append(req.Entitlements, ent)

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

		//build up profiles based on the requested entitlements and put them into ini format

		//call granteds merge method to merge together the new profiles and the current aws config state. This should handle duplicates

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

		//prints response diag messages
		printdiags.Print(res.Msg.Diagnostics, nil)

		si.Stop()

		clio.Debugw("BatchEnsure response", "response", res)

		var isAWSAccess bool

		for _, g := range res.Msg.Grants {
			if g.Grant.Target.Eid.Type == "AWS::Account" {
				isAWSAccess = true
				break
			}
		}

		disableConfig := c.Bool("skip-local-config-update")

		if !disableConfig && isAWSAccess {
			// update any local config files like ~/.aws/config
			accountClient := grantedv1alpha1.NewFromConfig(cfg)

			err = updateAWSProfiles(ctx, res.Msg.Grants, accountClient)
			if err != nil {
				clio.Errorf("error updating local AWS config: %s", err.Error())
			}
		}

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

func updateAWSProfiles(ctx context.Context, grants []*accessv1alpha1.GrantState, accountClient awsv1alpha1connect.ProfileRegistryServiceClient) error {
	awsConfig, filePath, err := awsconfig.Load()
	if err != nil {
		return err
	}

	for _, g := range grants {
		if g.Grant.Target.Eid.Type != "AWS::Account" {
			continue
		}

		accountID := g.Grant.Target.Eid.Id
		roleName := g.Grant.Role.Name

		profileFromCF, err := accountClient.GetProfileForAccountAndRole(ctx, &connect.Request[awsv1alpha1.GetProfileForAccountAndRoleRequest]{
			Msg: &awsv1alpha1.GetProfileForAccountAndRoleRequest{
				AccountId: g.Grant.Target.Eid.Id,
				RoleName:  g.Grant.Role.Name,
			},
		})
		if connect.CodeOf(err) == connect.CodeNotFound || connect.CodeOf(err) == connect.CodeUnimplemented {
			clio.Warnf("could not update AWS profile in %s for account %s and role %s because this version of Common Fate does not implement the Granted Profile Registry API, you may need to update your deployment: %s", filePath, accountID, roleName, err.Error())
			continue
		}
		if err != nil {
			clio.Warnf("could not update AWS profile in %s for account %s and role %s: %s", filePath, accountID, roleName, err.Error())
			continue
		}

		clio.Infof("adding %s to your AWS config file (%s)\t[target=%s, role=%s]", profileFromCF.Msg.Profile.Name, filePath, g.Grant.Target.Display(), roleName)

		// build up a new section for each profile being added
		err = awsconfig.Merge(awsconfig.MergeOpts{
			Config:            awsConfig,
			ProfileName:       profileFromCF.Msg.Profile.Name,
			ProfileAttributes: profileFromCF.Msg.Profile.Attributes,
		})
		if err != nil {
			return err
		}
	}

	err = awsConfig.SaveTo(filePath)
	if err != nil {
		return err
	}

	return nil
}
