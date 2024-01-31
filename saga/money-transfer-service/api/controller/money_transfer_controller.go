package controller

import (
	"context"
	shared "kingstonduy/demo-temporal/saga"
	app "kingstonduy/demo-temporal/saga/money-transfer-service"
	"kingstonduy/demo-temporal/saga/money-transfer-service/bootstrap"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type MoneyTransferController struct {
	MoneyTransferUsecase domain.MoneyTransferUsecase
	Env                  *bootstrap.Env
}

func (lc *MoneyTransferController) MoneyTransfer(c client.Client) gin.HandlerFunc {
	fn := func(g *gin.Context) {
		var transferInfo domain.TransactionInfo
		err := g.BindJSON(&transferInfo)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
			return
		}

		log.Printf("💡Request %+v\n", transferInfo)

		option := client.StartWorkflowOptions{
			ID:        lc.Env.Workflow + "_" + time.Now().String(),
			TaskQueue: lc.Env.TaskQueue,
		}
		transferInfo.TransactionId = option.ID
		_, err = c.ExecuteWorkflow(context.Background(), option, app.MoneyTransferWorkflow, transferInfo)
		if err != nil {
			log.Fatalf("Unable to execute %s workflow\n, error=%s", option.ID, err)
		}

		// err = we.Get(context.Background(), nil)
		// if err != nil {
		// 	g.JSON(http.StatusBadRequest, gin.H{
		// 		"message": err.Error(),
		// 	})
		// } else {
		// 	g.JSON(http.StatusAccepted, gin.H{
		// 		"message": "ok",
		// 	})
		// }
	}
	return gin.HandlerFunc(fn)
}

func MoneyTransferWorkflow(ctx workflow.Context, usecase1 domain.MoneyTransferUsecase, info shared.TransactionInfo) (err error) {
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
	err = workflow.ExecuteActivity(ctx, usecase1.ValidateAccount, info).Get(ctx, nil)
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
