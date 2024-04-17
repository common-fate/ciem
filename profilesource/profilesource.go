package profilesource

import (
	"context"

	"connectrpc.com/connect"
	"github.com/common-fate/awsconfigfile"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	awsv1alpha1 "github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1"
	"github.com/common-fate/sdk/gen/granted/registry/aws/v1alpha1/awsv1alpha1connect"
)

// Source reads available AWS SSO profiles from the Common Fate API.
// It implements the awsconfigfile.Source interface
type Source struct {
	SSORegion    string
	StartURL     string
	Client       awsv1alpha1connect.ProfileRegistryServiceClient
	DashboardURL string
	Entitlements []*accessv1alpha1.EntitlementInput
}

func (s Source) GetProfiles(ctx context.Context) ([]awsconfigfile.SSOProfile, error) {
	profiles := []awsconfigfile.SSOProfile{}
	// add all options to our profile map
	for _, ent := range s.Entitlements {

		//lookup target

		acc, err := s.Client.GetProfileForAccountAndRole(ctx, &connect.Request[awsv1alpha1.GetProfileForAccountAndRoleRequest]{})
		if err != nil {
			return nil, err
		}

		p := awsconfigfile.SSOProfile{
			AccountID:     acc.Msg.Profile.Name,
			AccountName:   acc.Msg.Profile.Name,
			RoleName:      ent.Role.String(),
			SSOStartURL:   s.StartURL,
			SSORegion:     s.SSORegion,
			GeneratedFrom: "commonfate",
			CommonFateURL: s.DashboardURL,
		}
		profiles = append(profiles, p)

	}

	return profiles, nil
}
