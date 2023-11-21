package authz

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"
)

// EntityJSON is a JSON representation of entities.
// It matches the Cedar Rust SDK JSON format.
type EntityJSON struct {
	UID     UID            `json:"uid"`
	Attrs   map[string]any `json:"attrs"`
	Parents []UID          `json:"parents"`
}

type BatchPutEntityJSONInput struct {
	Entities []EntityJSON
}

func (c *Client) BatchPutEntityJSON(ctx context.Context, input BatchPutEntityJSONInput) (*authzv1alpha1.BatchPutEntityResponse, error) {
	var req = &authzv1alpha1.BatchPutEntityRequest{
		Universe: "default",
		Entities: []*authzv1alpha1.Entity{},
	}

	for _, e := range input.Entities {
		parsed, err := transformJSONToEntity(e)
		if err != nil {
			return nil, err
		}

		req.Entities = append(req.Entities, parsed)
	}

	res, err := c.raw.BatchPutEntity(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return res.Msg, nil
}

func transformJSONToEntity(e EntityJSON) (*authzv1alpha1.Entity, error) {
	res := authzv1alpha1.Entity{
		Uid:        e.UID.ToAPI(),
		Attributes: []*authzv1alpha1.Attribute{},
		Parents:    []*authzv1alpha1.UID{},
	}

	for k, v := range e.Attrs {
		switch val := v.(type) {
		case string:
			res.Attributes = append(res.Attributes, &authzv1alpha1.Attribute{
				Key: k,
				Value: &authzv1alpha1.Value{
					Value: &authzv1alpha1.Value_Str{
						Str: val,
					},
				},
			})

		case int:
			res.Attributes = append(res.Attributes, &authzv1alpha1.Attribute{
				Key: k,
				Value: &authzv1alpha1.Value{
					Value: &authzv1alpha1.Value_Long{
						Long: int64(val),
					},
				},
			})
		default:
			return nil, fmt.Errorf("unhandled attribute type: %s", k)
		}
	}

	for _, v := range e.Parents {
		res.Parents = append(res.Parents, v.ToAPI())
	}

	return &res, nil
}
