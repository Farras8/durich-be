package requests

import "time"

type ShipmentCreateRequest struct {
	Tujuan   string    `json:"tujuan" binding:"required"`
	TglKirim time.Time `json:"tgl_kirim"`
}

type ShipmentAddItemRequest struct {
	LotID string  `json:"lot_id" binding:"required"`
	Qty   int     `json:"qty" binding:"required,min=1"`
	Berat float64 `json:"berat" binding:"required,min=0.01"`
}

type ShipmentRemoveItemRequest struct {
	DetailID string `json:"detail_id" binding:"required"`
}
