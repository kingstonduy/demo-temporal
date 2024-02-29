package route

import (
	"saga/money-transfer-service/api/controller"
	"saga/money-transfer-service/bootstrap"
	"saga/money-transfer-service/repository"
	"saga/money-transfer-service/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func NewMoneyTransferRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup, c client.Client) {
	ur := repository.NewTransactionRepository(db)
	sc := controller.MoneyTransferController{
		MoneyTransferUsecase: usecase.NewMoneyTransferUsecase(ur, env.NapasUrl, env.T24Url, env.LimitServiceUrl, timeout),
		Env:                  env,
	}
	group.POST("/moneytransfer", sc.MoneyTransfer(c))
}
