package main

import (
	"context"
	"log"
	"os"

	"demo-temporal/model"
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
		input := model.ParallelWorkflowInput{
			Cif1: "1",
			Cif2: "2",
		}

		_, err := c.ExecuteWorkflow(context.Background(), options, workflow.ParallelWorkFlow, input)
		if err != nil {
			log.Fatalln("Unable to execute parallel workflow", err)
		}

	} else if os.Args[1] == "async" {
		_, err := c.ExecuteWorkflow(context.Background(), options, workflow.AsyncWorkFlow)
		if err != nil {
			log.Fatalln("Unable to execute async workflow", err)
		}
	} else {
		log.Fatal("Invalid Argument")
	}

}
