package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	model "demo-temporal/model"

	"github.com/gin-gonic/gin"
)

var infoLogger *log.Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
var errorLogger *log.Logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

var cnt int = 4

var db = []model.Account{
	{Cif: "1", Balance: 1000, IsSms: true, IsEmail: true},
	{Cif: "2", Balance: 2000, IsSms: true, IsEmail: true},
	{Cif: "3", Balance: 3000, IsSms: true, IsEmail: true},
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
}

func getAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, db)
}

func getAccountById(c *gin.Context) {
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
	id := c.Param("id")

	for i, a := range db {
		if a.Cif == id {
			if a.IsSms == true {
				c.IndentedJSON(http.StatusConflict, gin.H{"message": "Sms already registered"})
				return
			} else {
				a.IsSms = true
				db[i] = a
				c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("register Sms for cif=%s successfully", id),
					"account": db[i]})
			}

			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "account not found"})
}

func registerEmail(c *gin.Context) {
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

func main() {
	router := gin.Default()
	router.POST("/account/register", register)
	router.GET("/account/:id", getAccountById)
	router.GET("/account/all", getAccounts)
	router.POST("/account/register/sms/:id", registerSms)
	router.POST("/account/register/email/:id", registerEmail)
	router.POST("/withdraw/:id/:amount", withdraw)
	router.POST("/deposit/:id/:amount", deposit)

	router.Run("localhost:8080")
}
