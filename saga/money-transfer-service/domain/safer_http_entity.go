package domain

type SaferRequest struct {
	TransactionId string `json:"transactionId"`
	AccountId     string `json:"accountId"`
	Amount        int64  `json:"amount"`
}

type SaferResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SaferHttpEntity interface {
	Post(url, in *SaferRequest, out *SaferResponse) error
}
