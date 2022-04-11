package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off calc-controller service
type Config struct {
	RMQConfig RabbitMQConfig
	WConfig   WorkersConfig
}

type WorkersConfig struct {
	CountPlus  int `envconfig:"TODO" default:"1"`
	CountMinus int `envconfig:"TODO" default:"1"`
	CountMulti int `envconfig:"TODO" default:"1"`
	CountDiv   int `envconfig:"TODO" default:"1"`
	CountMod   int `envconfig:"TODO" default:"1"`
}

// RabbitMQConfig stores configs for RabbitMQ connection
type RabbitMQConfig struct {
	Server string `envconfig:"RABBITMQ_SERVER" default:"localhost"`
	Port   string `envconfig:"RABBITMQ_PORT" default:"5672"`

	QNameResult string `envconfig:"QUERY_NAME_RESULT" default:"expr.result"`
	QNamePLus   string `envconfig:"QUERY_NAME_PLUS" default:"expr.plus"`
	QNameMinus  string `envconfig:"QUERY_NAME_MINUS" default:"expr.minus"`
	QNameMulti  string `envconfig:"QUERY_NAME_MULTI" default:"expr.multi"`
	QNameDiv    string `envconfig:"QUERY_NAME_DIV" default:"expr.div"`
	QNameMod    string `envconfig:"QUERY_NAME_MOD" default:"expr.mod"`
}

// InitConfig reads config variables from env and init *Config value
func InitConfig() (*Config, error) {
	var cfg = new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
