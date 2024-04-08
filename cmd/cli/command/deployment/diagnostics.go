package deployment

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"github.com/common-fate/cli/table"
	"github.com/common-fate/grab"
	"github.com/common-fate/sdk/config"
	diagnosticv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/diagnostic/v1alpha1"
	"github.com/common-fate/sdk/service/control/diagnostic"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

var diagnosticsCommand = cli.Command{
	Name:  "diagnostics",
	Usage: "Retrieve diagnostics about your deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text' or 'json')"},
	},
	Subcommands: []*cli.Command{&backgroundTasksCommand},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		outputFormat := c.String("output")

		if outputFormat != "text" && outputFormat != "json" {
			return errors.New("--output flag must be either 'text' or 'json'")
		}

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		// fetch each diagnostic through separate API calls, and then combine them here
		// into the main 'full' set of diagnostics.
		//
		// This has been implemented to avoid having a single massive 'AllDiagnostics' endpoint
		// which may be expensive to call - plus, if something is wrong the entire API may return with an error.
		var all diagnosticv1alpha1.AllDiagnostics

		client := diagnostic.NewFromConfig(cfg)

		tokenMetadata, err := client.GetOAuthTokenMetadata(ctx, connect.NewRequest(&diagnosticv1alpha1.GetOAuthTokenMetadataRequest{}))
		if err != nil {
			return err
		}

		all.OauthTokenMetadata = tokenMetadata.Msg

		switch outputFormat {
		case "text":
			fmt.Println("OAUTH TOKEN METADATA")
			tbl := table.New(os.Stdout)
			tbl.Columns("ID", "APPNAME", "EXPIRES")

			for _, t := range all.OauthTokenMetadata.Tokens {
				exp := "-"

				if !t.ExpiresAt.AsTime().IsZero() {
					exp = t.ExpiresAt.AsTime().Format(time.RFC3339)
				}

				tbl.Row(t.Id, t.AppName, exp)
			}

			err = tbl.Flush()
			if err != nil {
				return err
			}

		case "json":
			resJSON, err := protojson.Marshal(&all)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		}

		return nil
	},
}

func JobStateFromString(state string) (diagnosticv1alpha1.JobState, error) {
	switch state {
	case "available":
		return diagnosticv1alpha1.JobState_JOB_STATE_AVAILABLE, nil
	case "cancelled":
		return diagnosticv1alpha1.JobState_JOB_STATE_CANCELLED, nil
	case "completed":
		return diagnosticv1alpha1.JobState_JOB_STATE_COMPLETED, nil
	case "discarded":
		return diagnosticv1alpha1.JobState_JOB_STATE_DISCARDED, nil
	case "retryable":
		return diagnosticv1alpha1.JobState_JOB_STATE_RETRYABLE, nil
	case "running":
		return diagnosticv1alpha1.JobState_JOB_STATE_RUNNING, nil
	case "scheduled":
		return diagnosticv1alpha1.JobState_JOB_STATE_SCHEDULED, nil
	default:
		return diagnosticv1alpha1.JobState_JOB_STATE_UNSPECIFIED, nil
	}
}

type JobSummary struct {
	Id            string
	Name          string
	LastRun       time.Time
	NextRun       time.Time
	Errors        string
	CurrentStatus string
}

var backgroundTasksCommand = cli.Command{
	Name:  "background-jobs",
	Usage: "Retrieve diagnostics about your deployments background jobs",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "output", Value: "text", Usage: "output format ('text' or 'json')"},
		&cli.StringSliceFlag{Name: "kinds"},
		&cli.StringFlag{Name: "state", Usage: "valid states are ['available','cancelled','completed','discarded','retryable','running','scheduled']"},
		&cli.Int64Flag{Name: "count"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		outputFormat := c.String("output")

		if outputFormat != "text" && outputFormat != "json" {
			return errors.New("--output flag must be either 'text' or 'json'")
		}

		// if c.StringSlice("kinds") != nil && c.String("state") == "" {
		// 	return errors.New("flag --kinds must be used with --state flag")
		// }

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := diagnostic.NewFromConfig(cfg)

		if c.String("state") == "" && c.StringSlice("kinds") == nil {
			fmt.Println("Background Jobs")
			tbl := table.New(os.Stdout)
			tbl.Columns("ID", "KIND", "STATE", "LAST_STARTED_AT", "TIME_ELAPSED", "ERRORS")
			kinds, err := client.ListBackgroundJobKindSummary(ctx, connect.NewRequest(&diagnosticv1alpha1.ListBackgroundJobKindSummaryRequest{}))
			if err != nil {
				return err
			}

			//pulling the latest update from each kind of job
			for _, job := range kinds.Msg.Jobs {

				tbl.Row(job.Id, job.Kind, job.State, job.LastRun.AsTime().Local().String(), job.TimeElapsed.AsDuration().String(), job.Errors)
			}

			err = tbl.Flush()
			if err != nil {
				return err
			}
			return nil
		}

		state, err := JobStateFromString(c.String("state"))
		if err != nil {
			return err
		}

		backgroundJobs, err := client.ListBackgroundJobs(ctx, connect.NewRequest(&diagnosticv1alpha1.ListBackgroundJobsRequest{
			Kinds: c.StringSlice("kinds"),
			Count: grab.If(c.Int64("count") > 0, grab.Ptr(c.Int64("count")), grab.Ptr(int64(100))),
			State: grab.If(c.String("state") == "", nil, &state),
		}))
		if err != nil {
			return err
		}

		switch outputFormat {
		case "text":

			fmt.Println("Background Jobs")
			tbl := table.New(os.Stdout)
			tbl.Columns("ID", "KIND", "STATE", "OCCURED_AT")

			for _, job := range backgroundJobs.Msg.Jobs {
				tbl.Row(strconv.Itoa(int(job.Id)), job.Kind, job.State, job.CreatedAt.AsTime().String())
			}

			err = tbl.Flush()
			if err != nil {
				return err
			}

		case "json":
			resJSON, err := protojson.Marshal(backgroundJobs.Msg)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		}

		return nil
	},
}
