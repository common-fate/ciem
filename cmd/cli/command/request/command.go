package request

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/cmd/cli/command/request/access"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/request"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "request",
	Usage: "Request access to an entitlement",
	Subcommands: []*cli.Command{
		&access.Command,
		&approveCommand,
		&closeCommand,
		&revoke,
		&list,
	},
}

var revoke = cli.Command{
	Name:  "revoke",
	Usage: "Revoke an Access Request",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := request.NewFromConfig(cfg)

		id := c.String("id")

		res, err := client.RevokeAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.RevokeAccessRequestRequest{
			Id: id,
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}

		for _, g := range res.Msg.Warnings.GrantsPermissionDenied {
			clio.Warnf("could not revoke grant to %s with role %s: permission was denied", g.Target.Display(), g.Role.Display())
		}

		for _, g := range res.Msg.Warnings.GrantsInvalidStatus {
			clio.Warnf("could not revoke grant to %s with role %s: grant was in an invalid status (%s)", g.Target.Display(), g.Role.Display(), g.Status)
		}

		if res.Msg.Warnings.OK() {
			clio.Successf("revoked request")
		} else {
			clio.Warnf("revoked request with warnings")
		}

		return nil
	},
}
