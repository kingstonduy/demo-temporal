package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"saga-rabbitmq-notclean/money-transfer-service/config"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"

	"net/http"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func update(s string) model.SaferResponse {
	var req model.SaferRequest

	err := json.Unmarshal([]byte(s), &req)
	if err != nil {
		return model.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Code:       http.StatusBadRequest,
			Message:    err.Error(),
		}
	}
	log.Printf("💡Request %+v\n", req)

	db, err := config.GetDB()
	if err != nil {
		return model.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	var accountLimitEntity model.AccountLimitEntity
	err = db.Where("account_id = ?", req.FromAccountID).First(&accountLimitEntity).Error
	if err != nil {
		return model.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	accountLimitEntity.Amount += req.Amount
	if accountLimitEntity.Amount < 0 {
		return model.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Code:       http.StatusBadRequest,
			Message:    "Not enough money",
		}
	}

	err = db.Save(&accountLimitEntity).Error
	if err != nil {
		return model.SaferResponse{
			WorkflowID: req.WorkflowID,
			RunID:      req.RunID,
			Code:       http.StatusInternalServerError,
			Message:    err.Error(),
		}
	}

	return model.SaferResponse{
		WorkflowID: req.WorkflowID,
		RunID:      req.RunID,
		Code:       http.StatusOK,
		Message:    "update record napas success",
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ConsumeAndPublish(topic string, url string) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	concurrency := 10
	var wg sync.WaitGroup // used to coordinate when they are done, ie: if rabbit conn was closed
	wg.Add(concurrency)

	for x := 0; x < concurrency; x++ {
		go func() {
			defer wg.Done()
			for d := range msgs {
				log.Printf(" [*] Awaiting RPC requests")
				n := string(d.Body)
				failOnError(err, "Failed to convert body to integer")

				var response model.SaferResponse = update(n)

				// convert struct to json string
				responseStr, _ := json.Marshal(response)

				err = ch.PublishWithContext(ctx,
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(responseStr),
					})
				failOnError(err, "Failed to publish a message")

				d.Ack(false)
			}
		}()
	}
	wg.Wait() // when all goroutine's exit, the app exits
}

func main() {
	var RabbitMQ_URL = fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.GetConfig().RabbitMQ.User,
		config.GetConfig().RabbitMQ.Password,
		config.GetConfig().RabbitMQ.Host,
		config.GetConfig().RabbitMQ.Port,
	)
	ConsumeAndPublish(config.GetConfig().Limit.Queue, RabbitMQ_URL)
}
