package config

type Config struct {
	App    AppConfig
	Server ServerConfig
	DB     DatabaseConfig
}

type AppConfig struct {
	Name string `env:"APP_NAME" envDefault:"gosocialize"`
	Env  string `env:"APP_ENV" envDefault:"development"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT,required"`
}

type DatabaseConfig struct {
	URL string `env:"DB_URL,required"`
}
