package domain

type ValidateAccountInput struct {
	AccountId string `json:"accountId"`
}

type ValidateAccountOutput struct {
	AccountId   string `json:"accountId"`
	AccountName string `json:"accountName"`
	Amount      int64  `json:"amount"`
}
