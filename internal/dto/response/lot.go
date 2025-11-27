package response

import "time"

type LotResponse struct {
	ID              string    `json:"id"`
	JenisDurianID   string    `json:"jenis_durian_id"`
	JenisDurianNama string    `json:"jenis_durian_nama"`
	KondisiBuah     string    `json:"kondisi_buah"`
	BeratAwal   float64   `json:"berat_awal"`
	QtyAwal     int       `json:"qty_awal"`
	BeratSisa   float64   `json:"berat_sisa"`
	QtySisa     int       `json:"qty_sisa"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type LotDetailResponse struct {
	Header LotResponse        `json:"header"`
	Items  []LotItemResponse  `json:"items"`
}

type LotItemResponse struct {
	ID        string `json:"id"`
	KodeBuah  string `json:"kode_buah"`
	TglPanen  string `json:"tgl_panen"`
	AsalBlok  string `json:"asal_blok"`
}

type LotAddItemsResponse struct {
	CurrentQty int `json:"current_qty"`
}

type LotFinalizeResponse struct {
	ID         string  `json:"id"`
	QtyTotal   int     `json:"qty_total"`
	BeratTotal float64 `json:"berat_total"`
	Status     string  `json:"status"`
}
