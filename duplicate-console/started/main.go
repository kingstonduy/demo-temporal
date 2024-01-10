package main

import (
	"context"
	"log"
	"os"

	duplicate_console "kingstonduy/demo-temporal/duplicate-console"
	"kingstonduy/demo-temporal/duplicate-console/model"
	"kingstonduy/demo-temporal/duplicate-console/shared"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	optionsParallel := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow-parallel_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	optionsAsync := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow-async_" + uuid.New(),
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

		_, err := c.ExecuteWorkflow(context.Background(), optionsParallel, duplicate_console.ParallelWorkFlow, input)
		if err != nil {
			log.Fatalln("Unable to execute parallel workflow", err)
		}

		break

	case "async":
		_, err := c.ExecuteWorkflow(context.Background(), optionsAsync, duplicate_console.AsyncWorkFlow)
		if err != nil {
			log.Fatalln("Unable to execute async workflow", err)
		}

		break

	default:
		input := model.ParallelWorkflowInput{
			Cif1: "1",
			Cif2: "2",
		}

		_, err := c.ExecuteWorkflow(context.Background(), optionsAsync, duplicate_console.AsyncWorkFlow)
		if err != nil {
			log.Fatalln("Unable to execute async workflow", err)
		}

		_, err = c.ExecuteWorkflow(context.Background(), optionsParallel, duplicate_console.ParallelWorkFlow, input)
		if err != nil {
			log.Fatalln("Unable to execute parallel workflow", err)
		}

		log.Println("BOTH")
		break

	}

}
