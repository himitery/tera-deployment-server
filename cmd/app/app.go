package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"tera/deployment/internal/adapters/kafka"
	"tera/deployment/internal/domain/services"
	"tera/deployment/internal/usecases"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
)

func NewApp(conf *config.Config) *fx.App {
	return fx.New(
		fx.Provide(
			func() *config.Config { return conf },
			logger.Init,

			// adapters
			kafka.NewKafkaConsumer,

			// services
			services.NewEventProcessor,
		),
		fx.Invoke(
			registerHooks,
		),
	)
}

func registerHooks(
	lc fx.Lifecycle,
	log *zap.Logger,
	processor usecases.EventProcessor,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			defer func() { _ = log.Sync() }()

			if err := processor.Register(); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := processor.Close(); err != nil {
				return err
			}

			return nil
		},
	})
}
