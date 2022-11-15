package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigEnvVarName string = "SCS_BROKER_CONFIG"

type Config struct {
	Auth                *Auth               `yaml:"broker_auth"`
	BrokerName          string              `yaml:"broker_name"`
	BrokerID            string              `yaml:"broker_id"`
	ArtifactsDir        string              `yaml:"artifacts_directory"`
	CfConfig            *CfConfig           `yaml:"cloud_foundry_config"`
	Description         string              `yaml:"description"`
	LongDescription     string              `yaml:"long_description"`
	ProviderDisplayName string              `yaml:"provider_display_name"`
	DocumentationURL    string              `yaml:"documentation_url"`
	SupportURL          string              `yaml:"support_url"`
	DisplayName         string              `yaml:"display_name"`
	IconImage           string              `yaml:"icon_image"`
	InstanceSpaceGUID   string              `yaml:"instance_space_guid"`
	InstanceDomain      string              `yaml:"instance_domain"`
	Services            map[string]*Service `yaml:"services"`
}

type Auth struct {
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

type Service struct {
	ServiceName             string `yaml:"service_name"`
	ServiceID               string `yaml:"service_id"`
	ServicePlanID           string `yaml:"service_plan_id"`
	ServicePlanName         string `yaml:"service_plan_name"`
	ServiceDescription      string `yaml:"service_description"`
	ServiceDownloadURI      string `yaml:"service_download_uri"`
	ServiceOrganizationGUID string `yaml:"service_organization_guid"`
	ServiceImplementation   string `yaml:"service_implementation"`
}

type CfConfig struct {
	ApiUrl            string `yaml:"api_url"`
	SkipSslValidation bool   `yaml:"skip_ssl_validation"`
	CfUsername        string `yaml:"cf_username"`
	CfPassword        string `yaml:"cf_password"`
	UaaClientID       string `yaml:"uaa_client_id"`
	UaaClientSecret   string `yaml:"uaa_client_secret"`
}

var Parsed *Config = &Config{Services: make(map[string]*Service)}

func ParseConfig() error {
	configJson := os.Getenv(ConfigEnvVarName)
	if configJson == "" {
		return fmt.Errorf("%s not set", ConfigEnvVarName)
	}

	return yaml.Unmarshal([]byte(configJson), Parsed)
}

func GetServiceByServiceID(serviceID string) (*Service, error) {
	for _, service := range Parsed.Services {
		if service.ServiceID == serviceID {
			return service, nil
		}
	}

	return nil, fmt.Errorf("no valid service found for %s", serviceID)
}
