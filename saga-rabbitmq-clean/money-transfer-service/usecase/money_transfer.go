package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/ocb/mcs-money-transfer/domain"
)

type moneyTransferHandler struct {
}

func NewMoneyTransferHandler() domain.MoneyTransferHandler {
	return &moneyTransferHandler{}
}

// Handle implements domain.MoneyTransferHandler.
func (m *moneyTransferHandler) Handle(ctx context.Context, request *domain.MoneyTransferClientRequest) (*domain.MoneyTransferClientResponse, error) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	return &domain.MoneyTransferClientResponse{
		TransactionID: uuid.New().String(),
		FromAccountID: request.FromAccountID,
		ToAccountID:   request.ToAccountID,
		ToAccountName: nameGenerator.Generate(),
		Message:       "Transfer success",
		Amount:        request.Amount,
		Timestamp:     time.RFC1123,
	}, nil
}
