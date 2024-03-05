package bootstrap

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Config) GetAMQPConnection() *amqp.Connection {
	var conn *amqp.Connection = nil
	if conn != nil {
		return conn
	}

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
