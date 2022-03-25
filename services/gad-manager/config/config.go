package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off gad-manager service
type Config struct {
	REST    RESTConfig
	RPCCalc RPCConfigCalc
}

// RESTConfig stores configs for REST API Gin Server
type RESTConfig struct {
	Port string `envconfig:"GM_REST_PORT" default:"8080"`
}

// RPCConfigCalc stores configs for GRPC connecion
type RPCConfigCalc struct {
	Port string `envconfig:"GM_GRPC_PORT" default:"50051"`
}

// InitConfig reads config variables from env and init *Config value
func InitConfig() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
