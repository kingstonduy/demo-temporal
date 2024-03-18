package account

import (
	"context"

	"github.com/ocb/mcs-money-transfer/domain/account"
)

type checkBalanceHandler struct {
	accDb account.AccountBalanceData
}

func NewCheckBalanceHandler(
	accDb account.AccountBalanceData,
) account.CheckBalanceHandler {
	return &checkBalanceHandler{
		accDb: accDb,
	}
}

func (h *checkBalanceHandler) Handle(ctx context.Context, request *account.CheckBalanceRequest) (*account.CheckBalanceResponse, error) {
	// Business logic 1

	// Business logic 2

	balRes, err := h.accDb.GetBalance(ctx, request.Account)
	if err != nil {
		return nil, err
	}

	return &account.CheckBalanceResponse{
		Balance:  balRes.WorkingBalance,
		Currency: balRes.Currency,
	}, nil

}
