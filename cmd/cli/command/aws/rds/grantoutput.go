package rds

import (
	"github.com/common-fate/sdk/eid"
)

type GrantOutput struct {
	Grant          eid.EID `json:"grant" authz:"grant"`
	Name           string  `json:"name" authz:"name"`
	SSOStartURL    string  `json:"sso_start_url" authz:"sso_start_url"`
	SSORoleName    string  `json:"sso_role_name" authz:"sso_role_name"`
	SSORegion      string  `json:"sso_region" authz:"sso_region"`
	AccountID      string  `json:"account_id" authz:"account_id"`
	InstanceID     string  `json:"instance_id" authz:"instance_id"`
	InstanceRegion string  `json:"instance_region" authz:"instance_region"`
}

func (e GrantOutput) Parents() []eid.EID { return []eid.EID{e.Grant} }

func (e GrantOutput) EID() eid.EID { return eid.New(GrantOutputType, e.Grant.ID) }

const GrantOutputType = "CF::GrantOutput::AWSRDS"
