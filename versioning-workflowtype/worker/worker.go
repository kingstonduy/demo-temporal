package main

import (
	"log"

	versioning_workflowtype "kingstonduy/demo-temporal/versioning-workflowtype"
	"kingstonduy/demo-temporal/versioning-workflowtype/shared"

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

	w.RegisterWorkflow(versioning_workflowtype.SimpleWorkflow)
	w.RegisterActivity(versioning_workflowtype.GetInformation)

	w.RegisterWorkflow(versioning_workflowtype.SimpleWorkflow1)
	w.RegisterActivity(versioning_workflowtype.GetInformation1)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
