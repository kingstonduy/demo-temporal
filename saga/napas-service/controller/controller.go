package main

import (
	"log"
	"net/http"
	shared "saga"
	"saga/napas-service/service"
	"time"

	"github.com/gin-gonic/gin"
)

func verify(c *gin.Context) {
	time.Sleep(shared.SERVICE_TIMEOUT)
	var req shared.ValidateAccountInput
	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	log.Printf("💡Request %+v\n", req)

	napasEntity, err := service.VerifyAccount(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, napasEntity)
}

func update(c *gin.Context) {
	time.Sleep(shared.SERVICE_TIMEOUT)
	var req shared.SaferRequest
	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	log.Printf("💡Request %+v\n", req)

	err = service.UpdateMoney(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, shared.SaferResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

func HandleError(c *gin.Context, err error) {
	switch err.Error() {
	case "Not enough money":
		c.IndentedJSON(http.StatusBadRequest, shared.SaferResponse{
			Code:    http.StatusBadRequest,
			Message: "Not enough money",
		})
		return
	case "Cannot find account":
		c.IndentedJSON(http.StatusNotFound, shared.SaferResponse{
			Code:    http.StatusNotFound,
			Message: "Cannot find account",
		})
		return
	case "Cannot connect to database":
		c.IndentedJSON(http.StatusInternalServerError, shared.SaferResponse{
			Code:    http.StatusInternalServerError,
			Message: "Cannot connect to database",
		})
		return
	default:
		c.IndentedJSON(http.StatusBadRequest, shared.SaferResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid data format",
		})
		return
	}
}

func main() {
	router := gin.Default()
	router.POST("/api/v1/account/verify", verify)
	router.POST("/api/v1/account/update", update)
	router.Run(shared.NAPAS_SERVICE_HOST_PORT)
}
