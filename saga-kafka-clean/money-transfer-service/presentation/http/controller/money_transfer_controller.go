package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/http"
	"github.com/lengocson131002/go-clean/pkg/logger"
	"github.com/lengocson131002/go-clean/pkg/pipeline"
)

type MoneyTransferController struct {
	Logger logger.Logger
}

func NewMoneyTransferController(logger logger.Logger) *MoneyTransferController {
	return &MoneyTransferController{
		Logger: logger,
	}
}

func (c *MoneyTransferController) MoneyTransfer(ctx *fiber.Ctx) error {
	request := new(domain.MoneyTransferClientRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Logger.Warn("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := pipeline.Send[*domain.MoneyTransferClientRequest, *domain.MoneyTransferClientResponse](ctx.Context(), request)
	if err != nil {
		c.Logger.Warn("Failed to open account : %+v", err)
		return err
	}

	httpResp := http.SuccessResponse[*domain.MoneyTransferClientResponse](response)
	return ctx.Status(httpResp.Status).JSON(httpResp)
}
