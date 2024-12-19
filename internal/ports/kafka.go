package ports

import "tera/deployment/internal/domain/models"

type KafkaConsumer interface {
	Start(events chan<- *models.EventMessage) error
	Close() error
}

type KafkaProducer interface {
	Produce(key string, value any) error
}
