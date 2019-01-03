// +build wireinject

package main

import (
	"github.com/google/wire"
)

func initializeApp() (*Migrator, error) {
	wire.Build(
		NewFromOptions,
		NewToOptions,
		NewToRedisClient,
		NewFromRedisClient,
		NewMigrator,
	)
	return &Migrator{}, nil
}
