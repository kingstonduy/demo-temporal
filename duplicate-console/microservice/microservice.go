package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"kingstonduy/demo-temporal/duplicate-console/model"

	"github.com/gin-gonic/gin"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

var cnt int = 4

var db = []model.Account{
	{Cif: "1", Balance: 1000, IsSms: false, IsEmail: false},
	{Cif: "2", Balance: 2000, IsSms: false, IsEmail: false},
	{Cif: "3", Balance: 3000, IsSms: false, IsEmail: false},
}

func register(c *gin.Context) {

	newAccount := model.Account{
		Cif:     strconv.Itoa(cnt),
		Balance: 0,
		IsSms:   false,
		IsEmail: false,
	}
	cnt++

	db = append(db, newAccount)
	c.IndentedJSON(http.StatusCreated, newAccount)
	time.Sleep(time.Second * 5)
}

func getAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, db)
	time.Sleep(time.Second * 5)
}

func getAccountById(c *gin.Context) {
	time.Sleep(time.Second * 5)
	id := c.Param("id")

	for _, a := range db {
		if a.Cif == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})

}

func registerSms(c *gin.Context) {
	time.Sleep(time.Second * 5)
	id := c.Param("id")

	for i, a := range db {
		if a.Cif == id {
			a.IsSms = true
			db[i] = a
			c.IndentedJSON(http.StatusOK, db[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func registerEmail(c *gin.Context) {
	time.Sleep(time.Second * 5)
	id := c.Param("id")

	for i, a := range db {
		if a.Cif == id {
			if a.IsEmail == true {
				c.IndentedJSON(http.StatusConflict, gin.H{"message": "Email already registered"})
				return
			} else {
				a.IsEmail = true
				db[i] = a
				c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("register email for cif=%s successfully", id)})
			}

			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func withdraw(c *gin.Context) {
	time.Sleep(time.Second * 5)
	id := c.Param("id")
	amount, err := strconv.ParseFloat(c.Param("amount"), 64)
	if err != nil {
		errorLogger.Println("Invalid amount")
		return
	}

	for i, a := range db {
		if a.Cif == id {
			if a.Balance < amount {
				c.IndentedJSON(http.StatusConflict, gin.H{"message": "Insufficient balance"})
				return
			} else {
				a.Balance -= amount
				db[i] = a
				c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("withdraw money for cif=%s successfully", id)})
			}

			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func deposit(c *gin.Context) {
	time.Sleep(time.Second * 5)
	id := c.Param("id")
	amount, err := strconv.ParseFloat(c.Param("amount"), 64)
	if err != nil {
		errorLogger.Println("Invalid amount")
		return
	}

	for i, a := range db {
		if a.Cif == id {
			a.Balance += amount
			db[i] = a
			c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Deposit money for cif=%s successfully", id)})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})

}

func getBalanceById(c *gin.Context) {
	time.Sleep(time.Second * 10)
	id := c.Param("id")

	for _, a := range db {
		if a.Cif == id {
			c.IndentedJSON(http.StatusOK, a.Balance)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func main() {
	router := gin.Default()
	router.POST("/account/register", register)
	router.GET("/account/:id", getAccountById)
	router.GET("/account/all", getAccounts)
	router.GET("/account/balance/:id", getBalanceById)
	router.POST("/account/register/sms/:id", registerSms)
	router.POST("/account/register/email/:id", registerEmail)
	router.POST("/withdraw/:id/:amount", withdraw)
	router.POST("/deposit/:id/:amount", deposit)

	router.Run("localhost:8080")
}
