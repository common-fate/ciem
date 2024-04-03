package authz

import (
	"fmt"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/config"
	"github.com/common-fate/sdk/eid"
	"github.com/common-fate/sdk/service/authz"
	"github.com/common-fate/sdk/service/authz/batchauthz"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var evaluateCommand = cli.Command{
	Name:    "evaluate",
	Aliases: []string{"eval"},
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "principal", Aliases: []string{"p"}, Required: true},
		&cli.StringFlag{Name: "action", Aliases: []string{"a"}, Required: true},
		&cli.StringFlag{Name: "resource", Aliases: []string{"r"}, Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		principal, err := eid.Parse(c.String("principal"))
		if err != nil {
			return errors.Wrap(err, "parsing principal")
		}

		resource, err := eid.Parse(c.String("resource"))
		if err != nil {
			return errors.Wrap(err, "parsing resource")
		}

		action, err := eid.Parse(c.String("action"))
		if err != nil {
			return errors.Wrap(err, "parsing action")
		}

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := authz.NewFromConfig(cfg)

		batch := batchauthz.New(client.RawClient())

		req := authz.Request{
			Principal: principal,
			Action:    action,
			Resource:  resource,
		}
		err = batch.AddRequest(req)
		if err != nil {
			return err
		}

		err = batch.Authorize(ctx)
		if err != nil {
			return err
		}

		eval, err := batch.IsPermitted(req)
		if err != nil {
			return err
		}

		for key, annotations := range eval.Annotations.All() {
			for _, anno := range annotations {
				clio.Infof("%s: %s (%s)", key, anno.Value, anno.PolicyID)
			}
		}

		if !eval.Allowed {
			return fmt.Errorf("denied: %s", eval.ID)
		}

		clio.Successf("allowed: %s", eval.ID)

		return nil
	},
}
