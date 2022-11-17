package config

var ServerConfig *Config

func GetConfig() *Config {
	return ServerConfig
}

type Config struct {
	ServiceName string `json:"service_name" yaml:"serviceName"`
	Enable      bool   `json:"enable" yaml:"enable"`
}
