package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type Penjualan struct {
	bun.BaseModel `bun:"table:tb_penjualan,alias:penjualan"`

	ID           string     `bun:",pk" json:"id"`
	PengirimanID string     `bun:",notnull" json:"pengiriman_id"`
	BeratTerjual float64    `bun:",notnull" json:"berat_terjual"`
	HargaTotal   float64    `bun:",notnull" json:"harga_total"`
	TipeJual     string     `bun:",notnull" json:"tipe_jual"`
	CreatedAt    time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt    *time.Time `bun:",soft_delete,nullzero" json:"deleted_at,omitempty"`

	Pengiriman *Pengiriman `bun:"rel:belongs-to,join:pengiriman_id=id" json:"pengiriman,omitempty"`
}

func (p *Penjualan) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if p.ID == "" {
			p.ID = ksuid.New().String()
		}
	case *bun.UpdateQuery:
		p.UpdatedAt = time.Now()
	}
	return nil
}
