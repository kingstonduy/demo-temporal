package main

import (
	"context"
	"log"
	"time"

	"versionning-build-id/shared"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type Person struct {
	Name string
}

func main() {
	// connect to temporal server
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	ctx := context.Background()
	taskQueue := shared.TaskQueueName

	// set up version 1--------------------------------------------------------------------------------
	err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
		TaskQueue: taskQueue,
		Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
			BuildID: "1.0",
		},
	})
	// err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
	// 	TaskQueue: taskQueue,
	// 	Operation: &client.BuildIDOpAddNewCompatibleVersion{
	// 		BuildID:                   "1.0",
	// 		ExistingCompatibleBuildID: "1.0",
	// 	},
	// })
	if err != nil {
		log.Fatalln("Unable to update worker build id compatibility", err)
	}

	firstWorkflowID := "build-id-versioning-first_" + uuid.New()
	firstWorkflowOptions := client.StartWorkflowOptions{
		ID:                       firstWorkflowID,
		TaskQueue:                taskQueue,
		WorkflowExecutionTimeout: 5 * time.Minute,
	}
	_, err = c.ExecuteWorkflow(ctx, firstWorkflowOptions, shared.WorkflowName)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}

	time.Sleep(5 * time.Second)

	// set up version 2--------------------------------------------------------------------------------
	err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
		TaskQueue: taskQueue,
		Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
			BuildID: "2.0",
		},
	})
	// err = c.UpdateWorkerBuildIdCompatibility(ctx, &client.UpdateWorkerBuildIdCompatibilityOptions{
	// 	TaskQueue: taskQueue,
	// 	Operation: &client.BuildIDOpAddNewCompatibleVersion{
	// 		BuildID:                   "2.0",
	// 		ExistingCompatibleBuildID: "2.0",
	// 	},
	// })
	if err != nil {
		log.Fatalln("Unable to update build id compatability", err)
	}

	// Start a new workflow, note that it will run on the new 2.0 version, without the client
	// invocation changing at all!
	secondWorkflowID := "build-id-versioning-second_" + uuid.New()
	secondWorkflowOptions := client.StartWorkflowOptions{
		ID:                       secondWorkflowID,
		TaskQueue:                taskQueue,
		WorkflowExecutionTimeout: 5 * time.Minute,
	}
	_, err = c.ExecuteWorkflow(ctx, secondWorkflowOptions, shared.WorkflowName)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}

	// time.Sleep(5 * time.Second)
}
