# Database Transaction & Locking Test Guide

Panduan untuk testing implementasi database transaction dan locking di durich-be.

## 📁 Test Files

| File | Deskripsi |
|------|-----------|
| `tests/concurrent_test.go` | Go test untuk Lot concurrent requests |
| `tests/shipment_concurrent_test.go` | Go test untuk Shipment concurrent requests |
| `tests/sales_concurrent_test.go` | Go test untuk Sales concurrent requests |
| `tests/test_concurrent.sh` | Bash script untuk quick testing dengan curl |

## 🔧 Setup Sebelum Testing

### 1. Update Konfigurasi

Edit file `tests/concurrent_test.go` dan ganti nilai-nilai berikut:

```go
const (
    BaseURL   = "http://localhost:8081/v1"
    AuthToken = "YOUR_ACCESS_TOKEN"  // Dapat dari login API
)

var (
    testLotID     = "YOUR_LOT_ID"      // ID lot dengan status DRAFT
    testBlokID    = "YOUR_BLOK_ID"     // ID blok valid
    testPohonKode = "P001"             // Kode pohon valid
)
```

### 1b. Cara Mendapatkan Access Token

```bash
# Login untuk mendapatkan access token
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your_email","password":"your_password"}'

# Copy nilai "access_token" dari response JSON
```

### 2. Dapatkan Data Test

```sql
-- Cari lot dengan status DRAFT
SELECT id, kode, status FROM tb_stok_lot WHERE status = 'DRAFT' LIMIT 1;

-- Cari blok dan pohon
SELECT b.id as blok_id, p.kode as pohon_kode 
FROM tb_blok b 
JOIN tb_pohon p ON p.blok_id = b.id 
LIMIT 1;

-- Cari jenis durian
SELECT id, kode FROM tb_jenis_durian LIMIT 1;
```

## 🚀 Cara Running Tests

### Option 1: Go Test (Recommended)

```bash
cd c:\Users\mfarr\OneDrive\Dokumen\Durich\durich-be

# Run semua tests
go test -v ./tests/...
```

#### Lot Tests
```bash
go test -v ./tests/... -run TestConcurrentAddItemsToLot
go test -v ./tests/... -run TestConcurrentFinalizeLot
go test -v ./tests/... -run TestConcurrentCreateLot
go test -v ./tests/... -run TestTransactionRollback
```

#### Shipment Tests
```bash
go test -v ./tests/... -run TestConcurrentCreateShipment
go test -v ./tests/... -run TestConcurrentAddItemToShipment
go test -v ./tests/... -run TestConcurrentFinalizeShipment
go test -v ./tests/... -run TestConcurrentReceiveShipment
```

#### Sales Tests
```bash
go test -v ./tests/... -run TestConcurrentCreateSalesFromSameShipment
go test -v ./tests/... -run TestConcurrentUpdateSameSales
```


### Option 2: Bash Script

```bash
cd c:\Users\mfarr\OneDrive\Dokumen\Durich\durich-be\tests

# Beri permission executable
chmod +x test_concurrent.sh

# Run script
./test_concurrent.sh
```

### Option 3: Manual dari Postman

1. Buka 3-5 tab Postman
2. Siapkan request POST ke `/lots/{lot_id}/items`
3. Klik Send secara bersamaan di semua tab
4. Perhatikan response dan timing

## ✅ Expected Results

### Test 1: Concurrent AddItems
| Scenario | Expected |
|----------|----------|
| 5 request bersamaan | Semua sukses, sequence berbeda |
| Waktu total | Lebih lama dari sequential (karena locking) |

### Test 2: Concurrent Finalize
| Scenario | Expected |
|----------|----------|
| 3 request bersamaan | 1 sukses, 2 gagal "bukan DRAFT" |

### Test 3: Concurrent Create Lot
| Scenario | Expected |
|----------|----------|
| 5 request bersamaan | Semua sukses, kode unik |

### Test 4: Transaction Rollback
| Scenario | Expected |
|----------|----------|
| Invalid pohon_kode | Error, tidak ada data tersimpan |

## 🔍 Troubleshooting

### Semua Request Gagal
- Pastikan server berjalan (`go run main.go`)
- Periksa JWT token masih valid
- Verifikasi lot ID ada dan statusnya DRAFT

### Duplicate Kode Tercipta
- Ini menandakan locking TIDAK bekerja
- Periksa apakah `FOR UPDATE` ada di query
- Pastikan query berjalan dalam transaction

### Request Timeout
- Periksa apakah ada deadlock
- Tambah timeout di context
- Check database connections

## 📊 Cara Verifikasi di Database

```sql
-- Cek buah yang baru dibuat (untuk AddItems test)
SELECT * FROM tb_buah_raw ORDER BY created_at DESC LIMIT 10;

-- Cek lot codes (untuk Create test)
SELECT kode, created_at FROM tb_stok_lot ORDER BY created_at DESC LIMIT 10;

-- Cek tidak ada duplicate
SELECT kode, COUNT(*) 
FROM tb_stok_lot 
GROUP BY kode 
HAVING COUNT(*) > 1;
```
