package account

import "context"

// USECASE
type CheckBalanceRequest struct {
	Account string `json:"account"`
}

type CheckBalanceResponse struct {
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type CheckBalanceHandler interface {
	Handle(ctx context.Context, request *CheckBalanceRequest) (*CheckBalanceResponse, error)
}

// INTERFACES
type AccountBalanceResponse struct {
	Currency        string
	OpenActualBal   int64
	OnlineActualBal int64
	WorkingBalance  int64
}

// Get account balance using account number
type AccountBalanceData interface {
	GetBalance(cxt context.Context, accNumber string) (*AccountBalanceResponse, error)
}
