package deployment

import (
	"connectrpc.com/connect"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	resetv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/reset/v1alpha1"
	"github.com/common-fate/sdk/service/control/integration/reset"
	"github.com/urfave/cli/v2"
)

var backgroundJobsCommand = cli.Command{
	Name:        "background-jobs",
	Usage:       "Manage background jobs",
	Subcommands: []*cli.Command{&cancelJobCommand},
}

var cancelJobCommand = cli.Command{
	Name:  "cancel",
	Usage: "Cancel a background job",
	Flags: []cli.Flag{
		&cli.Int64Flag{Name: "id", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := reset.NewFromConfig(cfg)

		_, err = client.CancelBackgroundJob(ctx, connect.NewRequest(&resetv1alpha1.CancelBackgroundJobRequest{
			Id: c.Int64("id"),
		}))
		if err != nil {
			return err
		}

		clio.Success("cancelled job successfully")

		return nil
	},
}
