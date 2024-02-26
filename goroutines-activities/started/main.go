package main

import (
	"context"
	"log"

	goroutines_activities "goroutines-activities"
	"goroutines-activities/shared"

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
		ID:        "temporal-demo-workflow-goroutines_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	_, err = c.ExecuteWorkflow(context.Background(), option1, goroutines_activities.AsyncWorkFlow1)
	if err != nil {
		log.Fatalln("Unable to execute async workflow", err)
	}
}
