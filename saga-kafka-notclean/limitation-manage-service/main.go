package main

import (
	"encoding/json"
	"fmt"
	"log"
	shared "saga-kafka-notclean/config"

	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func update(req shared.SaferRequest) (model shared.SaferResponse, err error) {
	log.Printf("💡Request %+v\n", req)

	db, err := shared.GetDB()
	if err != nil {
		log.Println(err)
		return
	}

	var accountLimitEntity shared.AccountLimitEntity
	err = db.Where("account_id = ?", req.FromAccountID).First(&accountLimitEntity).Error
	if err != nil {
		log.Println(err)
		return
	}

	accountLimitEntity.Amount += req.Amount
	if accountLimitEntity.Amount < 0 {
		log.Println(err)
		return
	}

	err = db.Save(&accountLimitEntity).Error
	if err != nil {
		log.Println(err)
		return
	}

	return shared.SaferResponse{
		WorkflowID: req.WorkflowID,
		RunID:      req.RunID,
		SignalName: "limit",
		Code:       http.StatusOK,
		Message:    "update record napas success",
	}, nil
}

func Produce[T any](topic string, message T) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			shared.GetConfig().Kafka.BootstrapServer.Host,
			shared.GetConfig().Kafka.BootstrapServer.Port,
		),
	})
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	// Produce messages to topic (asynchronously)
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(jsonBytes),
	}, nil)

	// Wait for message deliveries
	p.Flush(15 * 1000)
	p.Close()
}

// consume from kafka then signal the workflow to continue
func ConsumeAndProduce[T any, K any](topicIn string, topicOut string, req T, res K) {

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

	c.SubscribeTopics([]string{topicIn}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			err := json.Unmarshal([]byte(string(msg.Value)), &req)
			if err != nil {
				log.Println(err)
				return
			}

			res, err = update(req)

			Produce(topicOut, res)
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
}

func main() {
	var req model.SaferRequest
	var res model.SaferResponse
	ConsumeAndProduce(shared.GetConfig().Limit.Kafka.Topic.In, req, res)
}
