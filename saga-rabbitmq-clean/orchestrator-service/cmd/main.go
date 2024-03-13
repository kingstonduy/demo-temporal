package main

import (
	"orchestrator-service/bootstrap"
	"orchestrator-service/infra/message_queue/rabbitmq"
	"orchestrator-service/infra/repository"
	usecase "orchestrator-service/usecase/money_transfer"
	"orchestrator-service/worker"
)

func main() {
	cfg := bootstrap.NewConfig()

	log := bootstrap.GetLogger()

	moneyTransferRepository := repository.NewMoneytransferRepository(cfg)
	moneyTransferMessageQueue := rabbitmq.NewMoneyTransferRabbitMQ(cfg)

	activities := usecase.NewMoneyTransferActivities(cfg, log, moneyTransferRepository, moneyTransferMessageQueue)

	moneyTransferWorker := worker.NewMoneyTransferWorker(activities, cfg)

	moneyTransferWorker.Run()
}
