package domain

import "context"

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
	// from service
	// ts send message
	// request id
	// trace id = khoi tao workflow
	// span id bind each service
	WorkflowID    string `json:"WorkflowID"`
	RunID         string `json:"RunID"`
	TransactionID string `json:"TransactionID"`
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	Amount        int64  `json:"Amount"`
}

type SaferResponse struct {
	// ts reply
	// duration
	// header request = header response
	WorkflowID string `json:"WorkflowID"`
	RunID      string `json:"RunID"`
	Code       int    `json:"Code"`
	Status     string `json:"Status"`
	Message    string `json:"Message"`
}

type NapasAccountResponse struct {
	AccountID   string `json:"AccountID"`
	AccountName string `json:"AccountName"`
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

type MoneyTransferRepository interface {
	Save(entity TransactionEntity) error
	Create(entity TransactionEntity) error
}

type MoneyTransferMessageQueue interface {
	SaferRequestResponse(input SaferRequest, queue string) (output SaferResponse, err error)
}

type MoneyTransferActivities interface {
	ValidateAccount(ctx context.Context, input SaferRequest) (output NapasAccountResponse, err error)
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
