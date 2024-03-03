package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

type SaferResponse struct {
	Code    int
	Message string
}

func main() {
	concurrent := 100
	var wg sync.WaitGroup

	// Specify the body of the POST request
	body := map[string]interface{}{
		"FromAccountID": "OCB12345",
		"ToAccountID":   "TMCP23456",
		"Amount":        1,
	}
	bodyBytes, _ := json.Marshal(body)

	// send post requests to localhost:7201/api/v1/moneytransfer and get response
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			response, err := http.Post("http://localhost:7201/api/v1/moneytransfer", "application/json", bytes.NewBuffer(bodyBytes))
			if err != nil {
				log.Printf(err.Error())
			} else {

				defer response.Body.Close()
				respBodyBytes, _ := io.ReadAll(response.Body)
				var resp SaferResponse
				json.Unmarshal(respBodyBytes, &resp)
				log.Printf("💡%s. %+v", response.Status, resp)

			}
		}()
	}

	wg.Wait()
}
