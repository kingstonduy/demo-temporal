package main

import (
	"context"
	"fmt"
	"log"

	"signal"
	"signal/shared"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type Person struct {
	Name string
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	optionsAsync := client.StartWorkflowOptions{
		ID:        "temporal-demo-workflow-signal_" + uuid.New(),
		TaskQueue: shared.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), optionsAsync, signal.SignalWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	if err != nil {
		log.Fatalln("Unable to execute async workflow", err)
	}

	fmt.Println("User input something")
	var msg string
	fmt.Scanln(&msg)

	err = c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), "MySignal", msg)
	fmt.Println("💡Sending signal = " + msg)

}
