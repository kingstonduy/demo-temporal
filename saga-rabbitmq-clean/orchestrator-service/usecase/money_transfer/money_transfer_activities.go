package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"orchestrator-service/bootstrap"
	"orchestrator-service/domain"
	pkg "orchestrator-service/pkg/logger"

	"go.temporal.io/sdk/temporal"
)

type MoneyTransferActivities struct {
	cfg        *bootstrap.Config
	log        pkg.Logger
	mq         domain.MoneyTransferMessageQueue
	repository domain.MoneyTransferRepository
}

func NewMoneyTransferActivities(
	cfg *bootstrap.Config,
	log pkg.Logger,
	repository domain.MoneyTransferRepository,
	messagequeue domain.MoneyTransferMessageQueue,
) domain.MoneyTransferActivities {
	return &MoneyTransferActivities{
		cfg:        cfg,
		log:        log,
		repository: repository,
		mq:         messagequeue,
	}
}

func (m *MoneyTransferActivities) ValidateAccount(ctx context.Context, input domain.SaferRequest) (output domain.NapasAccountResponse, err error) {

	m.log.Info("💡Validate Account activity starts")

	var response domain.SaferResponse
	response, err = m.mq.SaferRequestResponse(input, m.cfg.NapasAccount.Queue)
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

	m.log.Info("💡Validate Account activity successfully")
	return
}

func (m *MoneyTransferActivities) LimitCut(ctx context.Context, input domain.SaferRequest) error {

	fmt.Println("💡HIHI")
	m.log.Info("💡Limit cut Account activity starts")

	// fix this

	input.Amount = -int64(math.Abs(float64(input.Amount)))
	var response domain.SaferResponse
	response, err := m.mq.SaferRequestResponse(input, m.cfg.Limit.Queue)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}

	m.log.Info("💡Limit cut Account activity successfully")
	return nil
}

func (m *MoneyTransferActivities) LimitCutCompensate(ctx context.Context, input domain.SaferRequest) error {

	m.log.Info("💡Limit cut compensate activity starts")

	input.Amount = int64(math.Abs(float64(input.Amount)))
	var response domain.SaferResponse
	response, err := m.mq.SaferRequestResponse(input, m.cfg.Limit.Queue)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	m.log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func (m *MoneyTransferActivities) MoneyCut(ctx context.Context, input domain.SaferRequest) error {

	m.log.Info("💡Money cut Account activity starts")

	input.Amount = -int64(math.Abs(float64(input.Amount)))
	var response domain.SaferResponse
	response, err := m.mq.SaferRequestResponse(input, m.cfg.T24.Queue)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	m.log.Info("💡Money cut Account activity successfully")
	return nil
}

func (m *MoneyTransferActivities) MoneyCutCompensate(ctx context.Context, input domain.SaferRequest) error {

	m.log.Info("💡Money cut compensate activity starts")

	input.Amount = int64(math.Abs(float64(input.Amount)))
	var response domain.SaferResponse
	response, err := m.mq.SaferRequestResponse(input, m.cfg.T24.Queue)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	m.log.Info("💡Money cut compensate activity successfully")
	return nil
}

func (m *MoneyTransferActivities) UpdateMoney(ctx context.Context, input domain.SaferRequest) error {

	m.log.Info("💡Add money to receiver activity starts")

	input.Amount = int64(math.Abs(float64(input.Amount)))
	var response domain.SaferResponse
	response, err := m.mq.SaferRequestResponse(input, m.cfg.NapasMoney.Queue)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	m.log.Info("💡Add money to receiver activity successfully")
	return nil
}

func (m *MoneyTransferActivities) UpdateMoneyCompensate(ctx context.Context, input domain.SaferRequest) error {

	m.log.Info("💡update money napas cut compensate activity starts")

	input.Amount = -int64(math.Abs(float64(input.Amount)))
	var response domain.SaferResponse
	response, err := m.mq.SaferRequestResponse(input, m.cfg.NapasMoney.Queue)
	if err != nil {
		return err
	}
	if response.Code != 200 {
		return temporal.NewNonRetryableApplicationError("non retry", "0", err, nil)
	}
	m.log.Info("💡Update money compensate activity successfully")
	return nil
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
