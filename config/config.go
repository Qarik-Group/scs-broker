package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigEnvVarName string = "CONFIG_SERVER_BROKER_CONFIG"

type Config struct {
	ServiceName          string `yaml:"service_name"`
	ServiceID            string `yaml:"service_id"`
	SharedVMPlanID       string `yaml:"shared_vm_plan_id"`
	Host                 string `yaml:"host"`
	ServiceInstanceLimit int    `yaml:"service_instance_limit"`
	Description          string `yaml:"description"`
	LongDescription      string `yaml:"long_description"`
	ProviderDisplayName  string `yaml:"provider_display_name"`
	DocumentationURL     string `yaml:"documentation_url"`
	SupportURL           string `yaml:"support_url"`
	DisplayName          string `yaml:"display_name"`
	IconImage            string `yaml:"icon_image"`
}

func ParseConfig() (Config, error) {
	configJson := os.Getenv(ConfigEnvVarName)
	if configJson == "" {
		panic(ConfigEnvVarName + " not set")
	}
	var config Config

	return config, yaml.Unmarshal([]byte(configJson), &config)
}
