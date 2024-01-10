package main

import (
	"context"
	"kingstonduy/demo-temporal/async"
	"kingstonduy/demo-temporal/async/shared"

	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type Person struct {
	Name string
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	optionsAsync := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow-async_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	_, err = c.ExecuteWorkflow(context.Background(), optionsAsync, async.AsyncWorkFlow)
	if err != nil {
		log.Fatalln("Unable to execute async workflow", err)
	}
}
