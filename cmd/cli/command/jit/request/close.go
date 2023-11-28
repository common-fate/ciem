package request

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/request"
	"github.com/urfave/cli/v2"
)

var closeCommand = cli.Command{
	Name:  "close",
	Usage: "Close an Access Request",
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

		res, err := client.ApproveAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.ApproveAccessRequestRequest{
			Id: id,
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}

		for _, g := range res.Msg.Warnings.GrantsPermissionDenied {
			clio.Warnf("could not close grant to %s with role %s: permission was denied", g.Target.Display(), g.Role.Display())
		}

		for _, g := range res.Msg.Warnings.GrantsInvalidStatus {
			clio.Warnf("could not close grant to %s with role %s: grant was in an invalid status (%s)", g.Target.Display(), g.Role.Display(), g.Status)
		}

		if res.Msg.Warnings.OK() {
			clio.Successf("closed request")
		} else {
			clio.Warnf("closed request with warnings")
		}

		return nil

	},
}
