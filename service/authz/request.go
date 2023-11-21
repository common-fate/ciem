package authz

import authzv1alpha1 "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1"

type Request struct {
	Principal UID `json:"principal"`
	Action    UID `json:"action"`
	Resource  UID `json:"resource"`
}

func (r Request) ToAPI() *authzv1alpha1.AuthorizationRequest {
	return &authzv1alpha1.AuthorizationRequest{
		Principal: r.Principal.ToAPI(),
		Action:    r.Action.ToAPI(),
		Resource:  r.Resource.ToAPI(),
	}
}
