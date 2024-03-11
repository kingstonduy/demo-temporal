package bootstrap

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func GetAMQPConnection(c *Config) *amqp.Connection {
	var conn *amqp.Connection = nil

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)

	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ, %s", err)
	}

	return conn
}
