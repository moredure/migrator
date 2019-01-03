package main

import (
	"github.com/go-redis/redis"
	"os"
)

type FromRedisClient *redis.Client
type FromOptions *redis.Options

func NewFromOptions() (FromOptions, error) {
	return redis.ParseURL(os.Getenv("MICROREDIS_FROM_ADDRESS"))
}

func NewFromRedisClient(toOptions FromOptions) (FromRedisClient) {
	return redis.NewClient(toOptions)
}