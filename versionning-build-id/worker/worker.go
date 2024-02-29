package main

import (
	"log"
	"sync"

	versionning_build_id "versionning-build-id"
	"versionning-build-id/shared"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	taskQueue := shared.TaskQueueName
	log.Println("Using Task Queue name: ", taskQueue, "(Copy this!)")

	wg := sync.WaitGroup{}

	// tao worker 1 id = 1.0, thuc hien workflow1
	createAndRunWorker(c, taskQueue, "1.0", versionning_build_id.SimpleWorkflowV1, &wg)
	// tao worker 1 id = 2.0, thuc hien workflow2
	createAndRunWorker(c, taskQueue, "2.0", versionning_build_id.SimpleWorkflowV2, &wg)
	// tao worker 1 id = 3.0, thuc hien workflow3
	createAndRunWorker(c, taskQueue, "3.0", versionning_build_id.SimpleWorkflowV3, &wg)
	wg.Wait()
}

func createAndRunWorker(c client.Client, taskQueue string, buildID string, workflowFunc func(ctx workflow.Context) error, wg *sync.WaitGroup) {
	// create worker based  on buidlID, provide a unique identifier for a set of worker code
	w := worker.New(c, taskQueue, worker.Options{
		BuildID:                 buildID,
		UseBuildIDForVersioning: true,
	})

	// It's important that we register all the different implementations of the workflow using
	// the same name. This allows us to demonstrate what would happen if you were making changes
	// to this workflow code over time while keeping the same workflow name/type.
	w.RegisterWorkflowWithOptions(workflowFunc, workflow.RegisterOptions{Name: shared.WorkflowName})

	w.RegisterActivity(versionning_build_id.GetInformation)
	w.RegisterActivity(versionning_build_id.GetInformation1)
	w.RegisterActivity(versionning_build_id.GetInformation2)

	// run the worker
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	}()
}
