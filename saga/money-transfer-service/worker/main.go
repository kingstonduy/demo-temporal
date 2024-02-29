package main

import (
	"fmt"
	"log"
	"os"
	shared "saga"
	"saga/money-transfer-service/bootstrap"
	"saga/money-transfer-service/repository"
	"saga/money-transfer-service/temporal"
	"saga/money-transfer-service/usecase"
	"time"

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

	app := bootstrap.App()

	usecase := usecase.NewMoneyTransferUsecase(repository.NewTransactionRepository(app.Postgres),
		app.Env.NapasUrl, app.Env.T24Url, app.Env.LimitServiceUrl, time.Duration(app.Env.ContextTimeout))

	workflow := temporal.MoneyTransferWorkflow{
		Usecase: usecase,
	}

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, shared.TASKQUEUE, worker.Options{})

	w.RegisterWorkflow(workflow.NewMoneyTransferWorkflow)

	w.RegisterActivity(usecase.ValidateAccount)

	w.RegisterActivity(usecase.CompensateTransaction)
	w.RegisterActivity(usecase.UpdateStateCreated)

	w.RegisterActivity(usecase.LimitCut)
	w.RegisterActivity(usecase.LimitCutCompensate)

	w.RegisterActivity(usecase.UpdateStateLimitCut)

	w.RegisterActivity(usecase.MoneyCut)
	w.RegisterActivity(usecase.MoneyCutCompensate)

	w.RegisterActivity(usecase.UpdateStateMoneyCut)

	w.RegisterActivity(usecase.UpdateMoney)
	w.RegisterActivity(usecase.UpdateMoneyCompensate)

	w.RegisterActivity(usecase.UpdateStateTransactionCompleted)

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
