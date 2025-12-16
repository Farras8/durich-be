package tests

import (
	"fmt"
	"sync"
	"testing"
)

// ========================================
// SALES TEST DATA
// Untuk mendapatkan data test, jalankan query:
// SELECT id FROM tb_pengiriman WHERE status = 'RECEIVED' LIMIT 1;
// ========================================
var (
	testReceivedShipmentID = "YOUR_RECEIVED_SHIPMENT_ID" // Shipment yang sudah RECEIVED
)

// ========================================
// TEST 1: Concurrent Create Sales from Same Shipment
// Only ONE should succeed (1 shipment = 1 sales)
// ========================================
func TestConcurrentCreateSalesFromSameShipment(t *testing.T) {
	if testReceivedShipmentID == "YOUR_RECEIVED_SHIPMENT_ID" {
		t.Skip("Please update testReceivedShipmentID with a valid RECEIVED shipment ID")
	}

	concurrency := 3
	var wg sync.WaitGroup
	results := make(chan string, concurrency)
	startSignal := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			<-startSignal

			url := fmt.Sprintf("%s/sales", BaseURL)
			body := map[string]interface{}{
				"pengiriman_id": testReceivedShipmentID,
				"harga_total":   float64(1000000 + workerID*100000),
				"tipe_jual":     "RETAIL",
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

			results <- fmt.Sprintf("Worker %d: status=%d, success=%v, message=%s",
				workerID, resp.StatusCode, apiResp.Success, apiResp.Message)
		}(i)
	}

	close(startSignal)
	wg.Wait()
	close(results)

	t.Logf("=== Concurrent Create Sales Test Results ===")

	successCount := 0
	for result := range results {
		t.Log(result)
		if !containsError(result) && !containsText(result, "status=400") {
			successCount++
		}
	}

	t.Logf("Success count: %d", successCount)
	t.Log("Note: Only 1 should succeed, others should fail (1 shipment = 1 sales)")
}

// ========================================
// TEST 2: Concurrent Update Same Sales
// Tests locking during update
// ========================================
func TestConcurrentUpdateSameSales(t *testing.T) {
	t.Log("=== Concurrent Update Same Sales Test ===")
	t.Log("This test requires an existing sales record")
	t.Log("To test manually:")
	t.Log("1. Get an existing sales ID")
	t.Log("2. Send multiple PUT requests concurrently")
	t.Log("Expected: All should succeed but with locking (serialized)")
}

// Helper function
func containsText(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsBytesSubstring([]byte(s), []byte(substr)))
}

func containsBytesSubstring(s, substr []byte) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if string(s[i:i+len(substr)]) == string(substr) {
			return true
		}
	}
	return false
}
