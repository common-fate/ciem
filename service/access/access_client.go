package access

import (
	"github.com/common-fate/ciem/config"
	"github.com/common-fate/ciem/gen/proto/common_fate/v1alpha1/common_fatev1alpha1connect"
)

func NewFromConfig(cfg *config.Context) common_fatev1alpha1connect.AccessServiceClient {
	return common_fatev1alpha1connect.NewAccessServiceClient(cfg.HTTPClient, cfg.APIURL)
}
