package main

import (
	"github.com/go-redis/redis"
	"os"
)

type (
	FromOptions     *redis.Options
	FromRedisClient *redis.Client
)

func NewFromOptions() (FromOptions, error) {
	return redis.ParseURL(os.Getenv("MICROREDIS_FROM_ADDRESS"))
}

func NewFromRedisClient(toOptions FromOptions) FromRedisClient {
	return redis.NewClient(toOptions)
}
