package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Set up a test server URL
	serverURL := "http://localhost:8000" // Update with your server URL and port

	// Test cases for each endpoint
	testCases := []struct {
		name     string
		method   string
		endpoint string
		body     interface{}
		expected int
	}{
		// Example test cases, update with your actual endpoints
		{"CheckServerStatus", http.MethodGet, "/", nil, http.StatusOK},

		{"GetOutTransaction", http.MethodGet, "/transactions/0x8C8D7C46219D9205f056f28fee5950aD564d7461", nil, http.StatusOK},

		{"GetCurrentBlockNumber", http.MethodGet, "/blockNumber", nil, http.StatusOK},

		{"AddSubscriber", http.MethodPost, "/subscriptions", map[string]string{"address": "0x1234567890123456789012345678901234567890"}, http.StatusOK},

		{"UnsubscribeSubscriber", http.MethodPut, "/subscriptions", map[string]string{"address": "0x1234567890123456789012345678901234567890"}, http.StatusOK},
	}

	// Slice to hold test results
	var results []string

	for _, tc := range testCases {
		fmt.Printf("Running test: %s %s\n", tc.method, tc.endpoint)

		// Prepare request body if provided
		var reqBody []byte
		if tc.body != nil {
			reqBody, _ = json.Marshal(tc.body)
		}

		// Create HTTP request
		start := time.Now()
		req, err := http.NewRequest(tc.method, serverURL+tc.endpoint, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Printf("Failed to create request: %v\n", err)
			continue
		}

		// Perform the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Failed to perform request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != tc.expected {
			results = append(results, fmt.Sprintf("%s | FAIL | %d | %s", tc.name, resp.StatusCode, time.Since(start)))
		} else {
			results = append(results, fmt.Sprintf("%s | PASS | %d | %s", tc.name, resp.StatusCode, time.Since(start)))
		}
	}

	// Print test results
	fmt.Println("\nTest Results:")
	for _, result := range results {
		fmt.Println(result)
	}
}
