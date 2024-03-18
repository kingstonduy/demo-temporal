package usecase

import (
	"context"

	"github.com/lengocson131002/go-clean-core/logger"
	"github.com/ocb/mcs-money-transfer/bootstrap"
	"github.com/ocb/mcs-money-transfer/domain"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type moneyTransferHandler struct {
	client *client.Client
	log    logger.Logger
	cfg    *bootstrap.Config
}

func NewMoneyTransferHandler(
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
func (h *moneyTransferHandler) Handle(ctx context.Context, req *domain.MoneyTransferClientRequest) (*domain.MoneyTransferClientResponse, error) {
	// seed := time.Now().UTC().UnixNano()
	// nameGenerator := namegenerator.NewNameGenerator(seed)

	// return &domain.MoneyTransferClientResponse{
	// 	TransactionID: uuid.New().String(),
	// 	FromAccountID: request.FromAccountID,
	// 	ToAccountID:   request.ToAccountID,
	// 	ToAccountName: nameGenerator.Generate(),
	// 	Message:       "Transfer success",
	// 	Amount:        request.Amount,
	// 	Timestamp:     time.RFC1123,
	// }, nil

	var res = domain.MoneyTransferClientResponse{}

	var workflowInput = &domain.MoneyTransferWorkflowInput{
		TransactionID: uuid.New(),
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	options := client.StartWorkflowOptions{
		ID:        h.cfg.Temporal.Workflow + "-" + workflowInput.TransactionID,
		TaskQueue: h.cfg.Temporal.TaskQueue,
	}

	we, err := (*h.client).ExecuteWorkflow(context.Background(), options, h.cfg.Temporal.Workflow, workflowInput)
	if err != nil {
		return &domain.MoneyTransferClientResponse{}, err
	}

	err = we.Get(context.Background(), &res)
	if err != nil {
		return &domain.MoneyTransferClientResponse{}, err
	}

	return &res, nil
}
