package services

import (
	"fmt"
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
	producer ports.KafkaProducer
	events   chan any
}

func NewEventProcessor(
	events chan any,
	manager usecases.DeploymentManager,
	consumer ports.KafkaConsumer,
	producer ports.KafkaProducer,
) usecases.EventProcessor {
	return &EventProcessor{
		manager:  manager,
		consumer: consumer,
		producer: producer,
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
	if err := ctx.consumer.Close(); err != nil {
		return err
	}

	return nil
}

func (ctx *EventProcessor) process(data any) {
	logger.Info("new message received", zap.Any("message", data))

	switch message := data.(type) {
	case *models.KafkaMessage:
		ctx.processKafkaMessage(message)
	case *models.SystemMessage:
		ctx.processSystemMessage(message)
	}
}

func (ctx *EventProcessor) processKafkaMessage(message *models.KafkaMessage) {
	switch strings.ToLower(message.Action) {
	case "fetch":
		applications, err := ctx.manager.GetList()
		if err != nil {
			logger.Error("failed to fetch application list", zap.Error(err))
			return
		}

		logger.Info("events successfully processed", zap.Any("applications", applications))
	case "create":
		application, err := ctx.manager.Create(
			message.Service,
			message.Version,
			message.Namespace,
			message.Values,
		)
		if application != nil && err == nil {
			logger.Info("application created", zap.Any("application", application))
		}
	default:
		logger.Warn("unknown action", zap.String("action", message.Action))
	}
}

func (ctx *EventProcessor) processSystemMessage(message *models.SystemMessage) {
	for idx := 0; idx < 3; idx++ {
		if err := ctx.producer.Produce(message.Key, message.Value); err != nil {
			logger.Error(
				fmt.Sprintf("failed to produce message (try: %d)", idx),
				zap.Any("message", message),
				zap.Error(err),
			)
		}

		break
	}
}
