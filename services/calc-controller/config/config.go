package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off calc-controller service
type Config struct {
	REST    RMQConfig
	RPCCalc RPCConfigCalc
}

// RMQConfig stores configs for RabbitMQ connection
type RMQConfig struct {
	Port string `envconfig:"CC_RMQ_PORT" default:"5672"`
}

// RPCConfigCalc stores configs for GRPC connecion
type RPCConfigCalc struct {
	Port string `envconfig:"CC_GRPC_PORT" default:"50051"`
}

// InitConfig reads config variables from env and init *Config value
func InitConfig() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
