package presentation

import (
	"log"
	"orchestrator-service/bootstrap"
	"orchestrator-service/domain"
	usecase "orchestrator-service/usecase/money_transfer"
	"sync"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func MoneytransferWorker(activities domain.MoneyTransferActivities, c *client.Client) {
	// Create the client object just once per process

	workers := 1

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			// This worker hosts both Workflow and Activity functions
			w := worker.New((*c), bootstrap.GetConfig().Temporal.TaskQueue, worker.Options{
				// MaxConcurrentWorkflowTaskPollers: 1,
				// MaxConcurrentActivityTaskPollers: 1,
				Identity: uuid.New(),
			})

			w.RegisterWorkflowWithOptions(usecase.MoneyTransferWorkflow, workflow.RegisterOptions{
				Name: "MoneyTransferService",
			})

			w.RegisterActivity(activities.ValidateAccount)

			w.RegisterActivity(activities.CompensateTransaction)
			w.RegisterActivity(activities.UpdateStateCreated)

			w.RegisterActivity(activities.LimitCut)
			w.RegisterActivity(activities.LimitCutCompensate)

			w.RegisterActivity(activities.UpdateStateLimitCut)

			w.RegisterActivity(activities.MoneyCut)
			w.RegisterActivity(activities.MoneyCutCompensate)

			w.RegisterActivity(activities.UpdateStateMoneyCut)

			w.RegisterActivity(activities.UpdateMoney)
			w.RegisterActivity(activities.UpdateMoneyCompensate)

			w.RegisterActivity(activities.UpdateStateTransactionCompleted)

			// Start listening to the Task Queue
			err := w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("unable to start Worker", err)
			}
		}()
	}
	wg.Wait()
}
