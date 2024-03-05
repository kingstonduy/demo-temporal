package domain

import (
	"context"
	"time"

	"go.temporal.io/sdk/workflow"
)

type CLientRequest struct {
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	Amount        int64  `json:"Amount"`
}

type ClientResponse struct {
	TransactionID string `json:"TransactionID"`
	FromAccountID string `json:"FromAccountID"`
	// FromAccountName string `json:"FromAccountName"`
	ToAccountID   string `json:"ToAccountID"`
	ToAccountName string `json:"ToAccountName"`
	Message       string `json:"Message"`
	Amount        int64  `json:"Amount"`
	Timestamp     string `json:"Timestamp"`
}

type WorkflowInput struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

type WorkflowOutput struct {
	TransactionID   string
	FromAccountID   string
	FromAccountName string
	ToAccountID     string
	ToAccountName   string
	Message         string
	Amount          int64
	Timestamp       string
}

type SaferRequest struct {
	WorkflowID    string `json:"WorkflowID"`
	RunID         string `json:"RunID"`
	Action        string `json:"Action"`
	TransactionID string `json:"TransactionID"`
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	Amount        int64  `json:"Amount"`
}

type SaferResponse struct {
	WorkflowID string `json:"WorkflowID"`
	RunID      string `json:"RunID"`
	Action     string `json:"SignalName"`
	Code       int    `json:"Code"`
	Status     string `json:"Status"`
	Message    string `json:"Message"`
}

type NapasAccountResponse struct {
	AccountID   string `json:"AccountID"`
	AccountName string `json:"AccountName"`
}

type AwaitSignals interface {
	Listen(ctx workflow.Context)
	GetNextTimeout(ctx workflow.Context) (time.Duration, error)
}

type TransactionEntity struct {
	TransactionID string `gorm:"primaryKey";column:transaction_id`
	FromAccountID string `gorm:"column:from_account_id"`
	// FromAccountName string `gorm:"column:from_account_name"`
	ToAccountID   string `gorm:"column:to_account_id"`
	ToAccountName string `gorm:"column:to_account_name"`
	Message       string `gorm:"column:message"`
	Amount        int64  `gorm:"column:amount;type:bigint""`
	Timestamp     string `gorm:"column:timestamp"`
	State         string `gorm:"column:state"`
}

func (*TransactionEntity) TableName() string {
	return "transaction"
}

type MoneyTransferHandler interface {
	Handle(ctx context.Context, request *WorkflowInput) (*WorkflowOutput, error)
}

type MoneyTransferRepository interface {
	Save(entity TransactionEntity) error
	Create(entity TransactionEntity) error
}

type MoneyTransferActivities interface {
	ValidateAccount(ctx context.Context, input SaferRequest) error
	LimitCut(ctx context.Context, input SaferRequest) error
	LimitCutCompensate(ctx context.Context, input SaferRequest) error
	MoneyCut(ctx context.Context, input SaferRequest) error
	MoneyCutCompensate(ctx context.Context, input SaferRequest) error
	UpdateMoney(ctx context.Context, input SaferRequest) error
	UpdateMoneyCompensate(ctx context.Context, input SaferRequest) error
	UpdateStateCreated(ctx context.Context, input TransactionEntity) error
	UpdateStateLimitCut(ctx context.Context, input TransactionEntity) error
	UpdateStateMoneyCut(ctx context.Context, input TransactionEntity) error
	UpdateStateTransactionCompleted(ctx context.Context, input TransactionEntity) error
	CompensateTransaction(ctx context.Context, input TransactionEntity) error
}
