package main

import (
	traditionalway_getting_block "kingstonduy/demo-temporal/traditionalway-getting-block"
	"kingstonduy/demo-temporal/traditionalway-getting-block/shared"
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

	w1.RegisterWorkflow(traditionalway_getting_block.BlockingWorkflow)

	w1.RegisterActivity(traditionalway_getting_block.Withdraw)
	w1.RegisterActivity(traditionalway_getting_block.UserInputOtp)
	w1.RegisterActivity(traditionalway_getting_block.Notification)
	w1.RegisterActivity(traditionalway_getting_block.ResendOtp)
	w1.RegisterActivity(traditionalway_getting_block.LongAcitivity)

	err = w1.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker1", err)
	}
}
