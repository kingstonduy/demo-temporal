package pkg

import (
	"context"
	"log"
	"orchestrator-service/bootstrap"
	"time"

	"github.com/pborman/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.temporal.io/sdk/temporal"
)

func RequestAndReply(req string, topic string, cfg *bootstrap.Config) (res string, err error) {
	conn := bootstrap.GetAMQPConnection(cfg)

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: Failed to open a channel", err)
		return "", err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("%s: Failed to declare a queue", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("%s: Failed to register a consumer", err)
		return "", err
	}

	corrId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err != nil {
		log.Panicf("Failed to convert object to JSON: %s", err)
		return "", temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	err = ch.PublishWithContext(ctx,
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(req),
		})
	if err != nil {
		log.Panicf("%s: Failed to publish a message", err)
		return "", err
	}
	defer ch.Close()

	for d := range msgs {
		if corrId == d.CorrelationId {
			if err != nil {
				log.Panicf("%s: Failed to convert json to  object", err)
				return "", temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
			}
			res = string(d.Body)
			break
		}
	}
	return res, nil
}
