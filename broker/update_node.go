package broker

import (
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/types"

	"github.com/starkandwayne/scs-broker/broker/result"
	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) updateNode(node *unupdatedNode, rc *utilities.RegistryConfig, pipeline chan<- *result.UpdateApp) {
	cfClient, err := broker.GetClient()
	if err != nil {
		pipeline <- result.NewUpdateApp().WithError(err)
		return
	}

	appJSON := rc.ForNode(node.URL)
	trusted := make([]string, 0)

	for _, peer := range rc.Peers {
		trusted = append(trusted, peer.Host)
	}

	_, _, err = cfClient.UpdateApplicationEnvironmentVariables(node.App.GUID, ccv3.EnvironmentVariables{
		"SKIP_SSL_VALIDATION": *types.NewFilteredString(strconv.FormatBool(broker.Config.CfConfig.SkipSslValidation)),
		//"REQUIRED_AUDIENCE":       *types.NewFilteredString(fmt.Sprintf("%s.%v", kind, instanceId)),
		//"SPRING_PROFILES_ACTIVE":  *types.NewFilteredString(profileString.String()),
		"SPRING_APPLICATION_JSON": *types.NewFilteredString(appJSON),
		"TRUST_CERTS":             *types.NewFilteredString(strings.Join(trusted, ",")),
	})
	if err != nil {
		pipeline <- result.NewUpdateApp().WithError(err)
		return
	}

	pipeline <- result.NewUpdateApp().WithApp(node.App)
}
