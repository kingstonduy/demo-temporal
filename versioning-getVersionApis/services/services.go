package services

import (
	"log"
	"os"
	"time"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

func GetInformation() string {
	time.Sleep(time.Second * 10)
	return "💡Default"
}

func GetInformation1() string {
	time.Sleep(time.Second * 10)
	return "💡Version 1"
}

func GetInformation2() string {
	time.Sleep(time.Second * 10)
	return "💡Version 2"
}
