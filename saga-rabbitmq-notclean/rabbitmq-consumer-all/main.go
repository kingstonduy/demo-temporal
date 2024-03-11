package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saga-rabbitmq-notclean/config"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"

	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ConsumeAndPublish(topic string, conn *amqp.Connection) {
	concurrency := 1000
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
		concurrency, // prefetch count
		0,           // prefetch size
		false,       // global
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

	var wg sync.WaitGroup // used to coordinate when they are done, ie: if rabbit conn was closed
	wg.Add(concurrency)

	for x := 0; x < concurrency; x++ {
		go func() {
			defer wg.Done()
			for d := range msgs {
				// convert json string to struct
				var saferRequest model.SaferRequest
				fmt.Printf("💡Consume request from topic %s, message: %+v\n", topic, string(d.Body))
				err := json.Unmarshal(d.Body, &saferRequest)
				failOnError(err, "Failed to convert body to struct")

				var response = model.SaferResponse{
					WorkflowID: saferRequest.WorkflowID,
					RunID:      saferRequest.RunID,
					Code:       http.StatusOK,
					Status:     "OK",
					Message:    "Success",
				}

				if topic == config.GetConfig().NapasAccount.Queue {
					napasResponse := model.NapasAccountResponse{
						AccountID:   saferRequest.ToAccountID,
						AccountName: "Nguyen Van A",
					}
					jsonBytes, _ := json.Marshal(napasResponse)
					response.Message = string(jsonBytes)
				}

				fmt.Printf("💡Consume request from topic %s, message: %+v\n", topic, saferRequest)

				// convert struct to json string
				responseStr, _ := json.Marshal(response)
				fmt.Printf("💡Reply response from topic %s, message: %+v\n", topic, response)

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

	conn, err := amqp.Dial(RabbitMQ_URL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		ConsumeAndPublish(config.GetConfig().NapasMoney.Queue, conn)
	}()

	go func() {
		defer wg.Done()
		ConsumeAndPublish(config.GetConfig().NapasAccount.Queue, conn)
	}()

	go func() {
		defer wg.Done()
		ConsumeAndPublish(config.GetConfig().Limit.Queue, conn)
	}()

	go func() {
		defer wg.Done()
		ConsumeAndPublish(config.GetConfig().T24.Queue, conn)
	}()

	wg.Wait()
}
