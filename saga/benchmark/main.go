package main

import (
	"fmt"
	shared "kingstonduy/demo-temporal/saga"
	"sync"

	"github.com/pborman/uuid"
)

func sth() {
	fmt.Println("STH")
}

func main() {
	var input = shared.TransactionInfo{
		FromAccountId: "123456789",
		ToAccountId:   "987654321",
		Amount:        1000000,
	}

	input.TransactionId = uuid.New()

	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			url := fmt.Sprintf("http://localhost:7201/api/v1/moneytransfer")
			var responseType string
			err := shared.PostApi(url, &input, &responseType)
			if err != nil {
				fmt.Println("Cannot transfer money")
			} else {
				fmt.Println("Transfer money successfully")
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	// code sth
}
