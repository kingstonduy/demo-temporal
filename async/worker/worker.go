package main

import (
	"async"
	"async/shared"
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

	w1.RegisterWorkflow(async.AsyncWorkFlow)

	w1.RegisterActivity(async.Withdraw)
	w1.RegisterActivity(async.UserInputOtp)
	w1.RegisterActivity(async.Notification)
	w1.RegisterActivity(async.ResendOtp)
	w1.RegisterActivity(async.GetOcbInfo)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
