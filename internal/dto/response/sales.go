package response

import (
	"durich-be/internal/domain"
	"time"
)

type SalesResponse struct {
	ID           string    `json:"id"`
	TglTransaksi time.Time `json:"tgl_transaksi"`
	PengirimanID string    `json:"pengiriman_id"`
	BeratTerjual float64   `json:"berat_terjual"`
	HargaTotal   float64   `json:"harga_total"`
	TipeJual     string    `json:"tipe_jual"`
}

type SalesDetailResponse struct {
	ID             string                    `json:"id"`
	InfoPenjualan  SalesInfoResponse         `json:"info_penjualan"`
	InfoPengiriman SalesShipmentInfoResponse `json:"info_pengiriman"`
}

type SalesInfoResponse struct {
	HargaTotal   float64   `json:"harga_total"`
	BeratTerjual float64   `json:"berat_terjual"`
	TipeJual     string    `json:"tipe_jual"`
	CreatedAt    time.Time `json:"created_at"`
}

type SalesShipmentInfoResponse struct {
	ID      string                 `json:"id"`
	Tujuan  string                 `json:"tujuan"`
	Status  string                 `json:"status"`
	Details []ShipmentItemResponse `json:"details"`
}

func NewSalesResponse(s *domain.Penjualan) SalesResponse {
	return SalesResponse{
		ID:           s.ID,
		TglTransaksi: s.CreatedAt,
		PengirimanID: s.PengirimanID,
		BeratTerjual: s.BeratTerjual,
		HargaTotal:   s.HargaTotal,
		TipeJual:     s.TipeJual,
	}
}
