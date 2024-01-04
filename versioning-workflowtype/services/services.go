package services

import (
	"log"
	"os"
	"time"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

var timeout time.Duration = 20

func GetInformation() string {
	time.Sleep(time.Second * timeout)
	return "💡Default"
}

func GetInformation1() string {
	time.Sleep(time.Second * timeout)
	return "💡Version 1"
}

func GetInformation2() string {
	time.Sleep(time.Second * timeout)
	return "💡Version 2"
}
