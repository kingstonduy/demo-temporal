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

func (*TransactionEntity) TableName() string {
	return "transaction"
}

type TransactionRepository interface {
	Create(c context.Context, transaction TransactionEntity) error
	DeleteById(c context.Context, transaction TransactionEntity) error
	Update(c context.Context, transaction TransactionEntity) error
	GetByID(c context.Context, id string) (TransactionEntity, error)
}
