package main

import (
	"github.com/go-redis/redis"
	"os"
)

type (
	ToOptions     *redis.Options
	ToRedisClient *redis.Client
)

func NewToOptions() (ToOptions, error) {
	return redis.ParseURL(os.Getenv("MICROREDIS_TO_ADDRESS"))
}

func NewToRedisClient(toOptions ToOptions) ToRedisClient {
	return redis.NewClient(toOptions)
}
