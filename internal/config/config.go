package config

import "github.com/caarlos0/env/v11"

type Config struct {
	App      app
	Http     http
	Postgres postgres
}

type app struct {
	Debug         bool   `env:"DEBUG"`
	TemplatesPath string `env:"TEMPLATES_PATH"`
}

type http struct {
	Host string `env:"HTTP_HOST"`
	Port uint16 `env:"HTTP_PORT"`
}

type postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     uint16 `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DbName   string `env:"POSTGRES_DBNAME"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
