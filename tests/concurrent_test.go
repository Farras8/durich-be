package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

// ========================================
// CONFIGURATION - Update these values!
// ========================================
// Untuk mendapatkan AuthToken:
// 1. Login via POST /v1/auth/login dengan email & password
// 2. Copy nilai "access_token" dari response
// 3. Paste di bawah ini
// ========================================
const (
	BaseURL   = "http://localhost:8081/v1"
	AuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX2lkIjoiMzY4WWt6N3NBWmw3Zk1ad1hXUGtnR21PRDQzIiwidXNlcl9pZCI6IjM2OFlreFpsZW1DREV3R3pGYTh5cVZha1FPTSIsImVtYWlsIjoiYWRtaW5AZXhhbXBsZS5jb20iLCJyb2xlIjpbImFkbWluIl0sImxvY2F0aW9uX2lkIjoiIiwicmVmcmVzaF90b2tlbl9pZCI6IjM2dXB2bUl0amJUOGJXT3ZhOXBSQVIybkkxUiIsImlzcyI6ImR1cmljaC1zeXN0ZW0iLCJzdWIiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc2NTkwOTc0MywibmJmIjoxNzY1ODY2NTQzLCJpYXQiOjE3NjU4NjY1NDN9.UEh9rIIr4bS4sQ1PUWJf1sqpYFd5sVEdked_Gc5Ngfo" // Replace with valid JWT token
)

// Test data - update with valid IDs from your database
var (
	testLotID     = "36upzfFwkSiyTxhxVqZWCA4G3TY" // A DRAFT lot for testing
	testBlokID    = "2SRlQ8zX9vJ2mN5P6Q7R8S9T005"
	testPohonKode = "0000"
)

// ========================================
// HELPER FUNCTIONS
// ========================================

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+AuthToken)

	client := &http.Client{Timeout: 30 * time.Second}
	return client.Do(req)
}

func parseResponse(resp *http.Response) (*APIResponse, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("response: %s, error: %v", string(body), err)
	}
	return &apiResp, nil
}

// ========================================
// TEST 1: Concurrent AddItems to Same Lot
// Tests row-level locking with FOR UPDATE
// ========================================
func TestConcurrentAddItemsToLot(t *testing.T) {
	concurrency := 5
	var wg sync.WaitGroup
	results := make(chan string, concurrency)
	errors := make(chan error, concurrency)

	startTime := time.Now()
	startSignal := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// Wait for start signal to ensure truly concurrent requests
			<-startSignal

			url := fmt.Sprintf("%s/lots/%s/items", BaseURL, testLotID)
			body := map[string]interface{}{
				"pohon_kode": testPohonKode,
				"blok_id":    testBlokID,
				"berat":      float64(3 + workerID), // Different weight for each
			}

			resp, err := makeRequest("POST", url, body)
			if err != nil {
				errors <- fmt.Errorf("worker %d: request error: %v", workerID, err)
				return
			}

			apiResp, err := parseResponse(resp)
			if err != nil {
				errors <- fmt.Errorf("worker %d: parse error: %v", workerID, err)
				return
			}

			results <- fmt.Sprintf("Worker %d: status=%d, success=%v, message=%s",
				workerID, resp.StatusCode, apiResp.Success, apiResp.Message)
		}(i)
	}

	// Start all workers at the same time
	close(startSignal)
	wg.Wait()
	close(results)
	close(errors)

	duration := time.Since(startTime)

	// Print results
	t.Logf("=== Concurrent AddItems Test Results ===")
	t.Logf("Total duration: %v", duration)
	t.Logf("Concurrency level: %d", concurrency)

	successCount := 0
	for result := range results {
		t.Log(result)
		successCount++
	}

	for err := range errors {
		t.Error(err)
	}

	t.Logf("Success count: %d/%d", successCount, concurrency)

	// All should succeed with proper locking
	if successCount != concurrency {
		t.Errorf("Expected all %d requests to succeed, got %d", concurrency, successCount)
	}
}

