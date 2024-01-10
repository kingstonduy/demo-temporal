package main

import (
	"crypto/rand"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"kingstonduy/demo-temporal/goroutines-activities/model"

	"github.com/gin-gonic/gin"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

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

func otpVerify(c *gin.Context) {
	var req model.RequestOtp
	if err := c.BindJSON(&req); err != nil {
		errorLogger.Println("Invalid data format", err)
		c.IndentedJSON(http.StatusBadRequest, "Invalid data format")
		return
	}

	res := model.ResponseOtp{
		Check: true,
	}

	c.IndentedJSON(http.StatusOK, res)
	return

	// if validateOtpService(req) {
	// 	infoLogger.Println("Otp is valid")
	// 	res.Check = true
	// 	c.IndentedJSON(http.StatusOK, res)
	// 	return
	// } else {
	// 	errorLogger.Println("Otp is invalid")
	// 	c.IndentedJSON(http.StatusUnauthorized, res)
	// 	return
	// }
}

func getNotification(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.IndentedJSON(http.StatusOK, "You have successfully withdraw money")
}

func getOcbInfo(c *gin.Context) {
	time.Sleep(time.Second * 10)
	c.IndentedJSON(http.StatusOK, "The Orient Commercial Joint Stock Bank is a large bank located in Vietnam. Its Swift code is ORCOVNVX. As of April 2007, 10%"+"of Oricombank was owned by BNP Paribas, a French banking group.")
}

func main() {
	router := gin.Default()
	router.POST("/withdraw", withdraw)
	router.POST("/otp/verify", otpVerify)
	router.GET("/otp/resend", resendOtp)
	router.GET("/notification", getNotification)
	router.GET("OCB/info", getOcbInfo)
	router.Run("localhost:8080")
}

func genOtpService() model.Otp {
	limit := 6
	b := make([]byte, limit)
	n, err := io.ReadAtLeast(rand.Reader, b, limit)
	if n != limit {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	return model.Otp{
		Otp:       string(b),
		Timestamp: time.Now().Unix(),
	}
}

func validateOtpService(newOtp model.RequestOtp) bool {
	if newOtp.Otp == otp.Otp && time.Now().Unix()-otp.Timestamp < 60 {
		return true
	}
	return false
}
