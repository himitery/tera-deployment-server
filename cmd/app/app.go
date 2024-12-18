package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
)

func NewApp(conf *config.Config) *fx.App {
	return fx.New(
		fx.Provide(
			func() *config.Config { return conf },
			logger.Init,
		),
		fx.Invoke(
			registerHooks,
		),
	)
}

func registerHooks(lc fx.Lifecycle, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer func() { _ = log.Sync() }()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
