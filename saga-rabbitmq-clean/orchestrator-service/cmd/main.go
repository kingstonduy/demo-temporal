package main

import (
	"orchestrator-service/bootstrap"
	"orchestrator-service/infra/rabbitmq"
	"orchestrator-service/infra/repository"
	"orchestrator-service/presentation"
	usecase "orchestrator-service/usecase/money_transfer"
)

func main() {
	cfg := bootstrap.NewConfig()

	log := bootstrap.GetLogger()

	moneyTransferRepository := repository.NewMoneytransferRepository(cfg)
	moneyTransferMessageQueue := rabbitmq.NewMoneyTransferRabbitMQ(cfg)

	activities := usecase.NewMoneyTransferActivities(cfg, log, moneyTransferRepository, moneyTransferMessageQueue)

	moneyTransferWorker := presentation.NewMoneyTransferWorker(activities, cfg)

	moneyTransferWorker.Run()
}
