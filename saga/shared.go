package shared

const TaskQueue = "money-transfer-service-task-queue"
const Workflow = "money-transfer-service-workflow"

const LIMITATION_SERVICE_URL = "http://localhost:7202/"
const T24_SERVICE_URL = "http://localhost:7203/"
const NAPAS_SERVICE_URL = "http://localhost:7204/"

const POSTGRES_URL = "jdbc:postgresql://localhost:5432/postgres"
const POSTGRES_HOST = "localhost"
const POSTGRES_PORT = 5432
const POSTGRES_USER = "postgres"
const POSTGRES_PASSWORD = "changeme"
const POSTGRES_DBNAME = "postgres"

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
	TransactionId string
	FromAccountId string
	ToAccountId   string
	Amount        int
	State         string
}

type AccountLimitEntity struct {
	AccountId string
	limit     int
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
