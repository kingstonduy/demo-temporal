package main

import (
	"context"
	"log"

	traditionalway_getting_block "traditionalway-getting-block"
	"traditionalway-getting-block/shared"

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
		ID:        "temporal-demo-workflow-tranditionalway-getting-block_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	_, err = c.ExecuteWorkflow(context.Background(), option1, traditionalway_getting_block.BlockingWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute async workflow", err)
	}
}
