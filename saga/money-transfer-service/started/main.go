package main

import (
	"context"
	shared "kingstonduy/demo-temporal/saga"
	app "kingstonduy/demo-temporal/saga/money-transfer-service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func validate(info shared.TransactionInfo) bool {
	return info.Amount >= 0
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	router := gin.Default()
	router.POST("/api/v1/moneytransfer", func(g *gin.Context) {

		var transferInfo shared.TransactionInfo
		err := g.BindJSON(&transferInfo)
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
			return
		}
		if !validate(transferInfo) {
			g.JSON(http.StatusBadRequest, gin.H{
				"message": "The amount of money must be positive",
			})
			return
		}

		log.Printf("💡Request %+v\n", transferInfo)

		option := client.StartWorkflowOptions{
			ID:        shared.WORKFLOW + "_" + time.Now().String(),
			TaskQueue: shared.TASKQUEUE,
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
	})

	router.Run(shared.MONEY_TRANSFER_SERVICE_HOST_PORT)
}
