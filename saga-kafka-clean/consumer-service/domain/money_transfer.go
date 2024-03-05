package domain

type SaferRequest struct {
	WorkflowID    string `json:"WorkflowID"`
	RunID         string `json:"RunID"`
	Action        string `json:"Action"`
	TransactionID string `json:"TransactionID"`
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	Amount        int64  `json:"Amount"`
}

type SaferResponse struct {
	WorkflowID string `json:"WorkflowID"`
	RunID      string `json:"RunID"`
	Action     string `json:"SignalName"`
	Code       int    `json:"Code"`
	Status     string `json:"Status"`
	Message    string `json:"Message"`
}

type MoneyTransferHandler interface {
	Handle(request SaferResponse) error
}
