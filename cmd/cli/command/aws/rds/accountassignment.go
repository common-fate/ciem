package rds

import (
	"github.com/common-fate/sdk/eid"
)

type AccountAssignment struct {
	ID            string  `json:"id" authz:"id"`
	Name          string  `json:"name" authz:"name"`
	Grant         eid.EID `json:"grant" authz:"grant"`
	StartURL      string  `json:"start_url" authz:"start_url"`
	IDCRegion     string  `json:"idc_region" authz:"idc_region"`
	Account       eid.EID `json:"account" authz:"account"`
	PermissionSet eid.EID `json:"permission_set" authz:"permission_set"`
	Organization  eid.EID `json:"organization" authz:"organization"`
}

func (e AccountAssignment) Parents() []eid.EID { return []eid.EID{e.Organization, e.Grant} }

func (e AccountAssignment) EID() eid.EID { return eid.New(AccountAssignmentType, e.ID) }

const AccountAssignmentType = "CF::Integration::AWSRDS::Grant::AccountAssignment"
