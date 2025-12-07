package response

import (
	"durich-be/internal/domain"
	"time"
)

type ShipmentResponse struct {
	ID         string    `json:"id"`
	Kode       string    `json:"kode"`
	Tujuan     string    `json:"tujuan"`
	TglKirim   time.Time `json:"tgl_kirim"`
	Status     string    `json:"status"`
	CreatedBy  string    `json:"created_by"`
	TotalItems int       `json:"total_items"`
	TotalBerat float64   `json:"total_berat"`
	CreatedAt  time.Time `json:"created_at"`
}

type ShipmentItemResponse struct {
	ID          string  `json:"id"`
	LotID       string  `json:"lot_id"`
	KodeLot     string  `json:"kode_lot"` // Added field
	JenisDurian string  `json:"jenis_durian"`
	Grade       string  `json:"grade"`
	QtyAmbil    int     `json:"qty_ambil"`
	BeratAmbil  float64 `json:"berat_ambil"`
}

type ShipmentDetailResponse struct {
	Header ShipmentResponse       `json:"header"`
	Items  []ShipmentItemResponse `json:"items"`
}

func NewShipmentResponse(p *domain.Pengiriman) ShipmentResponse {
	totalItems := 0
	totalBerat := 0.0

	for _, d := range p.Details {
		totalItems++
		totalBerat += d.BeratAmbil
	}

	createdBy := p.CreatedBy
	if p.Creator != nil {
		createdBy = p.Creator.Email
	}

	return ShipmentResponse{
		ID:         p.ID,
		Kode:       p.Kode,
		Tujuan:     p.Tujuan,
		TglKirim:   p.TglKirim,
		Status:     p.Status,
		CreatedBy:  createdBy,
		TotalItems: len(p.Details),
		TotalBerat: totalBerat,
		CreatedAt:  p.CreatedAt,
	}
}
