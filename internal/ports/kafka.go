package ports

import (
	"tera/deployment/internal/domain/models"
)

type KafkaConsumer interface {
	Start(events chan<- any) error
	Close() error
}

type KafkaProducer interface {
	Produce(key models.Key, value any) error
}
