package main

import (
	"fmt"
	shared "saga"
	"sync"

	"github.com/pborman/uuid"
)

func sth() {
	fmt.Println("STH")
}

func main() {
	var input = shared.TransactionInfo{
		FromAccountId: "OCB12345",
		ToAccountId:   "TMCP23456",
		Amount:        0,
	}

	input.TransactionId = uuid.New()

	wg := sync.WaitGroup{}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			url := fmt.Sprintf("http://localhost:7201/api/v1/moneytransfer")
			var responseType shared.SaferResponse
			_ = shared.PostApi(url, &input, &responseType)
			fmt.Printf("responseType: %+v\n", responseType)
			defer wg.Done()
		}()
	}
	wg.Wait()
	// code sth
}
