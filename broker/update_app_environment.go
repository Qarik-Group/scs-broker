package broker

import (
	"fmt"
	"strconv"
	"strings"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/types"
)

// Updates the app enviornment variables for creating or updating an instance.
func (broker *SCSBroker) UpdateAppEnvironment(cfClient *ccv3.Client, app *ccv3.Application, info *ccv3.Info, kind string, instanceId string, jsonparams string, params map[string]string) error {

	var hostKeySetSSH bool = false
	var profiles []string
	envVarToSet := make(ccv3.EnvironmentVariables)
	for key, value := range params {

		envVarToSet[key] = *types.NewFilteredString(value)

		if key == "SPRING_CLOUD_CONFIG_SERVER_GIT_URI" {
			profiles = append(profiles, "git")
		}

		if key == "SPRING_CLOUD_CONFIG_SERVER_GIT_HOSTKEY" {
			hostKeySetSSH = true
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

	}

	var profileString strings.Builder
	for index, profile := range profiles {
		profileString.WriteString(profile)

		if index < len(profiles)-1 {
			profileString.WriteString(", ")
		}
	}

	envVarToSet["SPRING_CLOUD_CONFIG_SERVER_GIT_IGNORELOCALSSHSETTINGS"] = *types.NewFilteredString("true")

	if !hostKeySetSSH {
		envVarToSet["SPRING_CLOUD_CONFIG_SERVER_GIT_STRICTHOSTKEYCHECKING"] = *types.NewFilteredString("false")
	} else {
		envVarToSet["SPRING_CLOUD_CONFIG_SERVER_GIT_STRICTHOSTKEYCHECKING"] = *types.NewFilteredString("true")
	}

	envVarToSet["SPRING_APPLICATION_JSON"] = *types.NewFilteredString(jsonparams)
	envVarToSet["JWK_SET_URI"] = *types.NewFilteredString(fmt.Sprintf("%v/token_keys", info.UAA()))
	envVarToSet["SKIP_SSL_VALIDATION"] = *types.NewFilteredString(strconv.FormatBool(broker.Config.CfConfig.SkipSslValidation))
	envVarToSet["REQUIRED_AUDIENCE"] = *types.NewFilteredString(fmt.Sprintf("%s.%v", kind, instanceId))
	envVarToSet["SPRING_PROFILES_ACTIVE"] = *types.NewFilteredString(profileString.String())

	_, _, err := cfClient.UpdateApplicationEnvironmentVariables(app.GUID, envVarToSet)
	if err != nil {
		return err
	}

	return nil
}
