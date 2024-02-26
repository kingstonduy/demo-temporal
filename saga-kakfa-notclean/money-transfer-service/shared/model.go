package model

type CLientRequest struct {
	FromAccount string `json:"FromAccount"`
	ToAccount   string `json:"ToAccount"`
	Amount      int64  `json:"Amount"`
}

type WorkflowInput struct {
	TransactionID string
	FromAccount   string
	ToAccount     string
	Amount        int64
}

type ClientResponse struct {
	TransactionID   string `json:"TransactionID"`
	FromAccountId   string `json:"FromAccountId"`
	FromAccountName string `json:"FromAccountName"`
	ToAccountId     string `json:"ToAccountId"`
	ToAccountName   string `json:"ToAccountName"`
	Message         string `json:"Message"`
	Amount          int64  `json:"Amount"`
	Timestamp       string `json:"Timestamp"`
}

type SaferRequest struct {
	WorkflowID string `json:"WorkflowID"`
	RunID      string `json:"RunID"`
	AccountID  string `json:"AccountID"`
	Amount     int64  `json:"Amount"`
}

type SaferResponse struct {
	WorkflowID string `json:"WorkflowID"`
	RunID      string `json:"RunID"`
	Code       int    `json:"Code"`
	Status     string `json:"Status"`
}

type NapasAccountResponse struct {
	AccountID   string `json:"AccountID"`
	AccountName string `json:"AccountName"`
}

type TransactionEntity struct {
	TransactionId   string `gorm:"primaryKey";column:transaction_id`
	FromAccountId   string `gorm:"column:from_account_id"`
	FromAccountName string `gorm:"column:from_account_name"`
	ToAccountId     string `gorm:"column:to_account_id"`
	ToAccountName   string `gorm:"column:to_account_name"`
	Message         string `gorm:"column:message"`
	Amount          int64  `gorm:"column:amount;type:bigint""`
	Timestamp       string `gorm:"column:timestamp"`
	State           string `gorm:"column:state"`
}

func (*TransactionEntity) TableName() string {
	return "transaction"
}

type T24Entity struct {
	AccountId string `gorm:"primaryKey";column:account_id`
	Amount    int64  `gorm:"column:amount;type:bigint""`
}

func (T24Entity) TableName() string {
	return "t24"
}

type AccountLimitEntity struct {
	AccountId string `gorm:"primaryKey";column:account_id`
	Amount    int64  `gorm:"column:amount;type:bigint""`
}

func (AccountLimitEntity) TableName() string {
	return "limit_manage"
}

type NapasEntity struct {
	AccountId   string `gorm:"primaryKey";column:account_id`
	AccountName string `gorm:"column:account_name"`
	Amount      int64  `gorm:"column:amount;type:bigint""`
}

func (NapasEntity) TableName() string {
	return "napas"
}
