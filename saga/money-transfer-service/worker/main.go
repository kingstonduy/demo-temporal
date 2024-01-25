package main

import (
	"fmt"
	shared "kingstonduy/demo-temporal/saga"
	app "kingstonduy/demo-temporal/saga/money-transfer-service"
	"log"
	"os"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, shared.TASKQUEUE, worker.Options{})
	w.RegisterWorkflow(app.MoneyTransferWorkflow)

	w.RegisterActivity(app.ValidateAccount)

	w.RegisterActivity(app.CompensateTransaction)
	w.RegisterActivity(app.UpdateStateCreated)

	w.RegisterActivity(app.LimitCut)
	w.RegisterActivity(app.LimitCutCompensate)

	w.RegisterActivity(app.UpdateStateLimitCut)

	w.RegisterActivity(app.MoneyCut)
	w.RegisterActivity(app.MoneyCutCompensate)

	w.RegisterActivity(app.UpdateStateMoneyCut)

	w.RegisterActivity(app.UpdateMoney)
	w.RegisterActivity(app.UpdateMoneyCompensate)

	w.RegisterActivity(app.UpdateStateTransactionCompleted)

	// Start listening to the Task Queue
	go func() {
		var input string
		fmt.Scanln(&input)
		if input == "stop" {
			w.Stop()
		}
		os.Exit(0)
	}()

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

	// i want that when i type stop in the terminal, the worker will stop

}
