package infra

import (
	"log"
	"orchestrator-service/bootstrap"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	cfg *bootstrap.Config
}

func NewKafkaClient(cfg *bootstrap.Config) *KafkaClient {
	return &KafkaClient{
		cfg: cfg,
	}
}

func (c *KafkaClient) Produce(topic string, message string) error {
	log.Printf("💡Send to topic: %s, message: %s", topic, message)
	p := c.cfg.GetKafkaProducer()
	return p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
}
