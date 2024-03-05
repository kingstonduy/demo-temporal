package usecase

import (
	"encoding/json"
	"fmt"
	"orchestrator-service/domain"
	"orchestrator-service/infra/logger"
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
	moneyMoneyTransferActivities MoneyTransferActivities,
) (output domain.WorkflowOutput, err error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 0,
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

	var a AwaitSignals
	workflow.Go(ctx, a.Listen)

	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.ValidateAccount, saferRequest).Get(ctx, nil)
	if err != nil {
		return
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.NapasAccountSignal
	})
	// Cancellation
	if err != nil {
		return
	}
	if a.NapasAccountResponse.Code != 200 {
		return output, fmt.Errorf("NapasAccountResponse.Code: %d, NapasAccountResponse.Message: %s", a.NapasAccountResponse.Code, a.NapasAccountResponse.Message)
	}

	var napasAccountResponse domain.NapasAccountResponse
	err = json.Unmarshal([]byte(a.NapasAccountResponse.Message), &napasAccountResponse)
	if err != nil {
		return
	}

	var transactionEntity = domain.TransactionEntity{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountResponse.AccountName,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Amount:        input.Amount,
		Timestamp:     time.Now().String(),
		State:         "",
	}

	output = domain.WorkflowOutput{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountResponse.AccountName,
		Amount:        input.Amount,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
	}

	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(moneyMoneyTransferActivities.CompensateTransaction, transactionEntity)

	// -------------------------------------------------------------------------------------------------
	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.LimitCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.LimitSignal
	})
	// Cancellation
	if err != nil {
		return output, err
	}
	if a.LimitResponse.Code != 200 {
		return output, fmt.Errorf("LimitResponse.Code: %d, LimitResponse.Message: %s", a.LimitResponse.Code, a.LimitResponse.Message)
	}
	compensations.AddCompensation(moneyMoneyTransferActivities.LimitCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// -------------------------------------------------------------------------------------------------
	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.MoneyCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.T24Signal
	})
	// Cancellation
	if err != nil {
		return output, err
	}
	if a.T24Response.Code != 200 {
		return output, fmt.Errorf("T24Response.Code: %d, T24Response.Message: %s", a.T24Response.Code, a.T24Response.Message)
	}
	compensations.AddCompensation(moneyMoneyTransferActivities.MoneyCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.UpdateStateMoneyCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// -------------------------------------------------------------------------------------------------
	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.UpdateMoney, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// Wait for NapasAccountSignal
	err = workflow.Await(ctx, func() bool {
		return a.NapasMoneySignal
	})
	// Cancellation
	if err != nil {
		return output, err
	}
	if a.NapasMoneyResponse.Code != 200 {
		return output, fmt.Errorf("NapasMoneyResponse.Code: %d, NapasMoneyResponse.Message: %s", a.NapasMoneyResponse.Code, a.NapasMoneyResponse.Message)
	}
	compensations.AddCompensation(moneyMoneyTransferActivities.UpdateMoneyCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, moneyMoneyTransferActivities.UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	// -------------------------------------------------------------------------------------------------
	fmt.Printf("💡Workflow output", output)
	return output, err
}
