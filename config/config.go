package config

type Config struct {
	ServiceName                 string `yaml:"service_name"`
	ServiceID                   string `yaml:"service_id"`
	SharedVMPlanID              string `yaml:"shared_vm_plan_id"`
	Host                        string `yaml:"host"`
	DefaultConfigPath           string `yaml:"redis_conf_path"`
	ProcessCheckIntervalSeconds int    `yaml:"process_check_interval"`
	StartRedisTimeoutSeconds    int    `yaml:"start_redis_timeout"`
	InstanceDataDirectory       string `yaml:"data_directory"`
	PidfileDirectory            string `yaml:"pidfile_directory"`
	InstanceLogDirectory        string `yaml:"log_directory"`
	ServiceInstanceLimit        int    `yaml:"service_instance_limit"`
	Description                 string `yaml:"description"`
	LongDescription             string `yaml:"long_description"`
	ProviderDisplayName         string `yaml:"provider_display_name"`
	DocumentationURL            string `yaml:"documentation_url"`
	SupportURL                  string `yaml:"support_url"`
	DisplayName                 string `yaml:"display_name"`
	IconImage                   string `yaml:"icon_image"`
}
