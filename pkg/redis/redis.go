package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
}

type Client struct {
	config Config
	client *redis.Client
}

func New(config Config) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Password: config.Password,
		Username: config.User,
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return redisClient, nil
}
