package main

import (
	"context"
	"log"
	"os"

	"demo-temporal/shared"
	workflow "demo-temporal/workflow"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow",
		TaskQueue: shared.TaskQueueName,
	}

	if len(os.Args) <= 1 {
		log.Fatalln("Must specify name and language code as command-line arguments")
	}

	if os.Args[1] == "parallel" {
		// we, err := c.ExecuteWorkflow(context.Background(), options, workflow.ParallelWorkFlow, os.Args[2])
	} else {
		account, err := c.ExecuteWorkflow(context.Background(), options, workflow.AsyncWorkFlow)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}

		log.Println("Workflow completed", "ID", account)
	}

}
