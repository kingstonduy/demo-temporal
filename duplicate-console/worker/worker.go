package main

import (
	"log"

	"kingstonduy/demo-temporal/async/shared"
	duplicate_console "kingstonduy/demo-temporal/duplicate-console"

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

	w1.RegisterWorkflow(duplicate_console.ParallelWorkFlow)
	w1.RegisterWorkflow(duplicate_console.AsyncWorkFlow)

	w1.RegisterActivity(duplicate_console.RegisterAccount)
	w1.RegisterActivity(duplicate_console.RegisterSms)
	w1.RegisterActivity(duplicate_console.NotificationSms)
	w1.RegisterActivity(duplicate_console.GetBalanceById)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
