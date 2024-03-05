package infra

import (
	"consumer-service/bootstrap"
	"consumer-service/pkg/logger"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	cfg      bootstrap.Config
}

func NewKafkaConsumer(log logger.Logger, cfg bootstrap.Config) *KafkaConsumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			cfg.Kafka.BootstrapServer.Host,
			cfg.Kafka.BootstrapServer.Port,
		),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	return &KafkaConsumer{
		consumer: c,
		cfg:      cfg,
	}
}

func (k *KafkaConsumer) SubscribeTopics(topics []string) {
	k.consumer.SubscribeTopics(topics, nil)
}

func (k *KafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	return k.consumer.ReadMessage(-1)
}
