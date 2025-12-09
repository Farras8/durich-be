package requests

import (
	"time"
)

type ShipmentCreateRequest struct {
	TujuanID string    `json:"tujuan_id" binding:"required"`
	TglKirim time.Time `json:"tgl_kirim"`
}

type ShipmentAddItemRequest struct {
	LotID string `json:"lot_id" binding:"required"`
}

type ShipmentRemoveItemRequest struct {
	DetailID string `json:"detail_id" binding:"required"`
}

type ShipmentUpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

type ShipmentReceiveRequest struct {
	ReceivedDate time.Time `json:"received_date" binding:"required"`
	Details      []struct {
		LotID         string  `json:"lot_id" binding:"required"`
		BeratDiterima float64 `json:"berat_diterima" binding:"required,min=0"`
	} `json:"details" binding:"required,dive"`
}
