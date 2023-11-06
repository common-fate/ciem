package device

import (
	"github.com/common-fate/ciem/config"
	cfconnect "github.com/common-fate/ciem/gen/proto/commonfatecloud/v1alpha1/commonfatecloudv1alpha1connect"
)

func NewFromConfig(cfg *config.Context) cfconnect.DeviceServiceClient {
	return cfconnect.NewDeviceServiceClient(cfg.HTTPClient, cfg.APIURL)
}
