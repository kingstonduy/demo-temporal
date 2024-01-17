package main

import (
	"context"
	"kingstonduy/demo-temporal/async"
	"kingstonduy/demo-temporal/async/shared"
	money_transfer_service "kingstonduy/demo-temporal/saga/money-transfer-service"

	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	var config money_transfer_service.Config
	config.GetInstance()

	optionsAsync := client.StartWorkflowOptions{
		ID:        config.MONEY_TRANSFER_SERVICE.WORKFLOW_NAME + "_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	_, err = c.ExecuteWorkflow(context.Background(), optionsAsync, async.AsyncWorkFlow)
	if err != nil {
		log.Fatalf("Unable to execute %s workflow\n", optionsAsync.ID, err)
	}
}
