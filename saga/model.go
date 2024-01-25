package shared

type Compensations struct {
	compensations []any
	arguments     [][]any
}

type TransactionInfo struct {
	TransactionId string `json:"transactionId"`
	FromAccountId string `json:"fromAccountId"`
	ToAccountId   string `json:"toAccountId"  `
	Amount        int64  `json:"amount"`
}

type TransactionEntity struct {
	TransactionId string `gorm:"primaryKey";column:transaction_id`
	FromAccountId string `gorm:"column:from_account_id"`
	ToAccountId   string `gorm:"column:to_account_id"`
	Amount        int64  `gorm:"column:amount;type:bigint""`
	State         string `gorm:"column:state"`
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

type ValidateAccountInput struct {
	AccountId string `json:"accountId"`
}

type ValidateAccountOutput struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
	Amount      int64  `json:"amount"`
}

type SaferRequest struct {
	TransactionId string `json:"transactionId"`
	AccountId     string `json:"accountId"`
	Amount        int64  `json:"amount"`
}

type SaferResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
