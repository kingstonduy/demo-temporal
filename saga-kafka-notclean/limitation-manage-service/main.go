package main

import (
	"encoding/json"
	"fmt"
	"log"
	shared "saga-kafka-notclean/config"
	"sync"

	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ch = make(chan shared.SaferResponse)

func update(s string) (model shared.SaferResponse, err error) {
	var req shared.SaferRequest
	err = json.Unmarshal([]byte(s), &req)

	return shared.SaferResponse{
		WorkflowID: req.WorkflowID,
		RunID:      req.RunID,
		Action:     "limit",
		Code:       http.StatusOK,
		Message:    "update record napas success",
	}, nil
}

func Produce(topic string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			shared.GetConfig().Kafka.BootstrapServer.Host,
			shared.GetConfig().Kafka.BootstrapServer.Port,
		),
	})
	if err != nil {
		panic(err)
	}

	for message := range ch {
		// convert message into json string
		messageString, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}

		log.Printf("💡Response to topic %s, message = %s\n", topic, messageString)
		// Produce messages to topic (asynchronously)
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(messageString),
		}, nil)

		// Wait for message deliveries
	}

	defer p.Close()
	defer close(ch)
}

// consume from kafka then signal the workflow to continue
func Consume(topic string, req shared.SaferRequest) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			shared.GetConfig().Kafka.BootstrapServer.Host,
			shared.GetConfig().Kafka.BootstrapServer.Port,
		),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	log.Printf("💡Consume from topic: %s\n", topic)
	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)

		if err == nil {
			log.Println("📩 Received from kafka", string(msg.Value))
			res, _ := update(string(msg.Value))

			ch <- res
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
}

func Handler() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		Consume(shared.GetConfig().Limit.Kafka.Topic.In, shared.SaferRequest{})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Produce(shared.GetConfig().Limit.Kafka.Topic.Out)
	}()
	wg.Wait()
}

func main() {
	Handler()
}
