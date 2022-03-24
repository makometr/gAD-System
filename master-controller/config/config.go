package config

type Config struct {
	REST RESTConfig
}

type RESTConfig struct {
	PortREST string
}

func InitConfig() *Config {
	return &Config{REST: RESTConfig{PortREST: ":8080"}}
}
