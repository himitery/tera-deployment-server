package kafka

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

var applicationHeader = kafka.Header{
	Key:   "source",
	Value: []byte("tera-deployment-server"),
}
