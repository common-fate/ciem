package authz

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
)

type DecisionNotFoundError struct {
	Key string
}

func (e DecisionNotFoundError) Error() string {
	return fmt.Sprintf("decision not found for client_id: %s", e.Key)
}

type BatchAuthorizeInput struct {
	Requests map[string]Request
}

type BatchAuthorizeResponse struct {
	// Decisions in the response.
	// Usually you'll want GetDecision() to look up a specific decision
	// which checks it exists.
	Decisions map[string]*authzv1alpha1.Decision
}

// GetDecision looks up a decision by the client key. The client key is the key to the map
// that you provided in BatchAuthorizeInput.
func (r BatchAuthorizeResponse) GetDecision(key string) (*authzv1alpha1.Decision, error) {
	got, ok := r.Decisions[key]
	if !ok {
		return nil, DecisionNotFoundError{Key: key}
	}

	if got == nil {
		return nil, fmt.Errorf("decision %s was nil", key)
	}

	return got, nil
}

func (c *Client) BatchAuthorize(ctx context.Context, input BatchAuthorizeInput) (BatchAuthorizeResponse, error) {
	r := make([]*authzv1alpha1.AuthorizationRequest, 0)

	for _, item := range input.Requests {
		r = append(r, item.ToAPI())
	}

	res, err := c.raw.BatchAuthorize(ctx, connect.NewRequest(&authzv1alpha1.BatchAuthorizeRequest{
		Universe:    "default",
		Environment: "production",
		Requests:    r,
	}))
	if err != nil {
		return BatchAuthorizeResponse{}, err
	}

	ret := BatchAuthorizeResponse{
		Decisions: map[string]*authzv1alpha1.Decision{},
	}

	for _, d := range res.Msg.Decisions {
		ret.Decisions[d.ClientKey] = d
	}

	return ret, nil
}
