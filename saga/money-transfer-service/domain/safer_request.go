package domain

type SaferRequest struct {
	TransactionId string `json:"transactionId"`
	AccountId     string `json:"accountId"`
	Amount        int64  `json:"amount"`
}
