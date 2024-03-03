package main

import (
	"encoding/json"
	"fmt"
	"log"
	shared "saga-kafka-notclean/config"

	"net/http"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func verify(s string) shared.SaferResponse {
	var req shared.SaferRequest

	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
		}
	}
	log.Printf("💡Request %+v\n", req)

	db, err := shared.GetDB()
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	var napasEntity shared.NapasEntity
	err = db.Where("account_id = ?", req.ToAccountID).First(&napasEntity).Error
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	jsonData, err := json.Marshal(napasEntity)
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return shared.SaferResponse{
		WorkflowID: req.WorkflowID,
		RunID:      req.RunID,
		Action:     req.Action,
		Code:       http.StatusOK,
		Message:    string(jsonData),
	}
}

func update(s string) shared.SaferResponse {
	var req shared.SaferRequest

	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
		}
	}
	log.Printf("💡Request %+v\n", req)

	db, err := shared.GetDB()
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	var napasEntity shared.NapasEntity
	err = db.Where("account_id = ?", req.ToAccountID).First(&napasEntity).Error
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	napasEntity.Amount += req.Amount
	if napasEntity.Amount < 0 {
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusBadRequest,
			Message:    "Not enough money",
		}
	}

	err = db.Save(&napasEntity).Error
	if err != nil {
		log.Println(err.Error())
		return shared.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Action:     req.Action,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return shared.SaferResponse{
		WorkflowID: req.WorkflowID,
		RunID:      req.RunID,
		Action:     req.Action,
		Code:       http.StatusOK,
		Message:    "update record napas success",
	}
}

var ch = make(chan shared.SaferResponse)

func Produce(topic string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s",
			shared.GetConfig().Kafka.BootstrapServer.Host,
			shared.GetConfig().Kafka.BootstrapServer.Port,
		),
	})
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	for message := range ch {
		// convert message into json string
		messageString, err := json.Marshal(message)
		if err != nil {
			log.Println(err.Error())
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
		log.Println(err.Error())
		panic(err)
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Println("📩 Received from kafka", string(msg.Value))
			err := json.Unmarshal([]byte(string(msg.Value)), &req)
			if err != nil {
				log.Println(err.Error())
				log.Println(err)
				return
			}

			if req.Action == "napas-account" {
				res := verify(string(msg.Value))
				ch <- res
			} else if req.Action == "napas-money" {
				res := update(string(msg.Value))
				ch <- res
			}
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
		Consume(shared.GetConfig().Napas.Kafka.Topic.In, shared.SaferRequest{})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Produce(shared.GetConfig().Napas.Kafka.Topic.Out)
	}()
	wg.Wait()
}

func main() {
	Handler()
}
