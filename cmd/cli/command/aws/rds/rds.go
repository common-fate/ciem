package rds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/bufbuild/connect-go"

	"github.com/common-fate/clio"
	"github.com/common-fate/grab"

	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	entityv1alpha1 "github.com/common-fate/sdk/gen/commonfate/entity/v1alpha1"
	"github.com/common-fate/sdk/service/access/grants"
	"github.com/common-fate/sdk/service/entity"
	"github.com/common-fate/sdk/service/identity"
	"github.com/urfave/cli/v2"
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
		&cli.IntFlag{Name: "port", Value: 5432, Usage: "The local port to forward the database connect to"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := grants.NewFromConfig(cfg)
		idClient := identity.NewFromConfig(cfg)

		idRes, err := idClient.GetCallerIdentity(ctx, connect.NewRequest(&accessv1alpha1.GetCallerIdentityRequest{}))
		if err != nil {
			return err
		}

		target := c.String("target")

		matchedGrants, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*accessv1alpha1.Grant, *string, error) {
			res, err := client.QueryGrants(ctx, connect.NewRequest(&accessv1alpha1.QueryGrantsRequest{
				PageToken: grab.Value(nextToken),
				Principal: idRes.Msg.Principal.Eid,
				Target: &entityv1alpha1.EID{
					Type: "AWS::RDS::Instance",
					Id:   target,
				},
				Role: &entityv1alpha1.EID{
					Type: "CF::Database::Role",
					Id:   "ReadWrite",
				},
				Status: accessv1alpha1.GrantStatus_GRANT_STATUS_ACTIVE.Enum(),
			}))
			if err != nil {
				return nil, nil, err
			}

			return res.Msg.Grants, &res.Msg.NextPageToken, nil
		})
		if err != nil {
			return err
		}

		if len(matchedGrants) == 0 {
			return errors.New("no matching grants found")
		}
		if len(matchedGrants) != 1 {
			clio.Debug("found more than one matching grant")
		}

		grant := matchedGrants[0]

		children, err := grab.AllPages(ctx, func(ctx context.Context, nextToken *string) ([]*entityv1alpha1.Entity, *string, error) {
			res, err := client.QueryGrantChildren(ctx, connect.NewRequest(&accessv1alpha1.QueryGrantChildrenRequest{
				Id:        grant.Id,
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

		clio.Infof("starting database proxy on port %v", commandData.LocalPort)
		cmd := exec.Command("assume", formatAssumeCommandArgs(commandData)...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "FORCE_NO_ALIAS=true")

		err = cmd.Run()
		if err != nil {
			return err
		}

		return nil
	},
}

type CommandData struct {
	BastionInstance   BastionInstance
	AccountAssignment AccountAssignment
	LocalPort         string
}

func formatAssumeCommandArgs(data CommandData) []string {
	out := []string{
		"--sso",
		fmt.Sprintf("--sso-start-url=%s", data.AccountAssignment.StartURL),
		fmt.Sprintf("--sso-region=%s", data.AccountAssignment.IDCRegion),
		fmt.Sprintf("--account-id=%s", data.AccountAssignment.Account.ID),
		fmt.Sprintf("--role-name=%s", data.AccountAssignment.Grant.ID),
		"--exec",
		"--",
		"aws",
		"ssm",
		"start-session",
		fmt.Sprintf("--target=%s", data.BastionInstance.ID),
		"--document-name=AWS-StartPortForwardingSession",
		"--parameters",
		fmt.Sprintf(`{"portNumber":["5432"], "localPortNumber":["%s"]}`, data.LocalPort),
	}

	return out
}
