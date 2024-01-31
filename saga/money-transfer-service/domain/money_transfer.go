package domain

import "context"

type MoneyTransferWorkflowInput struct {
}

type MoneyTransferUsecase interface {
	ValidateAccount(ctx context.Context, input TransactionInfo) error
	UpdateStateCreated(ctx context.Context, input TransactionEntity) error
	CompensateTransaction(ctx context.Context, input TransactionEntity) error
	LimitCut(ctx context.Context, input TransactionInfo) error
	LimitCutCompensate(ctx context.Context, input TransactionInfo) error
	UpdateStateLimitCut(ctx context.Context, input TransactionEntity) error
	MoneyCut(ctx context.Context, input TransactionInfo) error
	MoneyCutCompensate(ctx context.Context, input TransactionInfo) error
	UpdateStateMoneyCut(ctx context.Context, input TransactionEntity) error
	UpdateMoney(ctx context.Context, input TransactionInfo) error
	UpdateMoneyCompensate(ctx context.Context, input TransactionInfo) error
	UpdateStateTransactionCompleted(ctx context.Context, input TransactionEntity) error
}
