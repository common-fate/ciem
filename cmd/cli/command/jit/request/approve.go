package request

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/request"
	"github.com/urfave/cli/v2"
)

var approveCommand = cli.Command{
	Name:  "approve",
	Usage: "Approve an Access Request",
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

		haserrors := printdiags.Print(res.Msg.Diagnostics, nil)
		if !haserrors {
			clio.Successf("approved request")
		}

		return nil

	},
}
