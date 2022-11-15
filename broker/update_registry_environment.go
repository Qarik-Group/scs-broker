package broker

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/types"

	"github.com/starkandwayne/scs-broker/broker/utilities"
)

func (broker *SCSBroker) UpdateRegistryEnvironment(cfClient *ccv3.Client, app *ccv3.Application, info *ccv3.Info, kind string, instanceId string, rc *utilities.RegistryConfig, params map[string]string) error {

	var profiles []string
	for key, value := range params {
		_, _, err := cfClient.UpdateApplicationEnvironmentVariables(app.GUID, ccv3.EnvironmentVariables{
			key: *types.NewFilteredString(value),
		})

		if key == "SPRING_CLOUD_CONFIG_SERVER_GIT_URI" {
			profiles = append(profiles, "git")
		}

		if key == "SPRING_CLOUD_CONFIG_SERVER_VAULT_HOST" {
			profiles = append(profiles, "vault")
		}

		if key == "SPRING_CLOUD_CONFIG_SERVER_COMPOSIT" {
			profiles = append(profiles, "composit")
		}

		if key == "SPRING_CLOUD_CONFIG_SERVER_CREDHUB" {
			profiles = append(profiles, "credhub")
		}

		if err != nil {
			return err
		}
	}

	var profileString strings.Builder
	for index, profile := range profiles {
		profileString.WriteString(profile)

		if index < len(profiles)-1 {
			profileString.WriteString(", ")
		}
	}

	peers, err := json.Marshal(rc.Peers)
	if err != nil {
		return err
	}

	_, _, err = cfClient.UpdateApplicationEnvironmentVariables(app.GUID, ccv3.EnvironmentVariables{
		"JWK_SET_URI":            *types.NewFilteredString(fmt.Sprintf("%v/token_keys", info.UAA())),
		"SKIP_SSL_VALIDATION":    *types.NewFilteredString(strconv.FormatBool(broker.Config.CfConfig.SkipSslValidation)),
		"REQUIRED_AUDIENCE":      *types.NewFilteredString(fmt.Sprintf("%s.%v", kind, instanceId)),
		"SPRING_PROFILES_ACTIVE": *types.NewFilteredString(profileString.String()),
		"PEERS":                  *types.NewFilteredString(string(peers)),
		"PEERING_MODE":           *types.NewFilteredString(rc.Mode),
	})
	if err != nil {
		return err
	}

	return nil
}
