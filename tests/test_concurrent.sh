#!/bin/bash
# =====================================================
# Script untuk test concurrent requests ke API
# =====================================================
# Usage: ./test_concurrent.sh
# =====================================================

# KONFIGURASI - Update nilai-nilai ini!
BASE_URL="http://localhost:8081/v1"
TOKEN="YOUR_JWT_TOKEN_HERE"  # Ganti dengan JWT token yang valid
LOT_ID="YOUR_LOT_ID"         # Ganti dengan ID lot yang statusnya DRAFT
BLOK_ID="YOUR_BLOK_ID"       # Ganti dengan ID blok yang valid
POHON_KODE="P001"

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  CONCURRENT API TEST SCRIPT${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# =====================================================
# TEST 1: Concurrent AddItems ke Lot yang sama
# =====================================================
echo -e "${GREEN}[TEST 1] Concurrent AddItems to Same Lot${NC}"
echo "Sending 5 concurrent requests..."

for i in {1..5}; do
  curl -s -X POST "${BASE_URL}/lots/${LOT_ID}/items" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{\"pohon_kode\":\"${POHON_KODE}\",\"blok_id\":\"${BLOK_ID}\",\"berat\":${i}.5}" \
    -w "\nWorker ${i}: HTTP %{http_code} - Time: %{time_total}s\n" &
done

wait
echo -e "${GREEN}[TEST 1] Complete!${NC}"
echo ""

# =====================================================
# TEST 2: Concurrent Create Lot
# =====================================================
echo -e "${GREEN}[TEST 2] Concurrent Create Lot${NC}"
echo "Sending 3 concurrent requests..."

JENIS_DURIAN_ID="YOUR_JENIS_DURIAN_ID"  # Update this

for i in {1..3}; do
  curl -s -X POST "${BASE_URL}/lots" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{\"jenis_durian_id\":\"${JENIS_DURIAN_ID}\",\"kondisi_buah\":\"A\"}" \
    -w "\nWorker ${i}: HTTP %{http_code} - Time: %{time_total}s\n" &
done

wait
echo -e "${GREEN}[TEST 2] Complete!${NC}"
echo ""

# =====================================================
# TEST 3: Sequential dengan delay (untuk perbandingan)
# =====================================================
echo -e "${GREEN}[TEST 3] Sequential Requests (for comparison)${NC}"
echo "Sending 3 sequential requests..."

for i in {1..3}; do
  echo -n "Request ${i}: "
  curl -s -X POST "${BASE_URL}/lots" \
    -H "Authorization: Bearer ${TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{\"jenis_durian_id\":\"${JENIS_DURIAN_ID}\",\"kondisi_buah\":\"A\"}" \
    -w "HTTP %{http_code} - Time: %{time_total}s\n"
done

echo -e "${GREEN}[TEST 3] Complete!${NC}"
echo ""

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  All tests completed!${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""
echo "Tips for analysis:"
echo "1. Check if all requests succeeded (HTTP 200/201)"
echo "2. Verify no duplicate codes were generated"
echo "3. Compare response times between concurrent vs sequential"
echo ""
