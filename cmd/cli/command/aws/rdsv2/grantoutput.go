package rdsv2

import (
	"github.com/common-fate/sdk/eid"
)

type GrantOutput struct {
	Grant          eid.EID `json:"grant" authz:"grant"`
	Name           string  `json:"name" authz:"name"`
	SSOStartURL    string  `json:"sso_start_url" authz:"sso_start_url"`
	SSORoleName    string  `json:"sso_role_name" authz:"sso_role_name"`
	SSORegion      string  `json:"sso_region" authz:"sso_region"`
	ProxyAccountID string  `json:"proxy_account_id" authz:"proxy_account_id"`
	ProxyRegion    string  `json:"proxy_region" authz:"proxy_region"`
}

func (e GrantOutput) Parents() []eid.EID { return []eid.EID{e.Grant} }

func (e GrantOutput) EID() eid.EID { return eid.New(GrantOutputType, e.Grant.ID) }

const GrantOutputType = "CF::GrantOutput::AWSRDSV2"
