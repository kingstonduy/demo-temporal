package usecase

import (
	"fmt"
	"orchestrator-service/domain"
	"orchestrator-service/pkg/logger"
	temporal1 "orchestrator-service/usecase/temporal"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MoneyTransferWorkflow(
	ctx workflow.Context,
	input domain.WorkflowInput,
	log logger.Logger,
	compensations temporal1.Compensations,
	moneyTransferActivities MoneyTransferActivities) (output domain.WorkflowOutput, err error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
			InitialInterval: time.Second * 5},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	defer func() {
		if err != nil {
			// activity failed, and workflow context is canceled
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			compensations.Compensate(disconnectedCtx, true)
		}
	}()

	fmt.Printf("💡Workflow input %+v\n", input)

	saferRequest := domain.SaferRequest{
		WorkflowID:    workflow.GetInfo(ctx).WorkflowExecution.ID,
		RunID:         workflow.GetInfo(ctx).WorkflowExecution.RunID,
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
	}

	var napasAccountRes domain.NapasAccountResponse
	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.ValidateAccount, saferRequest).Get(ctx, &napasAccountRes)
	if err != nil {
		return output, err
	}

	var transactionEntity = domain.TransactionEntity{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountRes.AccountName,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Amount:        input.Amount,
		Timestamp:     time.Now().String(),
		State:         "",
	}

	output = domain.WorkflowOutput{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountRes.AccountName,
		Amount:        input.Amount,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
	}

	fmt.Printf("💡Output: %+v\n", output)

	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(moneyTransferActivities.CompensateTransaction, transactionEntity)

	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.LimitCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(moneyTransferActivities.LimitCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}

	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.MoneyCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(moneyTransferActivities.MoneyCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.UpdateStateMoneyCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}

	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.UpdateMoney, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(moneyTransferActivities.UpdateMoneyCompensate, saferRequest)

	err = workflow.ExecuteActivity(ctx, moneyTransferActivities.UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}

	return output, err
}
