package requests

type CompanyCreateRequest struct {
	Kode string `json:"kode" binding:"required,max=5"`
	Nama string `json:"nama" binding:"required"`
}

type CompanyUpdateRequest struct {
	Nama string `json:"nama" binding:"required"`
}

type EstateCreateRequest struct {
	Kode      string `json:"kode" binding:"required,max=5"`
	Nama      string `json:"nama" binding:"required"`
	CompanyID string `json:"company_id" binding:"required"`
}

type EstateUpdateRequest struct {
	Nama      string `json:"nama" binding:"required"`
	CompanyID string `json:"company_id" binding:"required"`
}

type DivisiCreateRequest struct {
	Kode     string `json:"kode" binding:"required,max=5"`
	Nama     string `json:"nama" binding:"required"`
	EstateID string `json:"estate_id" binding:"required"`
}

type DivisiUpdateRequest struct {
	Nama     string `json:"nama" binding:"required"`
	EstateID string `json:"estate_id" binding:"required"`
}

type BlokCreateRequest struct {
	Kode     string `json:"kode" binding:"required"`
	NamaBlok string `json:"nama_blok" binding:"required"`
	DivisiID string `json:"divisi_id" binding:"required"`
}

type BlokUpdateRequest struct {
	NamaBlok string `json:"nama_blok" binding:"required"`
	DivisiID string `json:"divisi_id" binding:"required"`
}

type JenisDurianCreateRequest struct {
	Kode      string `json:"kode" binding:"required"`
	NamaJenis string `json:"nama_jenis" binding:"required"`
}

type JenisDurianUpdateRequest struct {
	NamaJenis string `json:"nama_jenis" binding:"required"`
}

type PohonCreateRequest struct {
	Kode string `json:"kode" binding:"required,max=5"`
	Nama string `json:"nama" binding:"required"`
}

type PohonUpdateRequest struct {
	Nama string `json:"nama" binding:"required"`
}
