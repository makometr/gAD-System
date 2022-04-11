package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off gad-manager service
type Config struct {
	GMConfig GadManagerConfig
	CCConfig CalcControllerConfig
}

// GadManagerConfig stores configs for REST API Gin Server
type GadManagerConfig struct {
	Server string `envconfig:"GAD_MANAGER_SERVER" default:"localhost"`
	Port   string `envconfig:"GAD_MANAGER_PORT" default:"8080"`
}

// CalcControllerConfig stores configs for GRPC connection
type CalcControllerConfig struct {
	Server string `envconfig:"CALCULATION_CONTROLLER_SERVER" default:"localhost"`
	Port   string `envconfig:"CALCULATION_CONTROLLER_PORT" default:"50051"`
}

// InitConfig reads config variables from env and init *Config value
func InitConfig() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
