package domain

import "context"

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

type TransactionRepository interface {
	Create(c context.Context, transaction *TransactionInfo) error
	Fetch(c context.Context) ([]TransactionInfo, error)
	GetByEmail(c context.Context, email string) (TransactionInfo, error)
	GetByID(c context.Context, id string) (TransactionInfo, error)
}
