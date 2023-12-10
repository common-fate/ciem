package close

import (
	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	"github.com/common-fate/sdk/service/access/request"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:  "close",
	Usage: "Perform a 'close' action on resources such as Access Requests",
	Subcommands: []*cli.Command{
		&requestCommand,
	},
}

var requestCommand = cli.Command{
	Name:  "request",
	Usage: "Close an Access Request",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringSliceFlag{Name: "grant-id", Usage: "Close specific Grants on the Access Request"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := request.NewFromConfig(cfg)

		id := c.String("id")

		res, err := client.CloseAccessRequest(ctx, connect.NewRequest(&accessv1alpha1.CloseAccessRequestRequest{
			Id: id,
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}

		haserrors := printdiags.Print(res.Msg.Diagnostics, nil)
		if !haserrors {
			clio.Successf("closed request")
		}

		return nil
	},
}
