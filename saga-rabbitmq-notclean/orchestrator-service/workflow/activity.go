package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"saga-rabbitmq-notclean/money-transfer-service/config"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

var RabbitMQ_URL = fmt.Sprintf("amqp://%s:%s@%s:%s/",
	config.GetConfig().RabbitMQ.User,
	config.GetConfig().RabbitMQ.Password,
	config.GetConfig().RabbitMQ.Host,
	config.GetConfig().RabbitMQ.Port,
)

func IsRetryableError(err error) bool {
	return (string)(err.Error()[0]) == "1"
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func RequestAndReply(topic string, url string, message string) (res string, err error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Panicf("%s: Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: Failed to open a channel", err)
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
	}

	corrId := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(message),
		})
	if err != nil {
		log.Panicf("%s: Failed to publish a message", err)
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			break
		}
	}

	return
}

func ValidateAccount(ctx context.Context, input model.WorkflowInput) (output model.NapasAccountResponse, err error) {
	log := activity.GetLogger(ctx)
	log.Info("💡Validate Account activity starts")

	inputJson, err := json.Marshal(input)

	resJson, err := RequestAndReply(config.GetConfig().NapasAccount.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &output)

	if err != nil {
		return output, temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Validate Account activity successfully")
	return
}

func LimitCut(ctx context.Context, input model.WorkflowInput) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut Account activity starts")

	inputJson, err := json.Marshal(input)
	var responseType model.SaferResponse

	resJson, err := RequestAndReply(config.GetConfig().Limit.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &responseType)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Limit cut Account activity successfully")
	return nil
}

func LimitCutCompensate(ctx context.Context, input model.WorkflowInput) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut compensate activity starts")

	inputJson, err := json.Marshal(input)
	var responseType model.SaferResponse

	resJson, err := RequestAndReply(config.GetConfig().Limit.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &responseType)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func MoneyCut(ctx context.Context, input model.WorkflowInput) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut Account activity starts")

	inputJson, err := json.Marshal(input)
	var responseType model.SaferResponse

	resJson, err := RequestAndReply(config.GetConfig().T24.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &responseType)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Money cut Account activity successfully")
	return nil
}

func MoneyCutCompensate(ctx context.Context, input model.WorkflowInput) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut compensate activity starts")

	inputJson, err := json.Marshal(input)
	var responseType model.SaferResponse

	resJson, err := RequestAndReply(config.GetConfig().T24.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &responseType)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Money cut compensate activity successfully")
	return nil
}

func UpdateMoney(ctx context.Context, input model.WorkflowInput) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Add money to receiver activity starts")

	inputJson, err := json.Marshal(input)
	var responseType model.SaferResponse

	resJson, err := RequestAndReply(config.GetConfig().NapasMoney.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &responseType)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Add money to receiver activity successfully")
	return nil
}

func UpdateMoneyCompensate(ctx context.Context, input model.WorkflowInput) error {
	log := activity.GetLogger(ctx)
	log.Info("💡update money napas cut compensate activity starts")

	inputJson, err := json.Marshal(input)
	var responseType model.SaferResponse

	resJson, err := RequestAndReply(config.GetConfig().NapasMoney.Queue, RabbitMQ_URL, string(inputJson))
	err = json.Unmarshal([]byte(resJson), &responseType)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Update money compensate activity successfully")
	return nil
}

func UpdateStateCreated(ctx context.Context, input model.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	input.State = "CREATED"

	db, _ := config.GetDB()
	err := db.Create(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Persist transaction to database successfully")
	return nil
}

func UpdateStateLimitCut(ctx context.Context, input model.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	input.State = "LIMIT_CUT"

	db, _ := config.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Update state successfully")
	return nil
}

func UpdateStateMoneyCut(ctx context.Context, input model.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	input.State = "MONEY_CUT"

	db, _ := config.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Update state successfully")
	return nil
}

func UpdateStateTransactionCompleted(ctx context.Context, input model.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	input.State = "COMPLETED"

	db, _ := config.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Update state  successfully")
	return nil
}

func CompensateTransaction(ctx context.Context, input model.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Compensate Transaction starts")

	input.State = "CANCEL"

	db, _ := config.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Compensate Transaction successfully")
	return nil
}
