package controller

import (
	"kingstonduy/demo-temporal/saga/money-transfer-service/bootstrap"
	"kingstonduy/demo-temporal/saga/money-transfer-service/domain"

	"github.com/gin-gonic/gin"
)

type MoneyTransferController struct {
	MoneyTransferUsecase domain.MoneyTransferUsecase
	Env                  *bootstrap.Env
}

func (lc *MoneyTransferController) MoneyTransfer(c *gin.Context) {

}
