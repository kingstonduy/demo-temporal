package worker

import (
	"log"
	"orchestrator-service/bootstrap"
	"orchestrator-service/domain"
	usecase "orchestrator-service/usecase/money_transfer"
	"sync"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type moneyTransferWorker struct {
	activities domain.MoneyTransferActivities
	cfg        *bootstrap.Config
}

func NewMoneyTransferWorker(activities domain.MoneyTransferActivities, cfg *bootstrap.Config) moneyTransferWorker {
	return moneyTransferWorker{
		activities: activities,
		cfg:        cfg,
	}
}

func (mw *moneyTransferWorker) Run() {
	// Create the client object just once per process

	workers := 10

	wg := &sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			// This worker hosts both Workflow and Activity functions
			w := worker.New(*bootstrap.GetTemporalClient(mw.cfg), mw.cfg.Temporal.TaskQueue, worker.Options{
				Identity: uuid.New(),
			})

			w.RegisterWorkflowWithOptions(usecase.MoneyTransferWorkflow, workflow.RegisterOptions{
				Name: "MoneyTransferService",
			})

			w.RegisterActivity(mw.activities.ValidateAccount)
			w.RegisterActivity(mw.activities.CompensateTransaction)
			w.RegisterActivity(mw.activities.UpdateStateCreated)
			w.RegisterActivity(mw.activities.LimitCut)
			w.RegisterActivity(mw.activities.LimitCutCompensate)
			w.RegisterActivity(mw.activities.UpdateStateLimitCut)
			w.RegisterActivity(mw.activities.MoneyCut)
			w.RegisterActivity(mw.activities.MoneyCutCompensate)
			w.RegisterActivity(mw.activities.UpdateStateMoneyCut)
			w.RegisterActivity(mw.activities.UpdateMoney)
			w.RegisterActivity(mw.activities.UpdateMoneyCompensate)
			w.RegisterActivity(mw.activities.UpdateStateTransactionCompleted)

			// Start listening to the Task Queue
			err := w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("unable to start Worker", err)
			}
		}()
	}
	wg.Wait()
}
