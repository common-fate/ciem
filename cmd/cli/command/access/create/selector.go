package create

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/printdiags"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	configv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/config/v1alpha1"
	"github.com/common-fate/sdk/service/control/config/accessselector"
	"github.com/common-fate/sdk/uid"
	"github.com/urfave/cli/v2"
)

var selectorCommand = cli.Command{
	Name: "selector",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Required: true},
		&cli.StringFlag{Name: "resource-type", Required: true},
		&cli.StringFlag{Name: "belonging-to", Required: true},
		&cli.StringFlag{Name: "when", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		belongingTo, err := uid.Parse(c.String("belonging-to"))
		if err != nil {
			return err
		}

		client := accessselector.NewFromConfig(cfg)

		res, err := client.CreateSelector(ctx, connect.NewRequest(&configv1alpha1.CreateSelectorRequest{
			Name:         c.String("name"),
			ResourceType: c.String("resource-type"),
			BelongingTo:  belongingTo.ToAPI(),
			When:         c.String("when"),
		}))

		clio.Debugw("result", "res", res)
		if err != nil {
			return err
		}

		printdiags.Print(res.Msg.Diagnostics, nil)

		fmt.Println(res.Msg.Selector.Id)

		return nil
	},
}
