package common

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RequestGET(baseURL string, rawParams map[string]string) ([]byte, error) {
	/*
	   Sends a GET request to the provided url and retrieves the response
	*/
	client := &http.Client{}
	params := url.Values{}
	for key := range rawParams {
		params.Add(key, rawParams[key])
	}
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}

	// Set the "Accept" header to request JSON response
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error in GET request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Request failed with status code:", resp.StatusCode)
		return nil, nil
	}

	return io.ReadAll(resp.Body)
}

func RequestPOST(url string, payload []byte) ([]byte, error) {
	/*
		Sends a POST request to the provided url and retrieves the response
	*/
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Error creating HTTP request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPEN_AI_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending HTTP request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response body", err)
			return nil, err
		}
		return body, nil
	}
	log.Println("Payload used: ", string(payload))
	log.Fatal("Error with status code: ", resp.StatusCode)
	return nil, nil
}

func CleanForJSON(text string) string {
	// Fix json respresentation of soon-to-be payload
	text = strings.ReplaceAll(text, "\n", "\\n")
	text = strings.ReplaceAll(text, "\"", "\\\"")
	return text
}

func ReadInput(reader *bufio.Reader) string {
	// Read a line of text including spaces
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return text
}

func PruneText(text string, limit int) string {
	if len(text) > limit {
		return text[len(text)-limit:]
	}
	return text
}
