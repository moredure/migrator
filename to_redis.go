package main

import (
	"github.com/go-redis/redis"
	"os"
)

type ToOptions *redis.Options

func NewToOptions() (ToOptions, error) {
	return redis.ParseURL(os.Getenv("MICROREDIS_TO_ADDRESS"))
}

type ToRedisClient *redis.Client

func NewToRedisClient(toOptions ToOptions) (ToRedisClient) {
	return redis.NewClient(toOptions)
}
