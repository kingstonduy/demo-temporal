package cmd

import (
	"kingstonduy/demo-temporal/saga/money-transfer-service/api/route"
	"kingstonduy/demo-temporal/saga/money-transfer-service/bootstrap"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Postgres

	defer app.CloseConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	route.Setup(env, timeout, db, gin, c)

	gin.Run(env.ServerHost)
}
