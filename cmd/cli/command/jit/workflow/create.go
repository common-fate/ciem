package workflow

import (
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/common-fate/ciem/config"
	configv1alpha1 "github.com/common-fate/ciem/gen/commonfate/control/config/v1alpha1"
	"github.com/common-fate/ciem/service/control/config/jitworkflow"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/types/known/durationpb"
)

var createCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Required: true},
		&cli.Int64Flag{Name: "priority", Required: true},
		&cli.DurationFlag{Name: "access-duration", Required: true},
		&cli.StringSliceFlag{Name: "match-aws-account-id"},
		&cli.StringSliceFlag{Name: "aws-role"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := jitworkflow.NewFromConfig(cfg)

		res, err := client.CreateJITWorkflow(ctx, connect.NewRequest(&configv1alpha1.CreateJITWorkflowRequest{
			Name:           c.String("name"),
			Priority:       c.Int64("priority"),
			AccessDuration: durationpb.New(c.Duration("access-duration")),
			Filters: []*configv1alpha1.Filter{
				{
					Filter: &configv1alpha1.Filter_AwsAccount{
						AwsAccount: &configv1alpha1.AWSAccountFilter{
							MatchAccountIds: c.StringSlice("match-aws-account-id"),
							Role:            c.String("aws-role"),
						},
					},
				},
			},
		}))
		if err != nil {
			return err
		}

		fmt.Println(res.Msg.Workflow.Id)

		return nil
	},
}
