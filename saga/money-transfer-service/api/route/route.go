package route

import (
	"kingstonduy/demo-temporal/saga/money-transfer-service/bootstrap"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, gin *gin.Engine) {
	publicRouter := gin.Group("/api/v1")

	NewMoneyTransferRouter(env, timeout, db, publicRouter)
}
