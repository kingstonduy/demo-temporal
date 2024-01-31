package domain

type AccountLimitEntity struct {
	AccountId string `gorm:"primaryKey";column:account_id`
	Amount    int64  `gorm:"column:amount;type:bigint""`
}

func (AccountLimitEntity) TableName() string {
	return "limit_manage"
}
