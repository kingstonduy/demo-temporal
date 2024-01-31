package temporal

import (
	"fmt"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type MoneyTransferWorkflow struct {
	Usecase domain.MoneyTransferUsecase
}

func (w *MoneyTransferWorkflow) NewMoneyTransferWorkflow(ctx workflow.Context, input domain.TransactionInfo) (err error) {
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

	fmt.Printf("💡Info %+v\n", input)
	// just read the database, dont need to compensate
	err = workflow.ExecuteActivity(ctx, w.Usecase.ValidateAccount, input).Get(ctx, nil)
	if err != nil {
		return err
	}

	var transactionEntity = domain.TransactionEntity{
		TransactionId: input.TransactionId,
		FromAccountId: input.FromAccountId,
		ToAccountId:   input.ToAccountId,
		Amount:        input.Amount,
		State:         "CREATED",
	}

	// ghi vao db trang thai CREATED
	err = workflow.ExecuteActivity(ctx, w.Usecase.UpdateStateCreated, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(w.Usecase.CompensateTransaction, &domain.TransactionEntity{
			TransactionId: input.TransactionId,
			FromAccountId: input.FromAccountId,
			ToAccountId:   input.ToAccountId,
			Amount:        input.Amount,
			State:         "CANCELLED",
		})
	}

	// tru han muc giao dich
	err = workflow.ExecuteActivity(ctx, w.Usecase.LimitCut, input).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(w.Usecase.LimitCutCompensate, input)
	}

	// ghi vao db trang thai LIMIT_CUT
	transactionEntity.State = "LIMIT_CUT"
	err = workflow.ExecuteActivity(ctx, w.Usecase.UpdateStateLimitCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// goi t24 cat tien tai khoan ocb
	err = workflow.ExecuteActivity(ctx, w.Usecase.MoneyCut, input).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(w.Usecase.MoneyCutCompensate, input)
	}

	// ghi vao db trang thai MONEY_CUT
	transactionEntity.State = "MONEY_CUT"
	err = workflow.ExecuteActivity(ctx, w.Usecase.UpdateStateMoneyCut, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// goi napas ghi co vao tai khoan thu huong
	err = workflow.ExecuteActivity(ctx, w.Usecase.UpdateMoney, input).Get(ctx, nil)
	if err != nil {
		return err
	} else {
		compensations.AddCompensation(w.Usecase.UpdateMoneyCompensate, input)
	}

	// ghi vao db trang thai COMPLETE
	transactionEntity.State = "COMPLETED"
	err = workflow.ExecuteActivity(ctx, w.Usecase.UpdateStateTransactionCompleted, transactionEntity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return err
}
