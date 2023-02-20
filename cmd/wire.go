//go:build wireinject
// +build wireinject

package main

import (
	"stream/app/middleware"
	"stream/config"
	"stream/router"

	"github.com/google/wire"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func wireApp(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*App, func(), error) {
	panic(
		wire.Build(
			middleware.ProviderSet,
			router.ProviderSet,
			newHttpServer,
			newApp,
		),
	)
}