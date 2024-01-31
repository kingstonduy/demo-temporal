package usecase

import (
	"context"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"
	"kingstonduy/demo-temporal/saga/money-transfer-service/internal/httputil"
	"time"
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
	var responseType domain.NapasEntity
	return httputil.PostApi(m.napasUrl, &domain.ValidateAccountInput{input.ToAccountId}, &responseType)
}

func (m *moneyTransferUsecase) UpdateStateCreated(ctx context.Context, input domain.TransactionEntity) error {
	return m.transactionRepository.Create(ctx, input)
}

func (m *moneyTransferUsecase) CompensateTransaction(ctx context.Context, input domain.TransactionEntity) error {
	return m.transactionRepository.Update(ctx, input)
}

func (m *moneyTransferUsecase) LimitCut(ctx context.Context, input domain.TransactionInfo) error {
	var responseType domain.NapasEntity
	return httputil.PostApi(m.limitServiceUrl,
		&domain.SaferRequest{
			TransactionId: input.TransactionId,
			AccountId:     input.FromAccountId,
			Amount:        input.Amount,
		}, &responseType)
}

func (m *moneyTransferUsecase) LimitCutCompensate(ctx context.Context, input domain.TransactionInfo) error {
	var responseType domain.NapasEntity

	return httputil.PostApi(m.limitServiceUrl,
		&domain.SaferRequest{
			TransactionId: input.TransactionId,
			AccountId:     input.FromAccountId,
			Amount:        -input.Amount, // compensate
		}, &responseType)
}

func (m *moneyTransferUsecase) UpdateStateLimitCut(ctx context.Context, input domain.TransactionEntity) error {
	return m.transactionRepository.Update(ctx, input)
}

func (m *moneyTransferUsecase) MoneyCut(ctx context.Context, input domain.TransactionInfo) error {
	var responseType domain.SaferResponse
	return httputil.PostApi(m.t24Url, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)
}

func (m *moneyTransferUsecase) MoneyCutCompensate(ctx context.Context, input domain.TransactionInfo) error {
	var responseType domain.SaferResponse
	return httputil.PostApi(m.t24Url, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.FromAccountId,
		Amount:        input.Amount,
	}, &responseType)
}

func (m *moneyTransferUsecase) UpdateStateMoneyCut(ctx context.Context, input domain.TransactionEntity) error {
	return m.transactionRepository.Update(ctx, input)
}

func (m *moneyTransferUsecase) UpdateMoney(ctx context.Context, input domain.TransactionInfo) error {
	var responseType domain.NapasEntity
	return httputil.PostApi(m.napasUrl, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.ToAccountId,
		Amount:        input.Amount,
	}, &responseType)
}

func (m *moneyTransferUsecase) UpdateMoneyCompensate(ctx context.Context, input domain.TransactionInfo) error {
	var responseType domain.NapasEntity
	return httputil.PostApi(m.napasUrl, &domain.SaferRequest{
		TransactionId: input.TransactionId,
		AccountId:     input.ToAccountId,
		Amount:        input.Amount,
	}, &responseType)
}

func (m *moneyTransferUsecase) UpdateStateTransactionCompleted(ctx context.Context, input domain.TransactionEntity) error {
	return m.transactionRepository.Update(ctx, input)
}
