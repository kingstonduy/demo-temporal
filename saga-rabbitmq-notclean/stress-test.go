package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type SaferResponse struct {
	Code    int
	Message string
}

var ch chan int64

// func OneRequest() (avg float64) {

// }

func MultipleRequest(numberOfTimes int, concurrent int) {
	ch = make(chan int64)
	// Specify the body of the POST request
	body := map[string]interface{}{
		"FromAccountID": "OCB12345",
		"ToAccountID":   "TMCP23456",
		"Amount":        1,
	}
	bodyBytes, _ := json.Marshal(body)
	var res float64
	for i := 1; i <= numberOfTimes; i++ {
		// send post requests to localhost:7201/api/v1/moneytransfer and get response
		for j := 0; j < concurrent; j++ {
			go func() {
				now := time.Now().UnixMilli()
				response, err := http.Post("http://localhost:7201/api/v1/moneytransfer", "application/json", bytes.NewBuffer(bodyBytes))
				if err != nil {
					log.Printf(err.Error())
				} else {
					defer response.Body.Close()
					respBodyBytes, _ := io.ReadAll(response.Body)
					var resp SaferResponse
					json.Unmarshal(respBodyBytes, &resp)
					end := time.Now().UnixMilli()
					log.Printf("💡Request  takes %f", float64(end-now)/float64(1000))
					ch <- (end - now)
				}
			}()
		}
		var sum int64
		for i := 0; i < concurrent; i++ {
			sum += <-ch
		}

		var avg float64
		avg = float64(sum) / float64(concurrent)

		res += avg

		fmt.Printf("Batch %d, Average second per request: %fs\n", i, avg*0.001)
	}

	res = res / float64(numberOfTimes)

	log.Printf("Total batch: %d. Average second per batch: %fs", numberOfTimes, res*0.001)
}

func main() {
	MultipleRequest(5, 50)
}
