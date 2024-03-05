package bootstrap

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (c *Config) GetKafkaProducer() *kafka.Producer {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			c.Kafka.BootstrapServer.Host,
			c.Kafka.BootstrapServer.Port,
		),
	})
	if err != nil {
		log.Fatalf("Error creating Kafka producer, %s", err)
	}

	return producer
}
