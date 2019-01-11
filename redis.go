package main

import (
	"github.com/go-redis/redis"
	"os"
)

type (
	FromOptions *redis.Options
	ToOptions   *redis.Options
	RedisClients struct {
		From *redis.Client
		To   *redis.Client
	}
)

func NewFromOptions() (FromOptions, error) {
	return redis.ParseURL(os.Getenv("MICROREDIS_FROM_ADDRESS"))
}

func NewToOptions() (ToOptions, error) {
	return redis.ParseURL(os.Getenv("MICROREDIS_TO_ADDRESS"))
}

func NewRedisClients(from FromOptions, to ToOptions) RedisClients {
	return RedisClients{
		From: redis.NewClient(from),
		To:   redis.NewClient(to),
	}
}
