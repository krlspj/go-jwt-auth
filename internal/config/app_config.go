package config

type AppConfig struct {
	// token duration in minutes
	TokenLifetime    int
	Database         string
	ConfigCollection string
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
