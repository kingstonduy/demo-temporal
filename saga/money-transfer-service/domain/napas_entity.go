package domain

type NapasEntity struct {
	AccountId   string `gorm:"primaryKey";column:account_id`
	AccountName string `gorm:"column:account_name"`
	Amount      int64  `gorm:"column:amount;type:bigint""`
}

func (NapasEntity) TableName() string {
	return "napas"
}
