package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetApi[T any](url string, responseType *T) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = convertResponseBodyIntoObject(resp.Body, responseType)
	if err != nil {
		return err
	}

	status := resp.StatusCode
	if status >= 400 {
		// API service dead
		message := fmt.Sprintf("HTTP Error %d: %+v", status, responseType)
		return errors.New(message)
	}

	return nil
}

func PostApi[T any, K any](url string, requestType *T, responseType *K) error {
	reqReader, err := convertObjectIntoRequestBody(requestType)
	if err != nil {
		return fmt.Errorf("%s%s", NONRETRYABLE_ERROR, err.Error())
	}

	resp, err := http.Post(url, "application/json", reqReader)
	if err != nil {
		return fmt.Errorf("%s%s", RETRYABLE_ERROR, err.Error())
	}
	defer resp.Body.Close()

	convertResponseBodyIntoObject(resp.Body, responseType)

	status := resp.StatusCode
	if status != 200 {
		// API service dead
		message := fmt.Sprintf("HTTP Error %d: %+v", status, responseType)
		return fmt.Errorf("%s%s", NONRETRYABLE_ERROR, message)
	}
	return nil
}

func convertObjectIntoRequestBody[T any](obj T) (io.Reader, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(obj)
	if err != nil {
		log.Fatal(err)
	}
	return &buf, nil
}

func convertResponseBodyIntoObject[T any](body io.ReadCloser, responseType *T) error {
	responseBody, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, &responseType)
	if err != nil {
		return err
	}
	return nil
}
