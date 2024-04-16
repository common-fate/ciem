package profilesource

import (
	"context"

	"connectrpc.com/connect"
	"github.com/common-fate/awsconfigfile"
	accessv1alpha1 "github.com/common-fate/sdk/gen/commonfate/access/v1alpha1"
	accountv1alpha1 "github.com/common-fate/sdk/gen/commonfate/control/account"
	"github.com/common-fate/sdk/gen/commonfate/control/account/accountv1alpha1connect"
)

// Source reads available AWS SSO profiles from the Common Fate API.
// It implements the awsconfigfile.Source interface
type Source struct {
	SSORegion    string
	StartURL     string
	Client       accountv1alpha1connect.AccountServiceClient
	DashboardURL string
	Entitlements []*accessv1alpha1.EntitlementInput
}

func (s Source) GetProfiles(ctx context.Context) ([]awsconfigfile.SSOProfile, error) {
	profiles := []awsconfigfile.SSOProfile{}
	// add all options to our profile map
	for _, ent := range s.Entitlements {

		//lookup target

		acc, err := s.Client.GetAWSAccountDetail(ctx, &connect.Request[accountv1alpha1.GetAWSAccountDetailRequest]{
			Msg: &accountv1alpha1.GetAWSAccountDetailRequest{
				AccountId: ent.Target.GetEid().GetId(),
			},
		})

		if err != nil {
			return nil, err
		}

		p := awsconfigfile.SSOProfile{
			AccountID:     acc.Msg.AccountId,
			AccountName:   acc.Msg.AccountName,
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
