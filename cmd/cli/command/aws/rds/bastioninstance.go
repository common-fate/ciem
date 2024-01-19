package rds

import "github.com/common-fate/sdk/eid"

type BastionInstance struct {
	ID           string  `json:"id" authz:"id"`
	Name         string  `json:"name" authz:"name"`
	Grant        eid.EID `json:"grant" authz:"grant"`
	Region       string  `json:"region" authz:"region"`
	Account      eid.EID `json:"account" authz:"account"`
	Organization eid.EID `json:"organization" authz:"organization"`
}

func (e BastionInstance) EID() eid.EID { return eid.New(BastionInstanceType, e.ID) }

func (e BastionInstance) Parents() []eid.EID { return []eid.EID{e.Organization, e.Grant} }

const BastionInstanceType = "CF::Integration::AWSRDS::Grant::BastionInstance"
