package ports

import "tera/deployment/internal/domain/models"

type KafkaConsumer interface {
	Start(events chan<- *models.EventMessage) error
	Close() error
}

type KafkaProducer interface {
	Produce(topic string, key, value []byte) error
}