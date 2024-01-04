package main

import (
	"context"
	"demo-temporal/shared"
	"demo-temporal/workflow"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/pborman/uuid"
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

	option1 := client.StartWorkflowOptions{
		ID:        "temporal-demo-simple-workflow1_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	option2 := client.StartWorkflowOptions{
		ID:        "temporal-demo-simple-workflow2_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	_, _ = c.ExecuteWorkflow(context.Background(), option1, workflow.SimpleWorkflow)

	time.Sleep(time.Second * 20)

	_, _ = c.ExecuteWorkflow(context.Background(), option2, workflow.SimpleWorkflow)
}
