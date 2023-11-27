package access

import (
	"github.com/common-fate/ciem/config"
	"github.com/common-fate/ciem/gen/commonfate/control/config/v1alpha1/configv1alpha1connect"
)

func NewFromConfig(cfg *config.Context) configv1alpha1connect.ConfigServiceClient {
	return configv1alpha1connect.NewConfigServiceClient(cfg.HTTPClient, cfg.APIURL)
}
