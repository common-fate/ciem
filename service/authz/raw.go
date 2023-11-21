package authz

import "github.com/common-fate/ciem/gen/commonfate/authz/v1alpha1/authzv1alpha1connect"

// RawClient returns the underlying Buf Connect client.
func (c *Client) RawClient() authzv1alpha1connect.AuthzServiceClient {
	return c.raw
}
