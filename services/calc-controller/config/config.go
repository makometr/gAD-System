package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off calc-controller service
type Config struct {
	RMQConfig RabbitMQConfig
	CCConfig  CalcControllerConfig
}

// RabbitMQConfig stores configs for RabbitMQ connection
type RabbitMQConfig struct {
	Server       string `envconfig:"RABBITMQ_SERVER" default:"localhost"`
	Port         string `envconfig:"RABBITMQ_PORT" default:"5672"`
	PubQueryName string `envconfig:"PUBLISH_QUERY_NAME" default:"cc-in"`
	SubQueryName string `envconfig:"SUBSCRIBE_QUERY_NAME" default:"cc-out"`
}

// CalcControllerConfig stores configs for GRPC connection
type CalcControllerConfig struct {
	Server string `envconfig:"CALCULATION_CONTROLLER_SERVER" default:"localhost"`
	Port   string `envconfig:"CALCULATION_CONTROLLER_PORT" default:":50051"`
}

// InitConfig reads config variables from env and init *Config value
func InitConfig() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
