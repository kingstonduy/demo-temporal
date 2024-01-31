package domain

type T24Entity struct {
	AccountId string `gorm:"primaryKey";column:account_id`
	Amount    int64  `gorm:"column:amount;type:bigint""`
}

func (T24Entity) TableName() string {
	return "t24"
}
