package broker

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"github.com/ess/hype"

	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) UpdateRegistryEnvironment(app *ccv3.Application, url string, rc *utilities.RegistryConfig) error {
	client, err := broker.GetClient()
	if err != nil {
		return err
	}

	routes, _, err := client.GetApplicationRoutes(app.GUID)
	if err != nil {
		return err
	}

	peers, err := json.Marshal(rc.Peers)
	if err != nil {
		return err
	}

	beast, err := hype.New(fmt.Sprintf("https://%s", routes[0].URL))
	if err != nil {
		return err
	}

	for _, peer := range rc.Peers {
		resp := beast.
			WithoutTLSVerification().
			Post("cf-config-peers", nil, peers).
			WithHeader(hype.NewHeader("X-Cf-App-Instance", fmt.Sprintf("%s:%d", app.GUID, peer.Index))).
			Response()

		if !resp.Okay() {
			return resp.Error()
		}
	}

	return nil
}
