package response

type LokasiPanen struct {
	KodeLengkap string `json:"kode_lengkap"`
	BlokID      string `json:"blok_id"`
	BlokNama    string `json:"blok_nama"`
	DivisiNama  string `json:"divisi_nama"`
	EstateNama  string `json:"estate_nama"`
	CompanyNama string `json:"company_nama"`
}

type JenisDurianDetail struct {
	ID   string `json:"id"`
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type BuahRawResponse struct {
	ID          string            `json:"id"`
	KodeBuah    string            `json:"kode_buah"`
	JenisDurian JenisDurianDetail `json:"jenis_durian"`
	LokasiPanen LokasiPanen       `json:"lokasi_panen"`
	PohonPanen  *string           `json:"pohon_panen"`
	KodePohon   string            `json:"kode_pohon"` // Added field
	LotKode     *string           `json:"kode_lot"`   // Added field
	TglPanen    string            `json:"tgl_panen"`
	CreatedAt   string            `json:"created_at"`
}

type PaginationMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

type PaginationResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
