package engine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/chehsunliu/poker"
)

func CompareCardSlices(slice1, slice2 []poker.Card) bool {
	// Check if lengths are the same
	if len(slice1) != len(slice2) {
		return false
	}
	// Check each element for equality
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func getRandomTrueKey(m map[int]bool) (int, error) {
	var keys []int
	for k, v := range m {
		if v {
			keys = append(keys, k)
		}
	}

	// Return false if no valid keys are found
	if len(keys) == 0 {
		return -1, errors.New("all seats are full")
	}

	randomIndex := rand.Intn(len(keys))

	return keys[randomIndex], nil
}

func getUserToken(username string, password string) (string, error) {
	url := fmt.Sprintf("https://%s/oauth/token", os.Getenv("AUTH0_DOMAIN"))

	// Prepare the request body
	bodyData := map[string]string{
		"grant_type": "password",
		"client_id":  os.Getenv("AUTH0_CLIENT_ID"),
		"audience":   os.Getenv("AUTH0_AUDIENCE"),
		"username":   username,
		"password":   password,
	}

	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body using io.ReadAll (instead of ioutil.ReadAll)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract the access token
	token, ok := jsonResponse["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return token, nil
}