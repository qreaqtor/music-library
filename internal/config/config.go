package config

type Config struct {
	Api      ApiConfig
	Postgres PostgresConfig

	Host string `env:"APP_HOST" env-required:"true"`
	Port int    `env:"APP_PORT" env-required:"true"`
	Env  string `env:"APP_ENV" env-required:"true"`
}

type ApiConfig struct {
	Version int `env:"API_VERSION" env-required:"true"`
}

type PostgresConfig struct {
	User     string `env:"POSTGRES_USER" env-required:"true"`
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DB       string `env:"POSTGRES_DB" env-required:"true"`

	Host string `env:"POSTGRES_HOST" env-required:"true"`
	Port int    `env:"POSTGRES_PORT" env-required:"true"`

	SSL bool `env:"POSTGRES_SSL" env-required:"true"`
}
