package response

import "time"

type TraceLotResponse struct {
	LotInfo LotTraceInfo     `json:"lot_info"`
	Fruits  []FruitTraceInfo `json:"fruits"`
}

type LotTraceInfo struct {
	ID          string    `json:"id"`
	JenisDurian string    `json:"jenis_durian"`
	KondisiBuah string    `json:"kondisi_buah"`
	QtyTotal    int       `json:"qty_total"`
	BeratTotal  float64   `json:"berat_total"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type FruitTraceInfo struct {
	BuahRawID     string          `json:"buah_raw_id"`
	KodeBuah      string          `json:"kode_buah"`
	TglPanen      string          `json:"tgl_panen"`
	LokasiPanen   LokasiTraceInfo `json:"lokasi_panen"`
	AddedToLotAt  time.Time       `json:"added_to_lot_at"`
}

type LokasiTraceInfo struct {
	Company string `json:"company"`
	Estate  string `json:"estate"`
	Divisi  string `json:"divisi"`
	Blok    string `json:"blok"`
	Pohon   string `json:"pohon"`
}

type TraceFruitResponse struct {
	FruitInfo FruitDetailInfo `json:"fruit_info"`
	Journey   FruitJourney    `json:"journey"`
}

type FruitDetailInfo struct {
	ID          string          `json:"id"`
	KodeBuah    string          `json:"kode_buah"`
	JenisDurian string          `json:"jenis_durian"`
	TglPanen    string          `json:"tgl_panen"`
	LokasiPanen LokasiTraceInfo `json:"lokasi_panen"`
}

type FruitJourney struct {
	Lot      *LotJourneyInfo      `json:"lot,omitempty"`
	Shipment *ShipmentJourneyInfo `json:"shipment,omitempty"`
	Sales    *SalesJourneyInfo    `json:"sales,omitempty"`
}

type LotJourneyInfo struct {
	ID          string    `json:"id"`
	KondisiBuah string    `json:"kondisi_buah"`
	JoinedAt    time.Time `json:"joined_at"`
	Status      string    `json:"status"`
}

type ShipmentJourneyInfo struct {
	ID       string    `json:"id"`
	Tujuan   string    `json:"tujuan"`
	TglKirim time.Time `json:"tgl_kirim"`
	Status   string    `json:"status"`
}

type SalesJourneyInfo struct {
	ID           string    `json:"id"`
	HargaTotal   float64   `json:"harga_total"`
	TipeJual     string    `json:"tipe_jual"`
	TglTransaksi time.Time `json:"tgl_transaksi"`
}

type TraceShipmentResponse struct {
	ShipmentInfo      ShipmentTraceInfo       `json:"shipment_info"`
	BreakdownByLokasi []BreakdownByLokasi     `json:"breakdown_by_location"`
	DetailedFruits    []DetailedFruitInfo     `json:"detailed_fruits"`
}

type ShipmentTraceInfo struct {
	ID          string    `json:"id"`
	Tujuan      string    `json:"tujuan"`
	TglKirim    time.Time `json:"tgl_kirim"`
	Status      string    `json:"status"`
	TotalQty    int       `json:"total_qty"`
	TotalBerat  float64   `json:"total_berat"`
}

type BreakdownByLokasi struct {
	Lokasi       string   `json:"lokasi"`
	JenisDurian  string   `json:"jenis_durian"`
	Grade        string   `json:"grade"`
	Qty          int      `json:"qty"`
	Berat        float64  `json:"berat"`
	PohonSources []string `json:"pohon_sources"`
}

type DetailedFruitInfo struct {
	BuahRawID string `json:"buah_raw_id"`
	KodeBuah  string `json:"kode_buah"`
	Lokasi    string `json:"lokasi"`
	Pohon     string `json:"pohon"`
	TglPanen  string `json:"tgl_panen"`
}
