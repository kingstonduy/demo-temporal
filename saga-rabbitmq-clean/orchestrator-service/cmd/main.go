package main

import (
	"orchestrator-service/bootstrap"
	"orchestrator-service/infra/repository"
	"orchestrator-service/presentation"
	usecase "orchestrator-service/usecase/money_transfer"
)

func main() {
	log := bootstrap.GetLogger()

	temporalClient := bootstrap.GetTemporalClient()

	db, _ := bootstrap.GetDB()

	repository := repository.NewMoneytransferRepository(db)

	activities := usecase.NewMoneyTransferActivities(log, repository)

	presentation.MoneytransferWorker(activities, temporalClient)
}
