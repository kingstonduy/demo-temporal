package domain

import "context"

type MoneyTransferClientRequest struct {
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	Amount        int64  `json:"Amount"`
}

type MoneyTransferClientResponse struct {
	TransactionID string `json:"TransactionID"`
	FromAccountID string `json:"FromAccountID"`
	ToAccountID   string `json:"ToAccountID"`
	ToAccountName string `json:"ToAccountName"`
	Message       string `json:"Message"`
	Amount        int64  `json:"Amount"`
	Timestamp     string `json:"Timestamp"`
}

type MoneyTransferWorkflowInput struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

type MoneyTransferWorkflowOutput struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	ToAccountName string
	Message       string
	Amount        int64
	Timestamp     string
}

type MoneyTransferHandler interface {
	Handle(ctx context.Context, request *MoneyTransferClientRequest) (*MoneyTransferClientResponse, error)
}
