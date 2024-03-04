package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"go.temporal.io/sdk/client"
)

type moneyTransferHandler struct {
	client *client.Client
	log    logger.Logger
	cfg    *bootstrap.Config
}

func NewMoneyTransferHandler(

	// mdRepo data.MasterDataRepository,
	client *client.Client,
	cfg *bootstrap.Config,
	logger logger.Logger,

) domain.MoneyTransferHandler {

	return &moneyTransferHandler{
		client: client,
		log:    logger,
		cfg:    cfg,
	}
}

// Handle implements domain.MoneyTransferHandler.
func (*moneyTransferHandler) Handle(ctx context.Context, req *domain.WorkflowInput) (*domain.WorkflowOutput, error) {
	panic("unimplemented")
}
