package services

import (
	"go.uber.org/zap"
	"tera/deployment/internal/domain/models"
	"tera/deployment/internal/ports"
	"tera/deployment/internal/usecases"
	"tera/deployment/pkg/logger"
)

type EventProcessor struct {
	consumer ports.KafkaConsumer
	events   chan *models.EventMessage
}

func NewEventProcessor(consumer ports.KafkaConsumer) usecases.EventProcessor {
	events := make(chan *models.EventMessage)

	return &EventProcessor{
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
}
