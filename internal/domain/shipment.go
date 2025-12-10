package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type Pengiriman struct {
	bun.BaseModel `bun:"table:tb_pengiriman,alias:p"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string            `bun:",notnull" json:"kode"`
	TujuanID  string            `bun:",notnull" json:"tujuan_id"`
	Tujuan    string            `bun:",notnull" json:"tujuan"`
	TglKirim  time.Time         `bun:",notnull" json:"tgl_kirim"`
	Status    string            `bun:",default:'SENT'" json:"status"`
	ReceivedAt *time.Time       `bun:",nullzero" json:"received_at,omitempty"`
	CreatedBy string            `bun:",notnull" json:"created_by"`
	CreatedAt time.Time         `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time         `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time        `bun:",soft_delete,nullzero" json:"deleted_at,omitempty"`

	Details      []PengirimanDetail `bun:"rel:has-many,join:id=pengiriman_id" json:"details,omitempty"`
	Creator      *User              `bun:"rel:belongs-to,join:created_by=id" json:"creator,omitempty"`
	TujuanDetail *TujuanPengiriman  `bun:"rel:belongs-to,join:tujuan_id=id" json:"tujuan_detail,omitempty"`
}

type PengirimanDetail struct {
	bun.BaseModel `bun:"table:tb_pengiriman_detail,alias:pd"`

	ID           string    `bun:",pk" json:"id"`
	PengirimanID string    `bun:",notnull" json:"pengiriman_id"`
	LotSumberID  string    `bun:",notnull" json:"lot_sumber_id"`
	QtyAmbil     int       `bun:",notnull" json:"qty_ambil"`
	BeratAmbil   float64   `bun:",notnull" json:"berat_ambil"`
	CreatedAt    time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Pengiriman *Pengiriman `bun:"rel:belongs-to,join:pengiriman_id=id" json:"pengiriman,omitempty"`
	Lot        *StokLot    `bun:"rel:belongs-to,join:lot_sumber_id=id" json:"lot,omitempty"`
}

func (p *PengirimanDetail) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if p.ID == "" {
			p.ID = ksuid.New().String()
		}
	}
	return nil
}

func (p *Pengiriman) BeforeAppendModel(_ context.Context, query bun.Query) error {
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
