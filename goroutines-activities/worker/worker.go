package main

import (
	"log"

	goroutines_activities "goroutines-activities"
	"goroutines-activities/shared"

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

	w1.RegisterWorkflow(goroutines_activities.AsyncWorkFlow1)

	w1.RegisterActivity(goroutines_activities.Withdraw)
	w1.RegisterActivity(goroutines_activities.UserInputOtp)
	w1.RegisterActivity(goroutines_activities.Notification)
	w1.RegisterActivity(goroutines_activities.ResendOtp)
	w1.RegisterActivity(goroutines_activities.LongAcitivity)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
