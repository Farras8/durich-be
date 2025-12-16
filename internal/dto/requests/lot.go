package requests

type LotCreateRequest struct {
	JenisDurianID string `json:"jenis_durian_id" binding:"required"`
	KondisiBuah   string `json:"kondisi_buah" binding:"required"`
}

type LotAddItemsRequest struct {
	PohonKode string  `json:"pohon_kode" binding:"required"` 
	BlokID    string  `json:"blok_id" binding:"required"`    
	Berat     float64 `json:"berat" binding:"required,gt=0"`
}

type LotRemoveItemRequest struct {
	BuahRawID string `json:"buah_raw_id" binding:"required"`
}

type LotFinalizeRequest struct {
}
