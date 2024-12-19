package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"strings"
	"tera/deployment/internal/ports"
	"tera/deployment/pkg/config"
	"tera/deployment/pkg/logger"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(conf *config.Config) ports.KafkaProducer {
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

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		logger.Error("NewProducer error", zap.Error(err))

		panic(err)
	}

	return &Producer{
		producer: producer,
		topic:    conf.Kafka.Topic,
	}
}

func (ctx *Producer) Produce(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("failed to marshal message value to JSON", zap.Error(err))

		return errors.New(fmt.Sprintf("failed to marshal value: %v", err))
	}

	events := make(chan kafka.Event)
	defer close(events)

	err = ctx.producer.Produce(&kafka.Message{
		Key:   []byte(key),
		Value: data,
		TopicPartition: kafka.TopicPartition{
			Topic:     lo.ToPtr(ctx.topic),
			Partition: kafka.PartitionAny,
		},
		Headers: []kafka.Header{applicationHeader},
	}, events)
	if err != nil {
		logger.Error("failed to produce Kafka message", zap.Error(err))
		return errors.New(fmt.Sprintf("failed to produce Kafka message: %v", err))
	}

	event := <-events
	switch e := event.(type) {
	case *kafka.Message:
		if e.TopicPartition.Error != nil {
			logger.Error("failed to produce message to Kafka topic",
				zap.String("topic", ctx.topic),
				zap.Error(e.TopicPartition.Error),
			)
			return errors.New(fmt.Sprintf("failed to produce message: %v", e.TopicPartition.Error))
		}
		logger.Info("successfully produced message",
			zap.String("topic", ctx.topic),
			zap.String("key", key),
			zap.Any("value", value),
		)
	case *kafka.Error:
		logger.Error("Kafka error during message production",
			zap.String("topic", ctx.topic),
			zap.Error(e),
		)
		return errors.New(fmt.Sprintf("kafka error: %v", e))
	default:
		logger.Error("unknown event received during Kafka production",
			zap.Any("event", e),
		)
		return errors.New(fmt.Sprintf("unknown event: %v", e))
	}

	return nil
}
