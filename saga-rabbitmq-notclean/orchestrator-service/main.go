package main

import (
	"fmt"
	"log"
	"os"
	"saga-rabbitmq-notclean/money-transfer-service/config"
	"saga-rabbitmq-notclean/orchestrator-service/workflow"

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
	w := worker.New(c, config.GetConfig().Temporal.TaskQueue, worker.Options{})

	w.RegisterWorkflow(workflow.MoneyTransferWorklow)

	w.RegisterActivity(workflow.ValidateAccount)

	w.RegisterActivity(workflow.CompensateTransaction)
	w.RegisterActivity(workflow.UpdateStateCreated)

	w.RegisterActivity(workflow.LimitCut)
	w.RegisterActivity(workflow.LimitCutCompensate)

	w.RegisterActivity(workflow.UpdateStateLimitCut)

	w.RegisterActivity(workflow.MoneyCut)
	w.RegisterActivity(workflow.MoneyCutCompensate)

	w.RegisterActivity(workflow.UpdateStateMoneyCut)

	w.RegisterActivity(workflow.UpdateMoney)
	w.RegisterActivity(workflow.UpdateMoneyCompensate)

	w.RegisterActivity(workflow.UpdateStateTransactionCompleted)

	// Start listening to the Task Queue
	go func() {
		var input string
		fmt.Scanln(&input)
		if input == "stop" {
			w.Stop()
			os.Exit(0)
		}
	}()

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}

	// i want that when i type stop in the terminal, the worker will stop

}
