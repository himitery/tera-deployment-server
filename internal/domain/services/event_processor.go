package services

import (
	"go.uber.org/zap"
	"strings"
	"tera/deployment/internal/domain/models"
	"tera/deployment/internal/ports"
	"tera/deployment/internal/usecases"
	"tera/deployment/pkg/logger"
)

type EventProcessor struct {
	manager  usecases.DeploymentManager
	consumer ports.KafkaConsumer
	events   chan *models.EventMessage
}

func NewEventProcessor(
	manager usecases.DeploymentManager,
	consumer ports.KafkaConsumer,
) usecases.EventProcessor {
	events := make(chan *models.EventMessage)

	return &EventProcessor{
		manager:  manager,
		consumer: consumer,
		events:   events,
	}
}

func (ctx *EventProcessor) Register() error {
	if err := ctx.consumer.Start(ctx.events); err != nil {
		return err
	}

	go func() {
		for message := range ctx.events {
			ctx.process(message)
		}
	}()

	return nil
}

func (ctx *EventProcessor) Close() error {
	close(ctx.events)
	if err := ctx.consumer.Close(); err != nil {
		return err
	}

	return nil
}

func (ctx *EventProcessor) process(message *models.EventMessage) {
	logger.Info("new message received", zap.Any("message", message))

	if strings.ToLower(message.Action) == "create" {
		application, err := ctx.manager.Create(
			message.Service,
			message.Version,
			message.Namespace,
			message.Values,
		)
		if application != nil && err == nil {
			logger.Info("application created", zap.Any("application", application))
		}
	}
}
