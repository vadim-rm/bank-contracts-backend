package config

import "github.com/caarlos0/env/v11"

type Config struct {
	App      app
	Http     http
	Postgres postgres
	Minio    minio
}

type app struct {
	Debug bool `env:"DEBUG"`
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

type minio struct {
	Endpoint   string `env:"MINIO_ENDPOINT"`
	Id         string `env:"MINIO_ID"`
	Secret     string `env:"MINIO_SECRET"`
	BucketName string `env:"MINIO_BUCKET_NAME"`
	BaseUrl    string `env:"MINIO_BASE_URL"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
