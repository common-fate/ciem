package deployment

import (
	"errors"
	"fmt"
	"slices"

	"connectrpc.com/connect"
	"github.com/common-fate/clio"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	diagnosticv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/diagnostic/v1alpha1"
	resetv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/integration/reset/v1alpha1"
	"github.com/common-fate/sdk/service/control/diagnostic"
	"github.com/common-fate/sdk/service/control/integration/reset"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var backgroundJobsCommand = cli.Command{
	Name:        "background-jobs",
	Usage:       "Manage background jobs",
	Subcommands: []*cli.Command{&cancelJobCommand, &batchCancelJobCommand, &retryJobCommand, &batchRetryJobCommand},
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

var retryJobCommand = cli.Command{
	Name:  "retry",
	Usage: "Retry a background job",
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

		_, err = client.RetryBackgroundJob(ctx, connect.NewRequest(&resetv1alpha1.RetryBackgroundJobRequest{
			Id: c.Int64("id"),
		}))
		if err != nil {
			return err
		}

		clio.Success("retried job successfully")

		return nil
	},
}

var batchCancelJobCommand = cli.Command{
	Name:  "cancel-batch",
	Usage: "Cancel background jobs",
	Flags: []cli.Flag{
		&cli.Int64SliceFlag{Name: "ids"},
		&cli.StringSliceFlag{Name: "kinds"},
		&cli.StringFlag{Name: "state", Usage: "valid states are ['available','retryable','running','scheduled']"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := reset.NewFromConfig(cfg)

		ids := c.Int64Slice("ids")
		kinds := c.StringSlice("kinds")
		if len(kinds) > 0 && c.String("state") == "" {
			return errors.New("--state is required when specifiying kinds")
		}
		if c.String("state") != "" {
			state, err := JobStateFromString(c.String("state"))
			if err != nil {
				return err
			}
			diagClient := diagnostic.NewFromConfig(cfg)
			backgroundJobs, err := diagClient.ListBackgroundJobs(ctx, connect.NewRequest(&diagnosticv1alpha1.ListBackgroundJobsRequest{
				Kinds: kinds,
				Count: grab.Ptr(int64(10000)),
				State: state,
			}))
			if err != nil {
				return err
			}
			for _, job := range backgroundJobs.Msg.Jobs {
				ids = append(ids, job.Id)
			}
		}

		slices.Sort(ids)
		ids = slices.Compact(ids)

		clio.Infow("Cancelling jobs", "count", len(ids), "job_ids", ids)
		for _, id := range ids {
			clio.Successf("cancelling job: %v", id)
			_, err = client.CancelBackgroundJob(ctx, connect.NewRequest(&resetv1alpha1.CancelBackgroundJobRequest{
				Id: id,
			}))
			if err != nil {
				clio.Errorw(fmt.Sprintf("failed to cancel job: %v", id), zap.Error(err))
				continue
			}
			clio.Successf("cancelled job: %v successfully", id)
		}

		return nil
	},
}

var batchRetryJobCommand = cli.Command{
	Name:  "retry-batch",
	Usage: "Retry background jobs",
	Flags: []cli.Flag{
		&cli.Int64SliceFlag{Name: "ids"},
		&cli.StringSliceFlag{Name: "kinds"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := reset.NewFromConfig(cfg)

		ids := c.Int64Slice("ids")
		kinds := c.StringSlice("kinds")

		diagClient := diagnostic.NewFromConfig(cfg)
		backgroundJobs, err := diagClient.ListBackgroundJobs(ctx, connect.NewRequest(&diagnosticv1alpha1.ListBackgroundJobsRequest{
			Kinds: kinds,
			Count: grab.Ptr(int64(10000)),
			State: diagnosticv1alpha1.JobState_JOB_STATE_RETRYABLE,
		}))
		if err != nil {
			return err
		}
		for _, job := range backgroundJobs.Msg.Jobs {
			ids = append(ids, job.Id)
		}

		slices.Sort(ids)
		ids = slices.Compact(ids)

		clio.Infow("Retrying jobs", "count", len(ids), "job_ids", ids)
		for _, id := range ids {
			clio.Successf("Retrying job: %v", id)
			_, err = client.RetryBackgroundJob(ctx, connect.NewRequest(&resetv1alpha1.RetryBackgroundJobRequest{
				Id: id,
			}))
			if err != nil {
				clio.Errorw(fmt.Sprintf("failed to retry job: %v", id), zap.Error(err))
				continue
			}
			clio.Successf("retried job: %v successfully", id)
		}

		return nil
	},
}
