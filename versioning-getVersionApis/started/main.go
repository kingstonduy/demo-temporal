package main

import (
	"context"
	"log"
	"time"

	versioning_getVersionApis "kingstonduy/demo-temporal/versioning-getVersionApis"
	"kingstonduy/demo-temporal/versioning-getVersionApis/shared"

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
		ID:        "temporal-demo-simple-workflow1_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	option2 := client.StartWorkflowOptions{
		ID:        "temporal-demo-simple-workflow2_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	_, _ = c.ExecuteWorkflow(context.Background(), option1, versioning_getVersionApis.SimpleWorkflow)

	time.Sleep(time.Second * 20)

	_, _ = c.ExecuteWorkflow(context.Background(), option2, versioning_getVersionApis.SimpleWorkflow)
}
