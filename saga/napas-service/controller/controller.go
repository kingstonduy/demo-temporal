package main

import (
	shared "kingstonduy/demo-temporal/saga"
	"kingstonduy/demo-temporal/saga/napas-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func verify(c *gin.Context) {
	var req shared.ValidateAccountInput
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Invalid data format")
		return
	}

	res, err := service.VerifyAccount(req)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, res)
}

func main() {
	router := gin.Default()
	router.POST("/verify", verify)

	router.Run()
}
