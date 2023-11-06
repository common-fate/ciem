package device

import (
	"github.com/common-fate/ciem/config"
	"github.com/common-fate/ciem/gen/commonfate/cloud/attest/v1alpha1/attestv1alpha1connect"
)

func NewFromConfig(cfg *config.Context) attestv1alpha1connect.AttestServiceClient {
	return attestv1alpha1connect.NewAttestServiceClient(cfg.HTTPClient, cfg.APIURL)
}
