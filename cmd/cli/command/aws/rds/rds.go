package rds

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/bufbuild/connect-go"
	accessCmd "github.com/common-fate/ciem/cmd/cli/command/access"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/grab"
	"github.com/common-fate/granted/pkg/assume"
	"github.com/common-fate/granted/pkg/cfaws"
	"github.com/fatih/color"

	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	entityv1alpha1 "github.com/common-fate/sdk/gen/commonfate/entity/v1alpha1"
	"github.com/common-fate/sdk/service/access"
	"github.com/common-fate/sdk/service/access/grants"
	"github.com/common-fate/sdk/service/entity"
	"github.com/urfave/cli/v2"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

var Command = cli.Command{
	Name:  "rds",
	Usage: "Perform RDS Operations",
	Subcommands: []*cli.Command{
		&proxyCommand,
	},
}

var proxyCommand = cli.Command{
	Name:  "proxy",
	Usage: "Run a database proxy",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "target", Required: true},
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "skip the confirmation prompt"},
		&cli.IntFlag{Name: "port", Value: 5432, Usage: "The local port to forward the database connect to"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}
		target := c.String("target")

		client := access.NewFromConfig(cfg)
		apiURL, err := url.Parse(cfg.APIURL)
		if err != nil {
			return err
		}
		req := accessv1alpha1.BatchEnsureRequest{
			Entitlements: []*accessv1alpha1.EntitlementInput{
				{
					Target: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Lookup{
							Lookup: target,
						},
					},
					Role: &accessv1alpha1.Specifier{
						Specify: &accessv1alpha1.Specifier_Eid{
							Eid: &entityv1alpha1.EID{
								Type: "CF::Database::Role",
								Id:   "ReadWrite",
							},
						},
					},
				},
			},
		}

		if !c.Bool("confirm") {

			// run the dry-run first
			hasChanges, err := accessCmd.DryRun(ctx, apiURL, client, &req, false)
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

		// tree := treeprint.New()

		names := map[eid.EID]string{}

		for _, g := range res.Msg.Grants {
			names[eid.New("Access::Grant", g.Grant.Id)] = g.Grant.Name

			exp := "<invalid expiry>"

			if g.Grant.ExpiresAt != nil {
				exp = accessCmd.ShortDur(time.Until(g.Grant.ExpiresAt.AsTime()))
			}

			switch g.Change {
			case accessv1alpha1.GrantChange_GRANT_CHANGE_ACTIVATED:
				color.New(color.BgHiGreen).Printf("[ACTIVATED]")
				color.New(color.FgGreen).Printf(" %s was activated for %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
				continue

			case accessv1alpha1.GrantChange_GRANT_CHANGE_EXTENDED:
				color.New(color.BgBlue).Printf("[EXTENDED]")
				color.New(color.FgBlue).Printf(" %s was extended for another %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
				continue

			case accessv1alpha1.GrantChange_GRANT_CHANGE_REQUESTED:
				color.New(color.BgHiYellow, color.FgBlack).Printf("[REQUESTED]")
				color.New(color.FgYellow).Printf(" %s requires approval: %s\n", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
				continue

			case accessv1alpha1.GrantChange_GRANT_CHANGE_PROVISIONING_FAILED:
				// shouldn't happen in the dry-run request but handle anyway
				color.New(color.FgRed).Printf("[ERROR] %s failed provisioning: %s\n", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
				continue
			}

			switch g.Grant.Status {
			case accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE:
				color.New(color.FgGreen).Printf("[ACTIVE] %s is already active for the next %s: %s\n", g.Grant.Name, exp, accessCmd.RequestURL(apiURL, g.Grant))
				continue
			case accessv1alpha1.GrantStatus_GRANT_STATUS_PENDING:
				color.New(color.FgWhite).Printf("[PENDING] %s is already pending: %s\n", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
				continue
			case accessv1alpha1.GrantStatus_GRANT_STATUS_CLOSED:
				color.New(color.FgWhite).Printf("[CLOSED] %s is closed but was still returned: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
				continue
			}

			color.New(color.FgWhite).Printf("[UNSPECIFIED] %s is in an unspecified status: %s\n. This is most likely due to an error in Common Fate and should be reported to our team: support@commonfate.io.", g.Grant.Name, accessCmd.RequestURL(apiURL, g.Grant))
		}

		printdiags.Print(res.Msg.Diagnostics, names)

		ensuredGrant := res.Msg.Grants[0]
		// if its not yet active, we can just exit the process
		if ensuredGrant.Grant.Status != accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE {
			clio.Debug("grant not yet active, exiting")
			return nil
		}

		grantsClient := grants.NewFromConfig(cfg)
		// idClient := identity.NewFromConfig(cfg)

		// idRes, err := idClient.GetCallerIdentity(ctx, connect.NewRequest(&accessv1alpha1.GetCallerIdentityRequest{}))
		// if err != nil {
		// 	return err
		// }

		// matchedGrants, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*accessv1alpha1.Grant, *string, error) {
		// 	res, err := grantsClient.QueryGrants(ctx, connect.NewRequest(&accessv1alpha1.QueryGrantsRequest{
		// 		PageToken: grab.Value(nextToken),
		// 		Principal: idRes.Msg.Principal.Eid,
		// 		Target: &entityv1alpha1.EID{
		// 			Type: "AWS::RDS::Instance",
		// 			Id:   target,
		// 		},
		// 		Role: &entityv1alpha1.EID{
		// 			Type: "CF::Database::Role",
		// 			Id:   "ReadWrite",
		// 		},
		// 		Status: accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE.Enum(),
		// 	}))
		// 	if err != nil {
		// 		return nil, nil, err
		// 	}

		// 	return res.Msg.Grants, &res.Msg.NextPageToken, nil
		// })
		// if err != nil {
		// 	return err
		// }

		// if len(matchedGrants) == 0 {
		// 	return errors.New("no matching grants found")
		// }
		// if len(matchedGrants) != 1 {
		// 	clio.Debug("found more than one matching grant")
		// }

		// grant := matchedGrants[0]

		children, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*entityv1alpha1.Entity, *string, error) {
			res, err := grantsClient.QueryGrantChildren(ctx, connect.NewRequest(&accessv1alpha1.QueryGrantChildrenRequest{
				Id:        ensuredGrant.Grant.Id,
				PageToken: grab.Value(nextToken),
			}))
			if err != nil {
				return nil, nil, err
			}
			return res.Msg.Entities, &res.Msg.NextPageToken, nil
		})
		if err != nil {
			return err
		}

		commandData := CommandData{
			LocalPort: strconv.Itoa((c.Int("port"))),
		}

		for _, child := range children {
			switch child.Eid.Type {
			case "CF::Integration::AWSRDS::Grant::BastionInstance":
				err = entity.Unmarshal(child, &commandData.BastionInstance)
				if err != nil {
					return err
				}
			case "CF::Integration::AWSRDS::Grant::AccountAssignment":
				err = entity.Unmarshal(child, &commandData.AccountAssignment)
				if err != nil {
					return err
				}
			default:
				clio.Debugf("found unexpected child entity type %s", child.Eid.Type)
			}
		}

		if commandData.BastionInstance.ID == "" {
			return errors.New("did not find a bastion host in query grant children response")
		}
		if commandData.AccountAssignment.ID == "" {
			return errors.New("did not find an account assignment in query grant children response")
		}

		sso := cfaws.AwsSsoAssumer{}
		profile := &cfaws.Profile{
			Name:        commandData.AccountAssignment.Grant.ID,
			ProfileType: sso.Type(),
			AWSConfig: awsConfig.SharedConfig{
				SSOAccountID: commandData.AccountAssignment.Account.ID,
				SSORoleName:  commandData.AccountAssignment.Grant.ID,
				SSORegion:    commandData.AccountAssignment.IDCRegion,
				SSOStartURL:  commandData.AccountAssignment.StartURL,
			},
			Initialised: true,
		}

		creds, err := profile.SSOLogin(ctx, cfaws.ConfigOpts{})
		if err != nil {
			return err
		}

		clio.Infof("starting database proxy on port %v", commandData.LocalPort)
		cmd := exec.Command("aws", formatSSMCommandArgs(commandData)...)
		clio.Debugw("running aws ssm command", "command", "assume "+strings.Join(formatSSMCommandArgs(commandData), " "))
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, assume.EnvKeys(creds, commandData.BastionInstance.Region)...)

		// Start the command in a separate goroutine
		err = cmd.Start()
		if err != nil {
			return err
		}

		// Set up a channel to receive OS signals
		sigs := make(chan os.Signal, 1)
		// Notify sigs on os.Interrupt (Ctrl+C)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		// Wait for a termination signal in a separate goroutine
		go func() {
			<-sigs
			clio.Info("Received interrupt signal, shutting down...")
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				clio.Error("Error sending SIGTERM to process:", err)
			}
		}()

		// Wait for the command to finish
		err = cmd.Wait()
		if err != nil {
			clio.Error("Proxy connection failed with error:", err)
		} else {
			clio.Info("Proxy connection closed successfully")
		}
		return nil
	},
}

type CommandData struct {
	BastionInstance   BastionInstance
	AccountAssignment AccountAssignment
	LocalPort         string
}

func formatSSMCommandArgs(data CommandData) []string {
	out := []string{
		"ssm",
		"start-session",
		fmt.Sprintf("--target=%s", data.BastionInstance.ID),
		"--document-name=AWS-StartPortForwardingSession",
		"--parameters",
		fmt.Sprintf(`{"portNumber":["5432"], "localPortNumber":["%s"]}`, data.LocalPort),
	}

	return out
}
