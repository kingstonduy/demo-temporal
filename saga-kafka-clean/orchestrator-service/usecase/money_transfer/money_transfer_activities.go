package usecase

import (
	"context"
	"encoding/json"
	"math"
	"orchestrator-service/bootstrap"
	"orchestrator-service/domain"
	infra "orchestrator-service/infra/kafka"
	"orchestrator-service/infra/logger"

	"go.temporal.io/sdk/temporal"
)

type MoneyTransferActivities struct {
	log        logger.Logger
	repository domain.MoneyTransferRepository
	kafka      *infra.KafkaClient
	cfg        *bootstrap.Config
}

func NewMoneyTransferActivities(
	log logger.Logger,
	repository domain.MoneyTransferRepository,
	kafka *infra.KafkaClient,
	cfg *bootstrap.Config,
) domain.MoneyTransferActivities {
	return &MoneyTransferActivities{
		log:        log,
		repository: repository,
		kafka:      kafka,
		cfg:        cfg,
	}
}

func (m *MoneyTransferActivities) ValidateAccount(ctx context.Context, input domain.SaferRequest) error {
	// convert struct to json string
	input.Action = "napas-account"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.Napas.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) LimitCut(ctx context.Context, input domain.SaferRequest) error {
	input.Amount = -int64(math.Abs(float64(input.Amount)))
	input.Action = "limit"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.Limit.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) LimitCutCompensate(ctx context.Context, input domain.SaferRequest) error {
	input.Amount = int64(math.Abs(float64(input.Amount)))
	input.Action = "limit"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.Limit.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) MoneyCut(ctx context.Context, input domain.SaferRequest) error {
	input.Amount = -int64(math.Abs(float64(input.Amount)))
	input.Action = "t24"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.T24.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) MoneyCutCompensate(ctx context.Context, input domain.SaferRequest) error {
	input.Amount = int64(math.Abs(float64(input.Amount)))
	input.Action = "t24"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.T24.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) UpdateMoney(ctx context.Context, input domain.SaferRequest) error {
	input.Amount = int64(math.Abs(float64(input.Amount)))
	input.Action = "napas-money"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.Napas.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) UpdateMoneyCompensate(ctx context.Context, input domain.SaferRequest) error {
	input.Amount = -int64(math.Abs(float64(input.Amount)))
	input.Action = "napas-money"

	jsonInput, _ := json.Marshal(input)
	return m.kafka.Produce(m.cfg.Napas.Kafka.Topic.In, string(jsonInput))
}

func (m *MoneyTransferActivities) UpdateStateCreated(ctx context.Context, input domain.TransactionEntity) error {

	m.log.Info("💡Persist transaction to database starts")

	input.State = "CREATED"

	err := m.repository.Create(input)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	m.log.Info("💡Persist transaction to database successfully")
	return nil
}

func (m *MoneyTransferActivities) UpdateStateLimitCut(ctx context.Context, input domain.TransactionEntity) error {

	m.log.Info("💡Persist transaction to database starts")

	input.State = "LIMIT_CUT"

	err := m.repository.Save(input)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	m.log.Info("💡Update state successfully")
	return nil
}

func (m *MoneyTransferActivities) UpdateStateMoneyCut(ctx context.Context, input domain.TransactionEntity) error {

	m.log.Info("💡Persist transaction to database starts")

	input.State = "MONEY_CUT"

	err := m.repository.Save(input)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	m.log.Info("💡Update state successfully")
	return nil
}

func (m *MoneyTransferActivities) UpdateStateTransactionCompleted(ctx context.Context, input domain.TransactionEntity) error {

	m.log.Info("💡Persist transaction to database starts")

	input.State = "COMPLETED"

	err := m.repository.Save(input)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	m.log.Info("💡Update state  successfully")
	return nil
}

func (m *MoneyTransferActivities) CompensateTransaction(ctx context.Context, input domain.TransactionEntity) error {

	m.log.Info("💡Compensate Transaction starts")

	input.State = "CANCEL"

	err := m.repository.Save(input)
	if err != nil {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	m.log.Info("💡Compensate Transaction successfully")
	return nil
}
