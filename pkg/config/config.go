package config

import "github.com/joho/godotenv"

type Config struct {
	App    AppConfig
	Server ServerConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name: optional("APP_NAME", "gosocialize-api"),
			Env:  optional("APP_ENV", "development"),
		},
		Server: ServerConfig{
			Port: required("PORT"),
		},
	}
}
