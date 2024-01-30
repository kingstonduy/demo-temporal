package usecase

import (
	"context"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"
	"time"
)

type moneyTransferUsecase struct {
	transactionRepository domain.TransactionRepository
	ContextTimeout        time.Duration
}

func NewMoneyTransferUsecase(transactionRepository domain.TransactionRepository, timeout time.Duration) domain.MoneyTransferUsecase {
	return &moneyTransferUsecase{
		transactionRepository: transactionRepository,
		ContextTimeout:        timeout,
	}
}

func (m *moneyTransferUsecase) ValidateAccount(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) UpdateStateCreated(ctx context.Context, input domain.TransactionEntity) error {
	return nil
}

func (m *moneyTransferUsecase) CompensateTransaction(ctx context.Context, input domain.TransactionEntity) error {
	return nil
}

func (m *moneyTransferUsecase) LimitCut(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) LimitCutCompensate(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) UpdateStateLimitCut(ctx context.Context, input domain.TransactionEntity) error {
	return nil
}

func (m *moneyTransferUsecase) MoneyCut(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) MoneyCutCompensate(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) UpdateStateMoneyCut(ctx context.Context, input domain.TransactionEntity) error {
	return nil
}

func (m *moneyTransferUsecase) UpdateMoney(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) UpdateMoneyCompensate(ctx context.Context, input domain.TransactionInfo) error {
	return nil
}

func (m *moneyTransferUsecase) UpdateStateTransactionCompleted(ctx context.Context, input domain.TransactionEntity) error {
	return nil
}
