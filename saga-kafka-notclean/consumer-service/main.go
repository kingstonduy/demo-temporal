package main

import (
	"context"
	"encoding/json"
	"fmt"
	shared "saga-kafka-notclean/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.temporal.io/sdk/client"
)

func main() {
	temporalClient, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", shared.GetConfig().Temporal.Host, shared.GetConfig().Temporal.Port),
	})
	if err != nil {
		panic(err)
	}
	defer temporalClient.Close()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			shared.GetConfig().Kafka.BootstrapServer.Host,
			shared.GetConfig().Kafka.BootstrapServer.Port,
		),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	topics := []string{
		shared.GetConfig().MoneyTransfer.Kafka.Topic.Out,
		shared.GetConfig().Limit.Kafka.Topic.Out,
		shared.GetConfig().Napas.Kafka.Topic.Out,
	}

	c.SubscribeTopics(topics, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var message shared.SaferResponse
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				panic(err)
			}

			temporalClient.SignalWorkflow(context.Background(), message.WorkflowID, message.RunID, message.SignalName, message)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
}
