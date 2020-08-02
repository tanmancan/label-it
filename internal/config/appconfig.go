package config

// BaseConfig application configuration type
type baseConfig struct {
	ConfigVersion int
	BuildVersion  string
}

// AppConfig application configuration
var AppConfig = baseConfig{
	BuildVersion:  "1.0.0",
	ConfigVersion: 1,
}
