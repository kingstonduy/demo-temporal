package main

import (
	"log"
	"net/http"
	"saga-kafka-notclean/limitation-manage-service/service"
	"saga-kafka-notclean/money-transfer-service/config"
	model "saga-kafka-notclean/money-transfer-service/shared"
	"time"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch err.Error() {
	case "Not enough money":
		c.IndentedJSON(http.StatusBadRequest, model.SaferResponse{
			Code:    http.StatusBadRequest,
			Message: "Not enough money",
		})
		return
	case "Cannot find account":
		c.IndentedJSON(http.StatusNotFound, model.SaferResponse{
			Code:    http.StatusNotFound,
			Message: "Cannot find account",
		})
		return
	case "Cannot connect to database":
		c.IndentedJSON(http.StatusInternalServerError, model.SaferResponse{
			Code:    http.StatusInternalServerError,
			Message: "Cannot connect to database",
		})
		return
	default:
		c.IndentedJSON(http.StatusBadRequest, model.SaferResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid data format",
		})
		return
	}
}

func limit(c *gin.Context) {
	time.Sleep(time.Second * 5)
	var req model.SaferRequest
	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	log.Printf("💡Request %+v\n", req)

	err = service.LimitService(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, model.SaferResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

func main() {
	router := gin.Default()
	router.POST("/api/v1/account/limit", limit)
	router.Run(config.GetConfig().Limit.Host + ":" + config.GetConfig().Limit.Port)
}
