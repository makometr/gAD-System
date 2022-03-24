package config

type Config struct {
	REST    RESTConfig
	RPCCalc RPCCalculate
}

type RESTConfig struct {
	Port string
}

type RPCCalculate struct {
	Port string
}

func InitConfig() *Config {
	return &Config{REST: RESTConfig{Port: ":8080"}, RPCCalc: RPCCalculate{Port: "localhost:50051"}}
}
