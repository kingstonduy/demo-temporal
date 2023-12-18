package main

import (
	"demo-temporal/activity"
	"demo-temporal/shared"
	"demo-temporal/workflow"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, shared.TaskQueueName, worker.Options{})

	w.RegisterWorkflow(workflow.ParallelWorkFlow)
	w.RegisterWorkflow(workflow.AsyncWorkFlow)

	w.RegisterActivity(activity.RegisterAccount)
	w.RegisterActivity(activity.RegisterSms)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
