package availabilityspec

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
		// the Common Fate API doesn't currently expose a ListAvailabilitySpecs method, so we use the
		// QueryResources API.
		res, err := client.QueryResources(ctx, connect.NewRequest(&resourcev1alpha1.QueryResourcesRequest{
			Type: "Access::AvailabilitySpec",
		}))
		if err != nil {
			return err
		}

		availabilitySpecs := allAvailabilitySpecs{AvailabilitySpecs: []availabilitySpec{}}

		for _, w := range res.Msg.Resources {
			availabilitySpecs.AvailabilitySpecs = append(availabilitySpecs.AvailabilitySpecs, availabilitySpec{
				ID:   w.Eid.Id,
				Name: w.Name,
			})
		}

		workflowsJSON, err := json.MarshalIndent(availabilitySpecs, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(workflowsJSON))
		return nil
	},
}

// allAvailabilitySpecs is used to ensure the output of the `cf workflows list`
// command remains stable, even if the API that we call changes from
// QueryResources to something else in future (such as ListAvailabilitySpecs).
type allAvailabilitySpecs struct {
	AvailabilitySpecs []availabilitySpec `json:"availabilitySpecs"`
}

type availabilitySpec struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newResourceClient(cfg *config.Context) resourcev1alpha1connect.ResourceServiceClient {
	return resourcev1alpha1connect.NewResourceServiceClient(cfg.HTTPClient, cfg.APIURL)
}
