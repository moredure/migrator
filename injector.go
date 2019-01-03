// +build wireinject

package main

import (
	"github.com/google/wire"
)

func initializeApp() (*Migrator, error) {
	panic(wire.Build(
		NewFromOptions,
		NewToOptions,
		NewToRedisClient,
		NewFromRedisClient,
		NewMigrator,
	))
}
