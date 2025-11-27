package response

import (
	"durich-be/internal/domain"
	"time"
)

type ShipmentResponse struct {
	ID         string    `json:"id"`
	Tujuan     string    `json:"tujuan"`
	TglKirim   time.Time `json:"tgl_kirim"`
	Status     string    `json:"status"`
	CreatedBy  string    `json:"created_by"`
	TotalItems int       `json:"total_items"`
	TotalBerat float64   `json:"total_berat"`
	CreatedAt  time.Time `json:"created_at"`
}

type ShipmentItemResponse struct {
	ID             string  `json:"id"`
	LotID          string  `json:"lot_id"`
	JenisDurian    string  `json:"jenis_durian"`
	QtyAmbil       int     `json:"qty_ambil"`
	BeratAmbil     float64 `json:"berat_ambil"`
}

type ShipmentDetailResponse struct {
	Header ShipmentResponse       `json:"header"`
	Items  []ShipmentItemResponse `json:"items"`
}

func NewShipmentResponse(p *domain.Pengiriman) ShipmentResponse {
	totalItems := 0
	totalBerat := 0.0

	for _, d := range p.Details {
		totalItems++ // or d.QtyAmbil depending on definition, usually items count = rows
		totalBerat += d.BeratAmbil
	}

	// If counting total quantity of fruits instead of rows:
	// totalItems = 0
	// for _, d := range p.Details { totalItems += d.QtyAmbil }
	// Let's stick to rows for "items" or maybe query logic handles totals.
	// Spec says: "total_items: 5 // Jumlah lot/detail". So it's row count.
	// Wait, spec example says "total_items: 5". If it's number of lots, it's row count.

	return ShipmentResponse{
		ID:         p.ID,
		Tujuan:     p.Tujuan,
		TglKirim:   p.TglKirim,
		Status:     p.Status,
		CreatedBy:  p.CreatedBy,
		TotalItems: len(p.Details),
		TotalBerat: totalBerat,
		CreatedAt:  p.CreatedAt,
	}
}
