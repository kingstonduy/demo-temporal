package controller

import (
	"context"
	"log"
	"net/http"
	"saga/money-transfer-service/bootstrap"
	"saga/money-transfer-service/domain"
	"saga/money-transfer-service/temporal"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
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

		workflow := temporal.MoneyTransferWorkflow{
			Usecase: lc.MoneyTransferUsecase,
		}

		log.Printf("WorkflowInput %+v\n", workflow)
		_, err = c.ExecuteWorkflow(context.Background(), option, workflow.NewMoneyTransferWorkflow, transferInfo)
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
