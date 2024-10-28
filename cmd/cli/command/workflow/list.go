package workflow

import (
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"github.com/common-fate/sdk/config"
	resourcev1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/resource/v1alpha1"
	"github.com/common-fate/sdk/gen/commonfate/control/resource/v1alpha1/resourcev1alpha1connect"
	"github.com/urfave/cli/v2"
)

var listCommand = cli.Command{
	Name: "list",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := config.LoadDefault(ctx)
		if err != nil {
			return err
		}

		client := newResourceClient(cfg)
		// the Common Fate API doesn't currently expose a ListAccessWorkflows method, so we use the
		// QueryResources API.
		res, err := client.QueryResources(ctx, connect.NewRequest(&resourcev1alpha1.QueryResourcesRequest{
			Type: "Access::Workflow",
		}))
		if err != nil {
			return err
		}

		workflows := allWorkflows{Workflows: []workflow{}}

		for _, w := range res.Msg.Resources {
			workflows.Workflows = append(workflows.Workflows, workflow{
				ID:   w.Eid.Id,
				Name: w.Name,
			})
		}

		workflowsJSON, err := json.MarshalIndent(workflows, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(workflowsJSON))
		return nil
	},
}

// allWorkflows is used to ensure the output of the `cf workflows list`
// command remains stable, even if the API that we call changes from
// QueryResources to something else in future (such as ListAccessWorkflows).
type allWorkflows struct {
	Workflows []workflow `json:"workflows"`
}

type workflow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newResourceClient(cfg *config.Context) resourcev1alpha1connect.ResourceServiceClient {
	return resourcev1alpha1connect.NewResourceServiceClient(cfg.HTTPClient, cfg.APIURL)
}
