package route

import (
	"saga/money-transfer-service/bootstrap"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, gin *gin.Engine, c client.Client) {
	publicRouter := gin.Group("/api/v1")

	NewMoneyTransferRouter(env, timeout, db, publicRouter, c)
}
