package main

import (
	"context"
	shared "kingstonduy/demo-temporal/saga"
	app "kingstonduy/demo-temporal/saga/money-transfer-service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	router := gin.Default()
	router.POST("/api/v1/moneytransfer", func(g *gin.Context) {

		var transferInfo shared.TransactionInfo
		err := g.BindJSON(&transferInfo)
		if err != nil {
			g.JSON(400, gin.H{
				"message": "Invalid request",
			})
			return
		}

		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create client", err)
		}
		defer c.Close()

		option := client.StartWorkflowOptions{
			ID:        shared.WORKFLOW + "_" + uuid.New(),
			TaskQueue: shared.TASKQUEUE,
		}

		_, err = c.ExecuteWorkflow(context.Background(), option, app.MoneyTransferWorkflow, transferInfo)
		if err != nil {
			log.Fatalf("Unable to execute %s workflow\n, error=%s", option.ID, err)
		}
	})

	router.Run(shared.MONEY_TRANSFER_SERVICE_URL)
}
