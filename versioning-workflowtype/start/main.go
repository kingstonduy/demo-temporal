package main

import (
	"context"
	"demo-temporal/shared"
	"demo-temporal/workflow"
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

	option1 := client.StartWorkflowOptions{
		ID:        "temporal-demo-versioning-workflowtype_ " + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}
	_, _ = c.ExecuteWorkflow(context.Background(), option1, workflow.SimpleWorkflow)

	// option2 := client.StartWorkflowOptions{
	// 	ID:        "temporal-demo-simple-workflowv2_" + uuid.New(),
	// 	TaskQueue: shared.TaskQueueName,
	// }
	// _, _ = c.ExecuteWorkflow(context.Background(), option2, workflow.SimpleWorkflow1)
}
