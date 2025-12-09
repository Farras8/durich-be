package response

import (
	"durich-be/internal/domain"
	"time"
)

type TujuanPengirimanResponse struct {
	ID        string     `json:"id"`
	Nama      string     `json:"nama"`
	Tipe      string     `json:"tipe"`
	Alamat    string     `json:"alamat"`
	Kontak    string     `json:"kontak"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewTujuanPengirimanResponse(t *domain.TujuanPengiriman) *TujuanPengirimanResponse {
	return &TujuanPengirimanResponse{
		ID:        t.ID,
		Nama:      t.Nama,
		Tipe:      t.Tipe,
		Alamat:    t.Alamat,
		Kontak:    t.Kontak,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func NewTujuanPengirimanListResponse(tujuans []domain.TujuanPengiriman) []*TujuanPengirimanResponse {
	var responses []*TujuanPengirimanResponse
	for _, t := range tujuans {
		responses = append(responses, NewTujuanPengirimanResponse(&t))
	}
	return responses
}
