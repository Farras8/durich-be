package requests

type SalesCreateRequest struct {
	PengirimanID string  `json:"pengiriman_id" binding:"required"`
	HargaTotal   float64 `json:"harga_total" binding:"required,min=0"`
	TipeJual     string  `json:"tipe_jual" binding:"required"`
}

type SalesUpdateRequest struct {
	HargaTotal float64 `json:"harga_total" binding:"omitempty,min=0"`
	TipeJual   string  `json:"tipe_jual"`
}
