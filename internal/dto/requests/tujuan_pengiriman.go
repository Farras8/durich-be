package requests

type CreateTujuanPengirimanRequest struct {
	Nama   string `json:"nama" binding:"required"`
	Tipe   string `json:"tipe" binding:"required"`
	Alamat string `json:"alamat"`
	Kontak string `json:"kontak"`
}

type UpdateTujuanPengirimanRequest struct {
	Nama   string `json:"nama" binding:"required"`
	Tipe   string `json:"tipe" binding:"required"`
	Alamat string `json:"alamat"`
	Kontak string `json:"kontak"`
}
