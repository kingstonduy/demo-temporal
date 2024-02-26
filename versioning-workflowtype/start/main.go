package main

import (
	"context"
	"log"
	versioning_workflowtype "versioning-workflowtype"
	"versioning-workflowtype/shared"

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
	_, _ = c.ExecuteWorkflow(context.Background(), option1, versioning_workflowtype.SimpleWorkflow)

	// option2 := client.StartWorkflowOptions{
	// 	ID:        "temporal-demo-simple-workflowv2_" + uuid.New(),
	// 	TaskQueue: shared.TaskQueueName,
	// }
	// _, _ = c.ExecuteWorkflow(context.Background(), option2, versioning_workflowtype.SimpleWorkflow1)
}
