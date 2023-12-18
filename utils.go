package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func ReceiveFromApi(url string, tp string) (string, error) {
	if tp == "GET" {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		content := string(body)
		status := resp.StatusCode
		if status >= 400 {
			// API service dead
			message := fmt.Sprintf("HTTP Error %d: %s", status, content)
			return "", errors.New(message)
		}

		return content, nil
	} else {
		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		content := string(body)
		status := resp.StatusCode
		if status >= 400 {
			// API service dead
			message := fmt.Sprintf("HTTP Error %d: %s", status, content)
			return "", errors.New(message)
		}

		return content, nil
	}

}
