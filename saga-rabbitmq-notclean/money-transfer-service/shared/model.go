package model

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
	TransactionID string `json:"TransactionID"`
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	Amount        int64  `json:"Amount"`
}

type SaferResponse struct {
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

type T24Entity struct {
	AccountID string `gorm:"primaryKey";column:account_id`
	Amount    int64  `gorm:"column:amount;type:bigint""`
}

func (T24Entity) TableName() string {
	return "t24"
}

type AccountLimitEntity struct {
	AccountID string `gorm:"primaryKey";column:account_id`
	Amount    int64  `gorm:"column:amount;type:bigint""`
}

func (AccountLimitEntity) TableName() string {
	return "limit_manage"
}

type NapasEntity struct {
	AccountID   string `gorm:"primaryKey";column:account_id`
	AccountName string `gorm:"column:account_name"`
	Amount      int64  `gorm:"column:amount;type:bigint""`
}

func (NapasEntity) TableName() string {
	return "napas"
}
