package workflow

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
	Usage: "Delete an Access Workflow",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "workflow-id", Required: true, Usage: "the workflow ID to delete"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := configsvc.NewFromConfig(cfg)

		res, err := client.AccessWorkflow().DeleteAccessWorkflow(ctx, connect.NewRequest(&configv1alpha1.DeleteAccessWorkflowRequest{
			Id: c.String("workflow-id"),
		}))
		if err != nil {
			return err
		}

		fmt.Printf("deleted %s\n", res.Msg.Id)

		return nil
	},
}
