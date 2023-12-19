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

	optionsParallel := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow-parallel",
		TaskQueue: shared.TaskQueueName,
	}

	optionsAsync := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow-async",
		TaskQueue: shared.TaskQueueName,
	}

	arg := "both"
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	switch arg {
	case "parallel":
		input := model.ParallelWorkflowInput{
			Cif1: "1",
			Cif2: "2",
		}

		_, err := c.ExecuteWorkflow(context.Background(), optionsParallel, workflow.ParallelWorkFlow, input)
		if err != nil {
			log.Fatalln("Unable to execute parallel workflow", err)
		}

		break

	case "async":
		_, err := c.ExecuteWorkflow(context.Background(), optionsAsync, workflow.AsyncWorkFlow)
		if err != nil {
			log.Fatalln("Unable to execute async workflow", err)
		}

		break

	default:
		input := model.ParallelWorkflowInput{
			Cif1: "1",
			Cif2: "2",
		}

		_, err := c.ExecuteWorkflow(context.Background(), optionsAsync, workflow.AsyncWorkFlow)
		if err != nil {
			log.Fatalln("Unable to execute async workflow", err)
		}

		_, err = c.ExecuteWorkflow(context.Background(), optionsParallel, workflow.ParallelWorkFlow, input)
		if err != nil {
			log.Fatalln("Unable to execute parallel workflow", err)
		}

		log.Println("BOTH")
		break

	}

}
