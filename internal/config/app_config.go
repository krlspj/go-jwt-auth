package config

type AppConfig struct {
	// token duration in minutes
	TokenLifetime int
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
