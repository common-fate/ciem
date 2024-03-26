package log

import (
	"fmt"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/common-fate/cli/table"
	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	authzv1alpha1 "github.com/common-fate/sdk/gen/commonfate/authz/v1alpha1"
	logv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/log/v1alpha1"
	entityv1alpha1 "github.com/common-fate/sdk/gen/commonfate/entity/v1alpha1"
	logsdk "github.com/common-fate/sdk/service/control/log"
)

var queryCommand = cli.Command{
	Name: "query",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "principal"},
		&cli.StringSliceFlag{Name: "action"},
		&cli.StringSliceFlag{Name: "resource"},
		&cli.StringFlag{Name: "page-token"},
		&cli.StringSliceFlag{Name: "tag", Usage: "tags to filter for (e.g. --tag=read_only=true)"},
		&cli.StringSliceFlag{Name: "not-tag", Usage: "tags to exclude events for (e.g. --not-tag=read_only=true)"},
		&cli.StringFlag{Name: "outcome", Usage: "filter for a particular authorization outcome (either 'allow' or 'deny')"},
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table' or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		var filters []*logv1alpha1.Filter

		for _, t := range c.StringSlice("tag") {
			split := strings.SplitN(t, "=", 2)
			if len(split) == 2 {
				filters = append(filters, &logv1alpha1.Filter{
					Filter: &logv1alpha1.Filter_Tag{
						Tag: &logv1alpha1.TagFilter{
							Key:        split[0],
							Value:      split[1],
							Comparison: logv1alpha1.BoolComparison_BOOL_COMPARISON_EQUAL,
						},
					},
				})
			} else {
				return fmt.Errorf("invalid tag format: %s. tags should be provided as --tag=key=value, for example: '--tag=read_only=true'", t)
			}
		}

		for _, t := range c.StringSlice("not-tag") {
			split := strings.SplitN(t, "=", 2)
			if len(split) == 2 {
				filters = append(filters, &logv1alpha1.Filter{
					Filter: &logv1alpha1.Filter_Tag{
						Tag: &logv1alpha1.TagFilter{
							Key:        split[0],
							Value:      split[1],
							Comparison: logv1alpha1.BoolComparison_TIME_COMPARISON_NOT_EQUAL,
						},
					},
				})
			} else {
				return fmt.Errorf("invalid tag format: %s. tags should be provided as --not-tag=key=value, for example: '--not-tag=read_only=true'", t)
			}
		}

		var principals, actions, resources []*entityv1alpha1.EID

		for _, e := range c.StringSlice("principal") {
			parsed, err := eid.Parse(e)
			if err != nil {
				return errors.Wrap(err, "parsing principal")
			}
			principals = append(principals, parsed.ToAPI())
		}
		if len(principals) > 0 {
			filters = append(filters, &logv1alpha1.Filter{
				Filter: &logv1alpha1.Filter_Principal{
					Principal: &logv1alpha1.EntityFilter{
						Ids: principals,
					},
				},
			})
		}

		for _, e := range c.StringSlice("action") {
			parsed, err := eid.Parse(e)
			if err != nil {
				return errors.Wrap(err, "parsing action")
			}
			actions = append(actions, parsed.ToAPI())
		}
		if len(actions) > 0 {
			filters = append(filters, &logv1alpha1.Filter{
				Filter: &logv1alpha1.Filter_Action{
					Action: &logv1alpha1.EntityFilter{
						Ids: actions,
					},
				},
			})
		}

		for _, e := range c.StringSlice("resource") {
			parsed, err := eid.Parse(e)
			if err != nil {
				return errors.Wrap(err, "parsing resource")
			}
			resources = append(resources, parsed.ToAPI())
		}
		if len(resources) > 0 {
			filters = append(filters, &logv1alpha1.Filter{
				Filter: &logv1alpha1.Filter_Resource{
					Resource: &logv1alpha1.EntityFilter{
						Ids: resources,
					},
				},
			})
		}

		outcome := c.String("outcome")

		switch outcome {
		case "deny", "denied":
			filters = append(filters, &logv1alpha1.Filter{
				Filter: &logv1alpha1.Filter_Decision{
					Decision: &logv1alpha1.DecisionFilter{
						Decision: authzv1alpha1.Decision_DECISION_DENY,
					},
				},
			})
		case "allow", "allowed":
			filters = append(filters, &logv1alpha1.Filter{
				Filter: &logv1alpha1.Filter_Decision{
					Decision: &logv1alpha1.DecisionFilter{
						Decision: authzv1alpha1.Decision_DECISION_ALLOW,
					},
				},
			})
		case "":

		default:
			return fmt.Errorf("invalid --outcome flag: %s. outcome filter must be either 'allow' or 'deny'", outcome)
		}

		client := logsdk.New(cfg).Evaluation()

		input := logv1alpha1.QueryEvaluationsRequest{
			PageToken: c.String("page-token"),
			Filters:   filters,
		}

		clio.Debugw("calling QueryEvaluations", "input", &input)

		res, err := client.QueryEvaluations(ctx, connect.NewRequest(&input))
		if err != nil {
			return err
		}

		output := c.String("output")
		switch output {
		case "table":
			w := table.New(os.Stdout)
			w.Columns("ID", "PRINCIPAL", "ACTION", "RESOURCE", "DECISION")

			for _, r := range res.Msg.Evaluations {
				w.Row(r.Id, r.Request.Principal.Display(), r.Request.Action.Display(), r.Request.Resource.Display(), r.Decision.String())
			}

			err = w.Flush()
			if err != nil {
				return err
			}

		case "json":
			resJSON, err := protojson.Marshal(res.Msg)
			if err != nil {
				return err
			}
			fmt.Println(string(resJSON))
		default:
			return errors.New("invalid --output flag, valid values are [json, table]")
		}

		return nil
	},
}
