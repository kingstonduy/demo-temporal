package cmd

import (
	"kingstonduy/demo-temporal/saga/money-transfer-service/api/route"
	"kingstonduy/demo-temporal/saga/money-transfer-service/bootstrap"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Postgres

	defer app.CloseConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, db, gin)

	gin.Run(env.ServerHost)
}
