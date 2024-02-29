package workflow

import (
	"fmt"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MoneyTransferService(ctx workflow.Context, input model.WorkflowInput) (output model.WorkflowOutput, err error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 0,
			InitialInterval: time.Second * 5},
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var compensations Compensations

	defer func() {
		if err != nil {
			// activity failed, and workflow context is canceled
			disconnectedCtx, _ := workflow.NewDisconnectedContext(ctx)
			compensations.Compensate(disconnectedCtx, true)
		}
	}()

	fmt.Printf("💡Workflow input %+v\n", input)

	saferRequest := model.SaferRequest{
		WorkflowID:    workflow.GetInfo(ctx).WorkflowExecution.ID,
		RunID:         workflow.GetInfo(ctx).WorkflowExecution.RunID,
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
	}

	var napasAccountRes model.NapasAccountResponse
	err = workflow.ExecuteActivity(ctx, ValidateAccount, saferRequest).Get(ctx, &napasAccountRes)
	if err != nil {
		return output, err
	}

	var transactionEntity = model.TransactionEntity{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountRes.AccountName,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Amount:        input.Amount,
		Timestamp:     time.Now().String(),
		State:         "",
	}

	output = model.WorkflowOutput{
		TransactionID: input.TransactionID,
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		ToAccountName: napasAccountRes.AccountName,
		Amount:        input.Amount,
		Message:       fmt.Sprintf("Transfering money from %s to %s", input.FromAccountID, input.ToAccountID),
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
	}

	fmt.Printf("💡Output: %+v\n", output)

	err = workflow.ExecuteActivity(ctx, UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(CompensateTransaction, transactionEntity)

	err = workflow.ExecuteActivity(ctx, LimitCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(LimitCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}

	err = workflow.ExecuteActivity(ctx, MoneyCut, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(MoneyCutCompensate, saferRequest)
	err = workflow.ExecuteActivity(ctx, UpdateStateMoneyCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}

	err = workflow.ExecuteActivity(ctx, UpdateMoney, saferRequest).Get(ctx, nil)
	if err != nil {
		return output, err
	}
	compensations.AddCompensation(UpdateMoneyCompensate, saferRequest)

	err = workflow.ExecuteActivity(ctx, UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return output, err
	}

	return output, err

}
