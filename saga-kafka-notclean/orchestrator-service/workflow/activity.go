package workflow

import (
	"context"
	"encoding/json"
	"log"
	"math"
	shared "saga-kafka-notclean/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.temporal.io/sdk/temporal"
)

func Produce(topic string, message string) error {
	log.Printf("💡Send to topic: %s, message: %s", topic, message)
	p := shared.GetKafkaProducer()
	return p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
}

func ValidateAccount(ctx context.Context, input shared.SaferRequest) error {
	// convert struct to json string
	input.Action = "napas-account"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().Napas.Kafka.Topic.In, string(jsonInput))
}

func LimitCut(ctx context.Context, input shared.SaferRequest) error {
	input.Amount = -int64(math.Abs(float64(input.Amount)))
	input.Action = "limit"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().Limit.Kafka.Topic.In, string(jsonInput))
}

func LimitCutCompensate(ctx context.Context, input shared.SaferRequest) error {
	input.Amount = int64(math.Abs(float64(input.Amount)))
	input.Action = "limit"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().Limit.Kafka.Topic.In, string(jsonInput))
}

func MoneyCut(ctx context.Context, input shared.SaferRequest) error {
	input.Amount = -int64(math.Abs(float64(input.Amount)))
	input.Action = "t24"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().T24.Kafka.Topic.In, string(jsonInput))
}

func MoneyCutCompensate(ctx context.Context, input shared.SaferRequest) error {
	input.Amount = int64(math.Abs(float64(input.Amount)))
	input.Action = "t24"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().T24.Kafka.Topic.In, string(jsonInput))
}

func UpdateMoney(ctx context.Context, input shared.SaferRequest) error {
	input.Amount = int64(math.Abs(float64(input.Amount)))
	input.Action = "napas-money"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().Napas.Kafka.Topic.In, string(jsonInput))
}

func UpdateMoneyCompensate(ctx context.Context, input shared.SaferRequest) error {
	input.Amount = -int64(math.Abs(float64(input.Amount)))
	input.Action = "napas-money"

	jsonInput, _ := json.Marshal(input)
	return Produce(shared.GetConfig().Napas.Kafka.Topic.In, string(jsonInput))
}

func UpdateStateCreated(ctx context.Context, input shared.TransactionEntity) error {
	input.State = "CREATED"

	db, _ := shared.GetDB()
	err := db.Create(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}

func UpdateStateLimitCut(ctx context.Context, input shared.TransactionEntity) error {

	input.State = "LIMIT_CUT"

	db, _ := shared.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}

func UpdateStateMoneyCut(ctx context.Context, input shared.TransactionEntity) error {

	input.State = "MONEY_CUT"

	db, _ := shared.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}

func UpdateStateTransactionCompleted(ctx context.Context, input shared.TransactionEntity) error {

	input.State = "COMPLETED"

	db, _ := shared.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}

func CompensateTransaction(ctx context.Context, input shared.TransactionEntity) error {

	input.State = "CANCEL"

	db, _ := shared.GetDB()
	err := db.Save(input).Error
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	return nil
}
