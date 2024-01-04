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

	w1 := worker.New(c, shared.TaskQueueName, worker.Options{})

	w1.RegisterWorkflow(workflow.SimpleWorkflow)

	w1.RegisterActivity(activity.GetInformation)
	w1.RegisterActivity(activity.GetInformation1)
	w1.RegisterActivity(activity.GetInformation2)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
