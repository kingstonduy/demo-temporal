package app

import (
	shared "kingstonduy/demo-temporal/saga"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func MoneyTransferWorkflow(ctx workflow.Context, info shared.TransactionInfo) (err error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
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

	// just read the database, dont need to compensate
	err = workflow.ExecuteActivity(ctx, ValidateAccount, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	var transactionEntity = shared.TransactionEntity{
		TransactionId: info.TransactionId,
		FromAccountId: info.FromAccountId,
		ToAccountId:   info.ToAccountId,
		Amount:        info.Amount,
		State:         "CREATED",
	}

	// ghi vao db trang thai CREATED
	compensations.AddCompensation(CompensateTransaction, &shared.TransactionEntity{
		TransactionId: info.TransactionId,
		FromAccountId: info.FromAccountId,
		ToAccountId:   info.ToAccountId,
		Amount:        info.Amount,
		State:         "CANCELLED",
	})
	err = workflow.ExecuteActivity(ctx, UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// tru han muc giao dich
	compensations.AddCompensation(LimitCutCompensate, info)
	err = workflow.ExecuteActivity(ctx, LimitCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	// ghi vao db trang thai LIMIT_CUT
	transactionEntity.State = "LIMIT_CUT"
	err = workflow.ExecuteActivity(ctx, UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// goi t24 cat tien tai khoan ocb
	compensations.AddCompensation(MoneyCutCompensate, info)
	err = workflow.ExecuteActivity(ctx, MoneyCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	// ghi vao db trang thai MONEY_CUT
	transactionEntity.State = "MONEY_CUT"
	err = workflow.ExecuteActivity(ctx, UpdateStateMoneyCut, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	// goi napas ghi co vao tai khoan thu huong
	compensations.AddCompensation(UpdateMoneyCompensate, info)
	err = workflow.ExecuteActivity(ctx, UpdateMoney, info).Get(ctx, nil)
	if err != nil {
		return err
	}

	// ghi vao db trang thai COMPLETE
	transactionEntity.State = "COMPLETED"
	err = workflow.ExecuteActivity(ctx, UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
