package requests

type LotCreateRequest struct {
	JenisDurianID string `json:"jenis_durian_id" binding:"required"`
	KondisiBuah   string `json:"kondisi_buah" binding:"required"`
}

type LotAddItemsRequest struct {
	BuahRawIDs []string `json:"buah_raw_ids" binding:"required,min=1,dive,required"`
}

type LotRemoveItemRequest struct {
	BuahRawID string `json:"buah_raw_id" binding:"required"`
}

type LotFinalizeRequest struct {
	BeratAwal float64 `json:"berat_awal" binding:"required,gt=0"`
}
