package broker

import (
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
)

func (broker *SCSBroker) createRegistrySpace(org ccv3.Organization, instanceID string) (ccv3.Space, error) {
	client, err := broker.GetClient()
	if err != nil {
		return ccv3.Space{}, err
	}

	space, _, err := client.CreateSpace(
		ccv3.Space{
			Name: "eureka-" + instanceID,
			Relationships: ccv3.Relationships{
				constant.RelationshipTypeOrganization: ccv3.Relationship{GUID: org.GUID},
			},
		},
	)

	if err != nil {
		return space, err
	}

	return space, nil
}
