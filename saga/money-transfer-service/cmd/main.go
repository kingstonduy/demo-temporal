package main

import (
	"fmt"
	"log"
	"saga/money-transfer-service/api/route"
	"saga/money-transfer-service/bootstrap"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	fmt.Printf("%+v\n", env)

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

	go gin.Run(env.ServerHost)
	fmt.Println("HIHI")
}
