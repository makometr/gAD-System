package config

import "github.com/kelseyhightower/envconfig"

// Config stores all configs off calc-controller service
type Config struct {
	RMQConfig RabbitMQConfig
	CCConfig  CalcControllerConfig
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
