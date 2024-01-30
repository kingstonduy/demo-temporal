package route

import (
	"kingstonduy/demo-temporal/saga/money-transfer-service/api/controller"
	"kingstonduy/demo-temporal/saga/money-transfer-service/bootstrap"
	"kingstonduy/demo-temporal/saga/money-transfer-service/repository"
	"kingstonduy/demo-temporal/saga/money-transfer-service/usecase"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewMoneyTransferRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	ur := repository.NewTransactionRepository(db)
	sc := controller.MoneyTransferController{
		MoneyTransferUsecase: usecase.NewMoneyTransferUsecase(ur, timeout),
		Env:                  env,
	}
	group.POST("/signup", sc.MoneyTransfer)
}
