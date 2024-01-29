package main

import (
	shared "kingstonduy/demo-temporal/saga"
	"kingstonduy/demo-temporal/saga/t24-service/service"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	switch err.Error() {
	case "Not enough money":
		c.IndentedJSON(http.StatusBadRequest, shared.SaferResponse{
			Code:    http.StatusBadRequest,
			Message: "Not enough money",
		})
		return
	case "Cannot find account":
		log.Println("🔥🔥🔥🔥🔥🔥")
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

func amountCut(c *gin.Context) {
	var req shared.SaferRequest
	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	log.Printf("💡Request %+v\n", req)

	// note the difference
	req.Amount = -req.Amount

	err = service.AmountService(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	log.Println("💡OK")
	time.Sleep(shared.TIMEOUT)
	c.IndentedJSON(http.StatusOK, shared.SaferResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

func amountAdd(c *gin.Context) {
	time.Sleep(shared.TIMEOUT)
	var req shared.SaferRequest
	err := c.BindJSON(&req)
	if err != nil {
		HandleError(c, err)
		return
	}
	log.Printf("💡Request %+v\n", req)

	err = service.AmountService(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, shared.SaferResponse{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

func main() {
	router := gin.Default()
	router.POST("/api/v1/amount/cut", amountCut)
	router.POST("/api/v1/amount/add", amountAdd)
	router.Run(shared.T24_SERVICE_HOST_PORT)
}
