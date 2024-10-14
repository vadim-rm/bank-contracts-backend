package config

import (
	"github.com/caarlos0/env/v11"
	"time"
)

type Config struct {
	App      app
	Http     http
	Postgres postgres
	Minio    minio
	Redis    redis
	Jwt      jwt
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

type redis struct {
	Host     string `env:"REDIS_HOST"`
	Port     uint16 `env:"REDIS_PORT"`
	User     string `env:"REDIS_USER"`
	Password string `env:"REDIS_PASSWORD"`
}

type jwt struct {
	Token     string        `env:"JWT_TOKEN"`
	ExpiresIn time.Duration `env:"JWT_EXPIRES_IN"`
	Issuer    string        `env:"JWT_ISSUER"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
