package temporal

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
			MaximumAttempts: 5,
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
	err = workflow.ExecuteActivity(ctx, usecase.ValidateAccount, info).Get(ctx, nil)
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
	err = workflow.ExecuteActivity(ctx, UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(CompensateTransaction, &shared.TransactionEntity{
			TransactionId: info.TransactionId,
			FromAccountId: info.FromAccountId,
			ToAccountId:   info.ToAccountId,
			Amount:        info.Amount,
			State:         "CANCELLED",
		})
	}

	// tru han muc giao dich
	err = workflow.ExecuteActivity(ctx, LimitCut, info).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(LimitCutCompensate, info)
	}

	// ghi vao db trang thai LIMIT_CUT
	transactionEntity.State = "LIMIT_CUT"
	err = workflow.ExecuteActivity(ctx, UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// goi t24 cat tien tai khoan ocb
	err = workflow.ExecuteActivity(ctx, MoneyCut, info).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(MoneyCutCompensate, info)
	}

	// ghi vao db trang thai MONEY_CUT
	transactionEntity.State = "MONEY_CUT"
	err = workflow.ExecuteActivity(ctx, UpdateStateMoneyCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// goi napas ghi co vao tai khoan thu huong
	err = workflow.ExecuteActivity(ctx, UpdateMoney, info).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(UpdateMoneyCompensate, info)
	}

	// ghi vao db trang thai COMPLETE
	transactionEntity.State = "COMPLETED"
	err = workflow.ExecuteActivity(ctx, UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
