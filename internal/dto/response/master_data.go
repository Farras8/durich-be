package response

import "time"

type CompanyResponse struct {
	ID        string    `json:"id"`
	Kode      string    `json:"kode"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EstateResponse struct {
	ID        string           `json:"id"`
	Kode      string           `json:"kode"`
	Nama      string           `json:"nama"`
	CompanyID string           `json:"company_id"`
	Company   *CompanyResponse `json:"company,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type DivisiResponse struct {
	ID        string          `json:"id"`
	Kode      string          `json:"kode"`
	Nama      string          `json:"nama"`
	EstateID  string          `json:"estate_id"`
	Estate    *EstateResponse `json:"estate,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type BlokResponse struct {
	ID          string          `json:"id"`
	Kode        string          `json:"kode"`
	NamaBlok    string          `json:"nama_blok"`
	KodeLengkap string          `json:"kode_lengkap,omitempty"`
	DivisiID    string          `json:"divisi_id"`
	Divisi      *DivisiResponse `json:"divisi,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type JenisDurianResponse struct {
	ID        string    `json:"id"`
	Kode      string    `json:"kode"`
	NamaJenis string    `json:"nama_jenis"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PohonResponse struct {
	ID          string    `json:"id"`
	Kode        string    `json:"kode"`
	Nama        string    `json:"nama"`
	KodeLengkap string    `json:"kode_lengkap,omitempty"`
	BlokID      *string   `json:"blok_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
