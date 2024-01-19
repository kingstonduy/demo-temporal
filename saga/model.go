package shared

type Compensations struct {
	compensations []any
	arguments     [][]any
}

type TransactionInfo struct {
	TransactionId string
	FromAccountId string
	ToAccountId   string
	Amount        int
}

type TransactionEntity struct {
	TransactionId string `gorm:"primaryKey";column:transaction_id`
	FromAccountId string `gorm:"column:from_account_id"`
	ToAccountId   string `gorm:"column:to_account_id"`
	Amount        int    `gorm:"column:amount"`
	State         string `gorm:"column:state"`
}

type T24Entity struct {
	AccountId string `gorm:"primaryKey";column:account_id`
	Amount    int    `gorm:"column:amount"`
}

func (T24Entity) TableName() string {
	return "t24"
}

type AccountLimitEntity struct {
	AccountId string `gorm:"primaryKey";column:account_id`
	Amount    int    `gorm:"column:amount"`
}

func (AccountLimitEntity) TableName() string {
	return "limit_manage"
}

type NapasEntity struct {
	AccountId   string `gorm:"primaryKey";column:account_id`
	AccountName string `gorm:"column:account_name"`
	Amount      int    `gorm:"column:amount"`
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
	Amount      int    `json:"amount"`
}

type SaferRequest struct {
	TransactionId string `json:"transactionId"`
	AccountId     string `json:"accountId"`
	Amount        int    `json:"amount"`
}

type SaferResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
