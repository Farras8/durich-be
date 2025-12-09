package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type TujuanPengiriman struct {
	bun.BaseModel `bun:"table:tb_tujuan_pengiriman,alias:tp"`

	ID        string     `bun:",pk" json:"id"`
	Nama      string     `bun:",notnull" json:"nama"`
	Tipe      string     `bun:",notnull" json:"tipe"`
	Alamat    string     `bun:"" json:"alamat"`
	Kontak    string     `bun:"" json:"kontak"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:",soft_delete,nullzero" json:"deleted_at,omitempty"`
}

func (tp *TujuanPengiriman) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if tp.ID == "" {
			tp.ID = ksuid.New().String()
		}
	case *bun.UpdateQuery:
		tp.UpdatedAt = time.Now()
	}
	return nil
}
