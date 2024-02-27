package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"saga-rabbitmq-notclean/money-transfer-service/config"
	model "saga-rabbitmq-notclean/money-transfer-service/shared"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

// request, response to client. starts a workflow then wait for the workflow to finish
func Handler(c client.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		var clientReq model.CLientRequest
		err := ctx.BindJSON(&clientReq)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, error.Error(err))
		}

		var workflowInput = &model.WorkflowInput{
			TransactionID: uuid.New(),
			FromAccountID: clientReq.FromAccountID,
			ToAccountID:   clientReq.ToAccountID,
			Amount:        clientReq.Amount,
		}

		options := client.StartWorkflowOptions{
			ID:        config.GetConfig().Temporal.Workflow + "-" + workflowInput.TransactionID,
			TaskQueue: config.GetConfig().Temporal.TaskQueue,
		}

		we, err := c.ExecuteWorkflow(context.Background(), options, config.GetConfig().Temporal.Workflow, workflowInput)
		if err != nil {
			ctx.JSON(500, error.Error(err))
			return
		}

		var clientResponse model.ClientResponse
		err = we.Get(ctx, &clientResponse)
		if err != nil {
			ctx.JSON(500, error.Error(err))
			return
		}

		ctx.JSON(200, clientResponse)
	}
	return fn
}

func main() {
	config := config.GetConfig()

	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%s", config.Temporal.Host, config.Temporal.Port),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	g := gin.Default()
	publicRouter := g.Group("/api/v1")
	publicRouter.POST("/moneytransfer", Handler(c))

	g.Run(fmt.Sprintf("%s:%s", config.MoneyTransfer.Host, config.MoneyTransfer.Port))
}
