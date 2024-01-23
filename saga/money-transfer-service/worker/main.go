package main

import (
	shared "kingstonduy/demo-temporal/saga"
	app "kingstonduy/demo-temporal/saga/money-transfer-service"
	"log"

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

	// code a automatically register all activities. user reflection
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
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
