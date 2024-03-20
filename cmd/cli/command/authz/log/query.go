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
		&cli.StringFlag{Name: "principal"},
		&cli.StringFlag{Name: "action"},
		&cli.StringFlag{Name: "resource"},
		&cli.StringSliceFlag{Name: "tag", Usage: "tags to filter for (e.g. --tag=read_only=false)"},
		&cli.StringFlag{Name: "output", Value: "table", Usage: "output format ('table' or 'json')"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		var tags []*authzv1alpha1.Tag

		for _, t := range c.StringSlice("tag") {
			split := strings.SplitN(t, "=", 2)
			if len(split) == 2 {
				tags = append(tags, &authzv1alpha1.Tag{
					Key:   split[0],
					Value: split[1],
				})
			} else {
				return fmt.Errorf("invalid tag format: %s. tags should be provided as --tag=key=value, for example: '--tag=read_only=false'", t)
			}
		}

		var principal, action, resource *entityv1alpha1.EID

		principalArg := c.String("principal")
		if principalArg != "" {
			parsed, err := eid.Parse(principalArg)
			if err != nil {
				return errors.Wrap(err, "parsing principal")
			}
			principal = parsed.ToAPI()
		}

		actionArg := c.String("action")
		if actionArg != "" {
			parsed, err := eid.Parse(actionArg)
			if err != nil {
				return errors.Wrap(err, "parsing action")
			}
			action = parsed.ToAPI()
		}

		resourceArg := c.String("resource")
		if resourceArg != "" {
			parsed, err := eid.Parse(resourceArg)
			if err != nil {
				return errors.Wrap(err, "parsing resource")
			}
			resource = parsed.ToAPI()
		}

		client := logsdk.New(cfg).Evaluation()

		input := logv1alpha1.QueryEvaluationsRequest{
			Principal:    principal,
			Action:       action,
			Resource:     resource,
			MatchingTags: tags,
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
