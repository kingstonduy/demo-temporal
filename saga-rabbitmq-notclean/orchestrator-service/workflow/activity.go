package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"saga-rabbitmq-notclean/money-transfer-service/config"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"
	"time"

	"github.com/pborman/uuid"
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

func RequestAndReply[T any, K any](req T, res *K, topic string, url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Panicf("%s: Failed to connect to RabbitMQ", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: Failed to open a channel", err)
		return err
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
		return err
	}

	corrId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inputString, err := json.Marshal(req)
	if err != nil {
		log.Panicf("Failed to convert object to JSON: %s", err)
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
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
			Body:          []byte(inputString),
		})
	if err != nil {
		log.Panicf("%s: Failed to publish a message", err)
		return err
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			err = json.Unmarshal(d.Body, res)
			if err != nil {
				log.Panicf("%s: Failed to convert json to  object", err)
				return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
			}
			break
		}
	}

	return nil
}

func ValidateAccount(ctx context.Context, input model.SaferRequest) (output model.NapasAccountResponse, err error) {
	log := activity.GetLogger(ctx)
	log.Info("💡Validate Account activity starts")

	var response model.SaferResponse
	err = RequestAndReply(input, &response, config.GetConfig().NapasAccount.Queue, RabbitMQ_URL)
	if err != nil {
		return output, err
	}

	err = json.Unmarshal([]byte(response.Message), &output)
	if err != nil {
		return output, temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	if response.Code != 200 {
		return output, temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Validate Account activity successfully")
	return
}

func LimitCut(ctx context.Context, input model.SaferRequest) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut Account activity starts")

	// fix this

	input.Amount = -int64(math.Abs(float64(input.Amount)))
	var response model.SaferResponse
	err := RequestAndReply(input, &response, config.GetConfig().Limit.Queue, RabbitMQ_URL)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	log.Info("💡Limit cut Account activity successfully")
	return nil
}

func LimitCutCompensate(ctx context.Context, input model.SaferRequest) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut compensate activity starts")

	input.Amount = int64(math.Abs(float64(input.Amount)))
	var response model.SaferResponse
	err := RequestAndReply(input, &response, config.GetConfig().Limit.Queue, RabbitMQ_URL)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func MoneyCut(ctx context.Context, input model.SaferRequest) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut Account activity starts")

	input.Amount = -int64(math.Abs(float64(input.Amount)))
	var response model.SaferResponse
	err := RequestAndReply(input, &response, config.GetConfig().T24.Queue, RabbitMQ_URL)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	log.Info("💡Money cut Account activity successfully")
	return nil
}

func MoneyCutCompensate(ctx context.Context, input model.SaferRequest) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut compensate activity starts")

	input.Amount = int64(math.Abs(float64(input.Amount)))
	var response model.SaferResponse
	err := RequestAndReply(input, &response, config.GetConfig().T24.Queue, RabbitMQ_URL)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	log.Info("💡Money cut compensate activity successfully")
	return nil
}

func UpdateMoney(ctx context.Context, input model.SaferRequest) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Add money to receiver activity starts")

	input.Amount = int64(math.Abs(float64(input.Amount)))
	var response model.SaferResponse
	err := RequestAndReply(input, &response, config.GetConfig().NapasMoney.Queue, RabbitMQ_URL)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	log.Info("💡Add money to receiver activity successfully")
	return nil
}

func UpdateMoneyCompensate(ctx context.Context, input model.SaferRequest) error {
	log := activity.GetLogger(ctx)
	log.Info("💡update money napas cut compensate activity starts")

	input.Amount = -int64(math.Abs(float64(input.Amount)))
	var response model.SaferResponse
	err := RequestAndReply(input, &response, config.GetConfig().NapasMoney.Queue, RabbitMQ_URL)
	if err != nil {
		return err
	}
	if response.Code != 200 {
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
