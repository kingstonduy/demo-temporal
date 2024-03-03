package main

import (
	"log"
	shared "saga-kafka-notclean/config"
	"saga-kafka-notclean/orchestrator-service/workflow"
	"sync"

	"github.com/pborman/uuid"
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

	workers := 1

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			// This worker hosts both Workflow and Activity functions
			w := worker.New(c, shared.GetConfig().Temporal.TaskQueue, worker.Options{
				// MaxConcurrentWorkflowTaskPollers: 1,
				// MaxConcurrentActivityTaskPollers: 1,
				Identity: uuid.New(),
			})

			w.RegisterWorkflow(workflow.MoneyTransferWorkflow)

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
			err = w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("unable to start Worker", err)
			}
		}()
	}
	wg.Wait()

}
