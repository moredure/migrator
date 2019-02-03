package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

type (
	FromOptions  *redis.Options
	ToOptions    *redis.Options
	RedisClients struct {
		From *redis.Client
		To   *redis.Client
	}
)

const (
	RedisSource      = "REDIS_SOURCE"
	RedisDestination = "REDIS_DESTINATION"
)

func NewRedisOptionsFromEnv(key string) (*redis.Options, error) {
	url, ok := os.LookupEnv(key)
	if !ok {
		return nil, fmt.Errorf("environment variable $%s is not defined", key)
	}
	return redis.ParseURL(url)
}

func NewFromOptions() (FromOptions, error) {
	return NewRedisOptionsFromEnv(RedisSource)
}

func NewToOptions() (ToOptions, error) {
	return NewRedisOptionsFromEnv(RedisDestination)
}

func NewRedisClients(from FromOptions, to ToOptions) RedisClients {
	return RedisClients{
		From: redis.NewClient(from),
		To:   redis.NewClient(to),
	}
}