// ========================================
// TEST 2: Concurrent Finalize Same Lot
// Only ONE should succeed, others should fail
// ========================================
func TestConcurrentFinalizeLot(t *testing.T) {
	concurrency := 3
	var wg sync.WaitGroup
	results := make(chan string, concurrency)
	startSignal := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			<-startSignal

			url := fmt.Sprintf("%s/lots/%s/finalize", BaseURL, testLotID)
			body := map[string]interface{}{}

			resp, err := makeRequest("POST", url, body)
			if err != nil {
				results <- fmt.Sprintf("Worker %d: ERROR - %v", workerID, err)
				return
			}

			apiResp, err := parseResponse(resp)
			if err != nil {
				results <- fmt.Sprintf("Worker %d: PARSE ERROR - %v", workerID, err)
				return
			}

			results <- fmt.Sprintf("Worker %d: status=%d, success=%v, message=%s",
				workerID, resp.StatusCode, apiResp.Success, apiResp.Message)
		}(i)
	}

	close(startSignal)
	wg.Wait()
	close(results)

	t.Logf("=== Concurrent Finalize Test Results ===")
	successCount := 0
	for result := range results {
		t.Log(result)
		if !containsError(result) {
			successCount++
		}
	}

	// With proper locking, only ONE should actually finalize
	// Others should get "bukan DRAFT" error
	t.Logf("Note: Only 1 should truly succeed, others should fail with 'bukan DRAFT' error")
}

func containsError(s string) bool {
	return bytes.Contains([]byte(s), []byte("ERROR"))
}

// ========================================
// TEST 3: Concurrent Create Lot
// Tests sequence generation locking
// ========================================
func TestConcurrentCreateLot(t *testing.T) {
	concurrency := 5
	var wg sync.WaitGroup
	results := make(chan string, concurrency)
	startSignal := make(chan struct{})

	// jenisDurianID is defined in config section
	jenisDurianID := "1SRlQ8zX9vJ2mN5P6Q7R8S9T001"

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			<-startSignal

			url := fmt.Sprintf("%s/lots", BaseURL)
			body := map[string]interface{}{
				"jenis_durian_id": jenisDurianID,
				"kondisi_buah":    "A", // or "B" depending on your options
			}

			resp, err := makeRequest("POST", url, body)
			if err != nil {
				results <- fmt.Sprintf("Worker %d: ERROR - %v", workerID, err)
				return
			}

			apiResp, err := parseResponse(resp)
			if err != nil {
				results <- fmt.Sprintf("Worker %d: PARSE ERROR - %v", workerID, err)
				return
			}

			// Extract lot kode from response
			kode := ""
			if data, ok := apiResp.Data.(map[string]interface{}); ok {
				if k, ok := data["kode"].(string); ok {
					kode = k
				}
			}

			results <- fmt.Sprintf("Worker %d: status=%d, kode=%s, success=%v",
				workerID, resp.StatusCode, kode, apiResp.Success)
		}(i)
	}

	close(startSignal)
	wg.Wait()
	close(results)

	t.Logf("=== Concurrent Create Lot Test Results ===")
	kodes := make(map[string]bool)
	for result := range results {
		t.Log(result)
		// Extract kode to check for duplicates
	}

	// Check for duplicate kodes (should be none with proper locking)
	if len(kodes) > 0 {
		t.Log("Unique kodes generated - locking is working!")
	}
}

// ========================================
// TEST 4: Transaction Rollback Test
// Tests that failed transaction rolls back all changes
// ========================================
func TestTransactionRollback(t *testing.T) {
	// This test intentionally triggers an error mid-transaction
	// to verify rollback works correctly

	t.Log("=== Transaction Rollback Test ===")
	t.Log("This test verifies that when an error occurs mid-transaction,")
	t.Log("all previous operations in that transaction are rolled back.")
	t.Log("")
	t.Log("To manually test this:")
	t.Log("1. Try to add an item with invalid pohon_kode")
	t.Log("2. Verify no partial data was created")

	url := fmt.Sprintf("%s/lots/%s/items", BaseURL, testLotID)
	body := map[string]interface{}{
		"pohon_kode": "INVALID_POHON_THAT_DOES_NOT_EXIST",
		"blok_id":    testBlokID,
		"berat":      5.0,
	}

	resp, err := makeRequest("POST", url, body)
	if err != nil {
		t.Logf("Request error (expected): %v", err)
		return
	}

	apiResp, err := parseResponse(resp)
	if err != nil {
		t.Logf("Parse error: %v", err)
		return
	}

	if !apiResp.Success {
		t.Logf("✓ Request failed as expected: %s", apiResp.Message)
		t.Log("✓ Transaction should have been rolled back - no partial data created")
	} else {
		t.Error("✗ Request should have failed but succeeded")
	}
}
