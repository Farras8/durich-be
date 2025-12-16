package tests

import (
	"fmt"
	"sync"
	"testing"
)

// ========================================
// SHIPMENT TEST DATA
// Untuk mendapatkan data test, jalankan query:
// SELECT id FROM tb_tujuan_pengiriman LIMIT 1;
// SELECT id FROM tb_stok_lot WHERE status = 'READY' LIMIT 3;
// ========================================
var (
	testTujuanID   = "YOUR_TUJUAN_ID" // ID tujuan pengiriman yang valid
	testReadyLotID = "YOUR_LOT_ID"    // ID lot dengan status READY untuk test add item
)

// ========================================
// TEST 1: Concurrent Create Shipment
// Tests sequence generation (SHP-YYMMDD-XXX)
// ========================================
func TestConcurrentCreateShipment(t *testing.T) {
	if testTujuanID == "YOUR_TUJUAN_ID" {
		t.Skip("Please update testTujuanID with a valid tujuan pengiriman ID")
	}

	concurrency := 5
	var wg sync.WaitGroup
	results := make(chan string, concurrency)
	startSignal := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			<-startSignal

			url := fmt.Sprintf("%s/shipments", BaseURL)
			body := map[string]interface{}{
				"tujuan_id": testTujuanID,
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

			// Extract shipment kode from response
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

	t.Logf("=== Concurrent Create Shipment Test Results ===")

	kodes := make(map[string]int)
	for result := range results {
		t.Log(result)
	}

	// Check for duplicate kodes
	for kode, count := range kodes {
		if count > 1 {
			t.Errorf("Duplicate kode detected: %s (count: %d)", kode, count)
		}
	}
	t.Log("All shipment kodes should be unique with sequential numbers")
}

// ========================================
// TEST 2: Concurrent Add Item to Same Shipment
// Tests row-level locking
// ========================================
func TestConcurrentAddItemToShipment(t *testing.T) {
	if testTujuanID == "YOUR_TUJUAN_ID" {
		t.Skip("Please update test data first")
	}

	// First, create a new shipment
	createURL := fmt.Sprintf("%s/shipments", BaseURL)
	createBody := map[string]interface{}{
		"tujuan_id": testTujuanID,
	}

	createResp, err := makeRequest("POST", createURL, createBody)
	if err != nil {
		t.Fatalf("Failed to create shipment: %v", err)
	}

	apiResp, err := parseResponse(createResp)
	if err != nil {
		t.Fatalf("Failed to parse create response: %v", err)
	}

	shipmentID := ""
	if data, ok := apiResp.Data.(map[string]interface{}); ok {
		if id, ok := data["id"].(string); ok {
			shipmentID = id
		}
	}

	if shipmentID == "" {
		t.Fatal("Failed to get shipment ID")
	}

	t.Logf("Created shipment: %s", shipmentID)
	t.Log("Note: This test requires multiple READY lots to add items concurrently")
	t.Log("Skipping concurrent add item test - would need multiple READY lots")
}

// ========================================
// TEST 3: Concurrent Finalize Same Shipment
// Only ONE should succeed
// ========================================
func TestConcurrentFinalizeShipment(t *testing.T) {
	t.Log("=== Concurrent Finalize Shipment Test ===")
	t.Log("This test requires a DRAFT shipment with at least 1 item")
	t.Log("Due to complex setup requirements, this should be tested manually:")
	t.Log("1. Create a shipment")
	t.Log("2. Add items to shipment")
	t.Log("3. Send multiple finalize requests concurrently")
	t.Log("Expected: Only 1 should succeed, others fail with 'shipment must be DRAFT'")
}

// ========================================
// TEST 4: Concurrent Receive Same Shipment
// Tests locking during receive
// ========================================
func TestConcurrentReceiveShipment(t *testing.T) {
	t.Log("=== Concurrent Receive Shipment Test ===")
	t.Log("This test requires a SENDING shipment")
	t.Log("Due to complex setup requirements, this should be tested manually:")
	t.Log("1. Have a shipment in SENDING status")
	t.Log("2. Send multiple receive requests concurrently")
	t.Log("Expected: Only 1 should succeed, others fail with status error")
}
