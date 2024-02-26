package main

import (
	"log"

	versioning_getVersionApis "versioning-getVersionApis"
	"versioning-getVersionApis/shared"

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

	w1.RegisterWorkflow(versioning_getVersionApis.SimpleWorkflow)

	w1.RegisterActivity(versioning_getVersionApis.GetInformation)
	w1.RegisterActivity(versioning_getVersionApis.GetInformation1)
	w1.RegisterActivity(versioning_getVersionApis.GetInformation2)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
