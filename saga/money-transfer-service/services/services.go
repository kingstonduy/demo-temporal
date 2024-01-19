package main

import (
	"kingstonduy/demo-temporal/async/model"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

var otp model.Otp

func main() {
	router := gin.Default()

	router.Run("localhost:8080")
}
