package main

import (
	"kingstonduy/demo-temporal/async/model"
	app "kingstonduy/demo-temporal/saga/money-transfer-service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

var otp model.Otp

func withdraw(c *gin.Context) {
	time.Sleep(5 * time.Second)
	otp = genOtpService()
	infoLogger.Println("💡Sending Otp to user's email.The otp =", otp)
	c.IndentedJSON(http.StatusOK, "Otp sent")
}

func resendOtp(c *gin.Context) {
	otp = genOtpService()
	infoLogger.Println("💡Sending Otp to user's email.The otp =", otp)
	c.IndentedJSON(http.StatusOK, "Otp sent")
}

func verify(c *gin.Context) {
	var req app.ValidateAccountInput
	if err := c.BindJSON(&req); err != nil {
		errorLogger.Println("Invalid data format", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid data format")
		return
	}

	isValid, err := app.DB.SearchNapasByAccountId(req.AccountId)
	if err != nil {
		errorLogger.Println("🔥Error when searching Napas by account id", err)
		c.IndentedJSON(http.StatusInternalServerError, "🔥Error when searching Napas by account id")
	}

	res := app.ValidateAccountOutput{}
}

func main() {
	router := gin.Default()
	router.POST("/withdraw", withdraw)
	router.POST("/otp/verify", otpVerify)
	router.GET("/otp/resend", resendOtp)
	router.GET("/notification", getNotification)
	router.GET("OCB/info", getOcbInfo)
	router.POST("/verify", verify)
	router.Run("localhost:8080")
}
