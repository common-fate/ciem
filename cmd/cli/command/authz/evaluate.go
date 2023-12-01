package authz

import (
	"fmt"

	"github.com/common-fate/clio"
	"github.com/common-fate/sdk/service/authz"
	"github.com/common-fate/sdk/service/authz/batchauthz"
	"github.com/common-fate/sdk/service/authz/uid"
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

		principal, err := uid.Parse(c.String("principal"))
		if err != nil {
			return errors.Wrap(err, "parsing principal")
		}

		resource, err := uid.Parse(c.String("resource"))
		if err != nil {
			return errors.Wrap(err, "parsing resource")
		}

		action, err := uid.Parse(c.String("action"))
		if err != nil {
			return errors.Wrap(err, "parsing action")
		}

		client := authz.NewClient(authz.Opts{
			HTTPClient: newInsecureClient(),
			BaseURL:    "http://127.0.0.1:5050",
		})

		batch := batchauthz.New(client.RawClient())

		req := authz.Request{
			Principal: principal,
			Action:    action,
			Resource:  resource,
		}
		batch.AddRequest(req)

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
