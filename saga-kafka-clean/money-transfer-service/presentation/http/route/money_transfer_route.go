package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lengocson131002/go-clean/presentation/http/controller"
)

func RegisterMoneyTransferRoute(root *fiber.Router, t24Con *controller.MoneyTransferController) {

	(*root).Post("/moneyTransfer", t24Con.MoneyTransfer)
}
