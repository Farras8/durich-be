package requests

type BuahRawBulkCreateItem struct {
	JenisDurianID string  `json:"jenis_durian_id" binding:"required"`
	PohonPanenID  *string `json:"pohon_panen_id"`
	Jumlah        int     `json:"jumlah" binding:"required,min=1"`
}

type BuahRawBulkCreateRequest struct {
	TglPanen    string                  `json:"tgl_panen"`
	BlokPanenID string                  `json:"blok_panen_id" binding:"required"`
	Items       []BuahRawBulkCreateItem `json:"items" binding:"required,min=1,dive"`
}

type BuahRawCreateRequest struct {
	TglPanen      string  `json:"tgl_panen"`
	BlokPanenID   string  `json:"blok_panen_id" binding:"required"`
	JenisDurianID string  `json:"jenis_durian_id" binding:"required"`
	PohonPanenID  *string `json:"pohon_panen_id"`
}

type BuahRawUpdateRequest struct {
	TglPanen      string  `json:"tgl_panen"`
	BlokPanenID   string  `json:"blok_panen_id"`
	PohonPanenID  *string `json:"pohon_panen_id"`
	JenisDurianID string  `json:"jenis_durian_id"`
}
