package main

import (
	"log"

	"signal"
	"signal/shared"

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

	w1.RegisterWorkflow(signal.SignalWorkflow)

	w1.RegisterActivity(signal.BlockingActivity)
	w1.RegisterActivity(signal.InputActivity)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
