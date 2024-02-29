package usecase

import (
	"context"
	"fmt"
	shared "saga"
	"saga/money-transfer-service/domain"
	"saga/money-transfer-service/internal/httputil"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
)

type moneyTransferUsecase struct {
	transactionRepository domain.TransactionRepository
	napasUrl              string
	t24Url                string
	limitServiceUrl       string
	ContextTimeout        time.Duration
}

func NewMoneyTransferUsecase(transactionRepository domain.TransactionRepository, napasUrl string,
	t24Url string, limitServiceUrl string, timeout time.Duration) domain.MoneyTransferUsecase {
	return &moneyTransferUsecase{
		transactionRepository: transactionRepository,
		napasUrl:              napasUrl,
		t24Url:                t24Url,
		limitServiceUrl:       limitServiceUrl,
		ContextTimeout:        timeout,
	}
}

func (m *moneyTransferUsecase) ValidateAccount(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Validate Account activity starts")
	url := fmt.Sprintf("http://%s/api/v1/account/verify", m.napasUrl)
	var responseType domain.NapasEntity

	err := httputil.PostApi(url, &domain.ValidateAccountInput{input.ToAccountId}, &responseType)
	if err != nil {
		log.Error("🔥Validate Account activity failed")
		if shared.IsRetryableError(err) {
			return err
		} else {
			return temporal.NewNonRetryableApplicationError("non retry", shared.NONRETRYABLE_ERROR, err, nil)
		}
	}

	log.Info("💡Validate Account activity successfully")
	return nil
}

func (m *moneyTransferUsecase) UpdateStateCreated(ctx context.Context, input domain.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	err := m.transactionRepository.Create(ctx, input)
	if err != nil {
		log.Error("🔥Cannot connect to database")
		return err
	}
	return nil
}

func (m *moneyTransferUsecase) CompensateTransaction(ctx context.Context, input domain.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Compensate Transaction starts")

	err := m.transactionRepository.Update(ctx, input)
	if err != nil {
		log.Error("🔥Cannot create compensate transaction")
		return err
	}

	log.Info("💡Compensate Transaction successfully")
	return nil
}

func (m *moneyTransferUsecase) LimitCut(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut Account activity starts")
	url := fmt.Sprintf("http://%s/api/v1/account/limit", m.limitServiceUrl)
	var responseType domain.NapasEntity

	err := httputil.PostApi(url,
		&domain.SaferRequest{
			TransactionId: input.TransactionId,
			AccountId:     input.FromAccountId,
			Amount:        input.Amount,
		}, &responseType)
	if err != nil {
		log.Error("🔥Limit cut Account activity failed")
		if shared.IsRetryableError(err) {
			return err
		} else {
			return temporal.NewNonRetryableApplicationError("non retry", shared.NONRETRYABLE_ERROR, err, nil)
		}
	}

	log.Info("💡Limit cut Account activity successfully")
	return nil
}

func (m *moneyTransferUsecase) LimitCutCompensate(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Limit cut compensate activity starts")
	url := fmt.Sprintf("http://%s/api/v1/account/limit", m.limitServiceUrl)

	var responseType domain.NapasEntity
	err := httputil.PostApi(url,
		&domain.SaferRequest{
			TransactionId: input.TransactionId,
			AccountId:     input.FromAccountId,
			Amount:        -input.Amount, // compensate
		}, &responseType)
	if err != nil {
		log.Error("🔥Limit cut compensate activity failed")
		return err
	}

	log.Info("💡Limit cut compensate activity successfully")
	return nil
}

func (m *moneyTransferUsecase) UpdateStateLimitCut(ctx context.Context, input domain.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	err := m.transactionRepository.Update(ctx, input)
	if err != nil {
		log.Error("🔥Cannot update state")
		return err
	}

	log.Info("💡Update state successfully")
	return nil
}

func (m *moneyTransferUsecase) MoneyCut(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut Account activity starts")
	url := fmt.Sprintf("http://%s/api/v1/amount/cut", m.t24Url)
	var responseType domain.SaferResponse

	err := httputil.PostApi(url, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)
	if err != nil {
		log.Error("🔥Money cut Account activity failed")
		if shared.IsRetryableError(err) {
			return err
		} else {
			return temporal.NewNonRetryableApplicationError("non retry", shared.NONRETRYABLE_ERROR, err, nil)
		}
	}

	log.Info("💡Money cut Account activity successfully")
	return nil
}

func (m *moneyTransferUsecase) MoneyCutCompensate(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Money cut compensate activity starts")
	url := fmt.Sprintf("http://%s/api/v1/amount/add", m.t24Url)
	var responseType domain.SaferResponse

	err := httputil.PostApi(url, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)

	if err != nil {
		log.Error("🔥Money cut compensate activity failed")
		return err
	}

	log.Info("💡Money cut compensate activity successfully")
	return nil
}

func (m *moneyTransferUsecase) UpdateStateMoneyCut(ctx context.Context, input domain.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")

	err := m.transactionRepository.Update(ctx, input)
	if err != nil {
		log.Error("🔥Cannot update state")
		return err
	}

	log.Info("💡Update state successfully")
	return nil
}

func (m *moneyTransferUsecase) UpdateMoney(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Add money to receiver activity starts")

	url := fmt.Sprintf("http://%s/api/v1/account/update", m.napasUrl)
	var responseType domain.NapasEntity
	err := httputil.PostApi(url, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.ToAccountId,
		Amount:        input.Amount,
	}, &responseType)
	if err != nil {
		log.Error("🔥Add money to receiver activity failed")
		return err
	}

	log.Info("💡Add money to receiver activity successfully")
	return nil
}

func (m *moneyTransferUsecase) UpdateMoneyCompensate(ctx context.Context, input domain.TransactionInfo) error {
	log := activity.GetLogger(ctx)
	log.Info("💡update money napas cut compensate activity starts")
	url := fmt.Sprintf("http://%s/api/v1/account/update", m.napasUrl)
	var responseType domain.NapasEntity

	err := httputil.PostApi(url, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.ToAccountId,
		Amount:        input.Amount,
	}, &responseType)

	if err != nil {
		log.Error("🔥Update money  compensate activity failed")
		return err
	}
	log.Info("💡Update money compensate activity successfully")
	return nil
}

func (m *moneyTransferUsecase) UpdateStateTransactionCompleted(ctx context.Context, input domain.TransactionEntity) error {
	log := activity.GetLogger(ctx)
	log.Info("💡Persist transaction to database starts")
	err := m.transactionRepository.Update(ctx, input)
	if err != nil {
		log.Error("🔥Cannot update state")
		return err
	}
	log.Info("💡Update state  successfully")
	return nil
}
