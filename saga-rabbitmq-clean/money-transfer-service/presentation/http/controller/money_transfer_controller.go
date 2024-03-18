package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ocb/mcs-money-transfer/domain"
)

type MoneyTransferController struct {
}

func NewMoneyTransferController() *MoneyTransferController {
	return &MoneyTransferController{}
}

// Swagger information here
func (c *MoneyTransferController) TransferMoney(ctx *fiber.Ctx) error {
	return RequestHandler[*domain.MoneyTransferClientRequest, *domain.MoneyTransferClientResponse](ctx)
}
