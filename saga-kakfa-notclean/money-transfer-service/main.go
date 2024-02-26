package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kingstonduy/demo-temporal/saga-kakfa-notclean/money-transfer-service/config"
	model "kingstonduy/demo-temporal/saga-kakfa-notclean/money-transfer-service/shared"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
			FromAccount:   clientReq.FromAccount,
			ToAccount:     clientReq.ToAccount,
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

// consume from kafka then signal the workflow to continue
func Consume(cl client.Client, bootstrapServer string, topic string) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var saferResponse model.SaferResponse
			err := json.Unmarshal([]byte(string(msg.Value)), &saferResponse)
			if err != nil {
				log.Println(err)
				return
			}

			signalName := saferResponse.WorkflowID + "-" + saferResponse.RunID
			err = cl.SignalWorkflow(context.Background(), saferResponse.WorkflowID, saferResponse.RunID, signalName, nil)
			if err != nil {
				log.Println(err)
				return
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
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

	go func() {
		Consume(c,
			fmt.Sprintf("%s:%s", config.Kafka.BootstrapServer.Host, config.Kafka.BootstrapServer.Port),
			"money_transfer_reply_channel")
	}()

	g.Run(fmt.Sprintf("%s:%s", config.MoneyTransfer.Host, config.MoneyTransfer.Port))
}
