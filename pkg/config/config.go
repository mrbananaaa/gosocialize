package config

import "github.com/joho/godotenv"

type Config struct {
	App    AppConfig
	Server ServerConfig
	DB     DatabaseConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
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
		DB: DatabaseConfig{
			URL: required("DB_URL"),
		},
	}
}

func (c *Config) IsDev() bool {
	return c.App.Env == "development"
}
