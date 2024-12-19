package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strings"
	"tera/deployment/internal/domain/models"
	"tera/deployment/internal/ports"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
)

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewKafkaConsumer(conf *config.Config) ports.KafkaConsumer {
	bootstrapServers := lo.Map(conf.Kafka.BootstrapServers, func(server config.KafkaBootstrapServerConfig, _ int) string {
		return fmt.Sprintf("%s:%d", server.Host, server.Port)
	})

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(bootstrapServers, ","),
		"group.id":          "tera-deployment",
	}
	if conf.Kafka.Protocol == "SASL" {
		_ = kafkaConfig.SetKey("security.protocol", "SASL")
		_ = kafkaConfig.SetKey("sasl.mechanisms", conf.Kafka.Protocol)
		_ = kafkaConfig.SetKey("sasl.username", conf.Kafka.Sasl.Username)
		_ = kafkaConfig.SetKey("sasl.password", conf.Kafka.Sasl.Password)
	}

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		logger.Error("failed to create consumer", zap.Error(err))

		panic(err)
	}

	return &Consumer{
		consumer: consumer,
		topic:    conf.Kafka.Topic,
	}
}

func (ctx *Consumer) Start(events chan<- *models.EventMessage) error {
	if err := ctx.consumer.Subscribe(ctx.topic, nil); err != nil {
		return err
	}

	go func() {
		logger.Info("consumer started", zap.String("topic", ctx.topic))

		for {
			switch event := ctx.consumer.Poll(100).(type) {
			case *kafka.Message:
				if len(lo.FilterMap(event.Headers, func(item kafka.Header, _ int) (kafka.Header, bool) {
					return item, item.Key == applicationHeader.Key && string(item.Value) == string(applicationHeader.Value)
				})) != 0 {
					continue
				}

				var message *models.EventMessage
				if err := json.Unmarshal(event.Value, &message); err != nil {
					logger.Warn("failed to unmarshal event", zap.Error(err))
					continue
				}

				events <- message
			case kafka.Error:
				logger.Error("kafka error", zap.Error(event))
			}
		}
	}()

	return nil
}

func (ctx *Consumer) Close() error {
	if !ctx.consumer.IsClosed() {
		if err := ctx.consumer.Close(); err != nil {
			logger.Error("failed to close consumer", zap.Error(err))

			return errors.Wrap(err, "failed to close consumer")
		}
	}

	return nil
}
