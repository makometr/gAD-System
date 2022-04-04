package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off calc-controller service
type Config struct {
	RMQConfig RabbitMQConfig
}

type WorkersConfig struct {
	CountPlus  int `envconfig:"TODO" default:"1"`
	CountMinus int `envconfig:"TODO" default:"1"`
	CountMulti int `envconfig:"TODO" default:"1"`
	CountMod   int `envconfig:"TODO" default:"1"`
}

// RabbitMQConfig stores configs for RabbitMQ connection
type RabbitMQConfig struct {
	Server       string `envconfig:"RABBITMQ_SERVER" default:"localhost"`
	Port         string `envconfig:"RABBITMQ_PORT" default:"5672"`
	PubQueryName string `envconfig:"PUBLISH_QUERY_NAME" default:"cc-in"`
	SubQueryName string `envconfig:"SUBSCRIBE_QUERY_NAME" default:"cc-out"`

	QNamePLus  string
	QNameMinus string
	QNameMulti string
	QNameMod   string
}

// InitConfig reads config variables from env and init *Config value
func InitConfig() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
