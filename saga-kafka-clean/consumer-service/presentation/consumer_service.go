package presentation

import (
	"consumer-service/bootstrap"
	"consumer-service/domain"
	"consumer-service/pkg/logger"
	handlers "consumer-service/usecase"
)

func ConsumerService(log logger.Logger, cfg *bootstrap.Config) {
	var forever chan int

	moneyTransferHandler := handlers.NewMoneyTransferHandler(log, cfg)
	moneyTransferHandler.Handle(domain.SaferResponse{})

	<-forever
}
