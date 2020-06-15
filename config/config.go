package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigEnvVarName string = "CONFIG_SERVER_BROKER_CONFIG"

type Config struct {
	Auth                Auth     `yaml:"broker_auth"`
	ServiceName         string   `yaml:"service_name"`
	ServiceID           string   `yaml:"service_id"`
	BasicPlanId         string   `yaml:"basic_plan_id"`
	BasicPlanName       string   `yaml:"basic_plan_name"`
	CfConfig            CfConfig `yaml:"cloud_foundry_config"`
	Description         string   `yaml:"description"`
	LongDescription     string   `yaml:"long_description"`
	ProviderDisplayName string   `yaml:"provider_display_name"`
	DocumentationURL    string   `yaml:"documentation_url"`
	SupportURL          string   `yaml:"support_url"`
	DisplayName         string   `yaml:"display_name"`
	IconImage           string   `yaml:"icon_image"`
	InstanceSpaceGUID   string   `yaml:"instance_space_guid"`
	InstanceDomain      string   `yaml:"instance_domain"`
}

type Auth struct {
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

type CfConfig struct {
	ApiUrl   string `yaml:"api_url"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

func ParseConfig() (Config, error) {
	configJson := os.Getenv(ConfigEnvVarName)
	if configJson == "" {
		panic(ConfigEnvVarName + " not set")
	}
	var config Config

	return config, yaml.Unmarshal([]byte(configJson), &config)
}
