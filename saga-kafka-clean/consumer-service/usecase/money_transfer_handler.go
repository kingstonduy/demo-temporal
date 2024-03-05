package handlers

import (
	"consumer-service/bootstrap"
	"consumer-service/domain"
	"consumer-service/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type moneyTransferHandler struct {
	log logger.Logger
	cfg *bootstrap.Config
}

func NewMoneyTransferHandler(log logger.Logger, cfg *bootstrap.Config) domain.MoneyTransferHandler {
	return &moneyTransferHandler{
		log: log,
		cfg: cfg,
	}
}

// Handle implements domain.MoneyTransferHandler.
func (h *moneyTransferHandler) Handle(message domain.SaferResponse) error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			h.cfg.Kafka.BootstrapServer.Host,
			h.cfg.Kafka.BootstrapServer.Port,
		),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	topics := []string{
		h.cfg.T24.Kafka.Topic.Out,
		h.cfg.Limit.Kafka.Topic.Out,
		h.cfg.Napas.Kafka.Topic.Out,
	}

	fmt.Println("topics", topics)

	c.SubscribeTopics(topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				panic(err)
			}
			log.Printf("💡Consume message = %+v", message)
			log.Printf("💡Signal to temporal workflowID = %s, workflowRunID = %s, signal = %s, message = %+v",
				message.WorkflowID, message.RunID, message.Action, message)

			(*h.cfg.GetTemporalClient()).SignalWorkflow(context.Background(), message.WorkflowID, message.RunID, message.Action, message)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
	return nil
}
