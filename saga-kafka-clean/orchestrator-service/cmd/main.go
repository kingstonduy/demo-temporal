package main

import (
	"orchestrator-service/bootstrap"
	infra "orchestrator-service/infra/kafka"
	"orchestrator-service/infra/repository"
	"orchestrator-service/presentation"
	usecase "orchestrator-service/usecase/money_transfer"
)

func main() {
	log := bootstrap.GetLogger()

	config := bootstrap.NewConfig()

	temporalClient := config.GetTemporalClient()

	db, _ := config.GetDB()

	repository := repository.NewMoneytransferRepository(db)

	kafkaClient := infra.NewKafkaClient(config)

	activities := usecase.NewMoneyTransferActivities(log, repository, kafkaClient, config)

	presentation.MoneytransferWorker(activities, temporalClient)
}
