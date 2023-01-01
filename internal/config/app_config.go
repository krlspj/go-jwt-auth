package config

type AppConfig struct {
	// token duration in minutes
	TokenLifetime    int
	DatabaseName     string
	ConfigCollection string
	UsersCollection  string
	Collections      []string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
