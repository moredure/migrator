package main

import (
	"github.com/go-redis/redis"
	"os"
	"fmt"
)

type (
	FromOptions *redis.Options
	ToOptions   *redis.Options
	RedisClients struct {
		From *redis.Client
		To   *redis.Client
	}
)

const (
	REDIS_SOURCE      = "REDIS_SOURCE"
	REDIS_DESTINATION = "REDIS_DESTINATION"
)

func NewRedisOptionsFromEnv(key string) (*redis.Options, error) {
	url, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("environment variable $%s is not defined", key)
	}
	return redis.ParseURL(url)
}

func NewFromOptions() (FromOptions, error) {
	return NewRedisOptionsFromEnv(REDIS_SOURCE)
}

func NewToOptions() (ToOptions, error) {
	return NewRedisOptionsFromEnv(REDIS_DESTINATION)
}

func NewRedisClients(from FromOptions, to ToOptions) RedisClients {
	return RedisClients{
		From: redis.NewClient(from),
		To:   redis.NewClient(to),
	}
}
