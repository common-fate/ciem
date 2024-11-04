package availabilityspec

import (
	"fmt"

	"connectrpc.com/connect"
	"github.com/common-fate/sdk/config"
	configv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/config/v1alpha1"
	"github.com/common-fate/sdk/service/control/configsvc"
	"github.com/urfave/cli/v2"
)

var deleteCommand = cli.Command{
	Name:  "delete",
	Usage: "Delete an Availability Spec",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true, Usage: "the Availability Spec ID to delete"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := configsvc.NewFromConfig(cfg)

		res, err := client.AvailabilitySpec().DeleteAvailabilitySpec(ctx, connect.NewRequest(&configv1alpha1.DeleteAvailabilitySpecRequest{
			Id: c.String("id"),
		}))
		if err != nil {
			return err
		}

		fmt.Printf("deleted %s\n", res.Msg.Id)

		return nil
	},
}
