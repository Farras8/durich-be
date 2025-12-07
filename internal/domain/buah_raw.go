package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type BuahRaw struct {
	bun.BaseModel `bun:"table:tb_buah_raw,alias:buah_raw"`

	ID                string       `bun:",pk" json:"id"`
	KodeBuah          string       `bun:",notnull" json:"kode_buah"`
	JenisDurian       string       `bun:",notnull" json:"jenis_durian"`
	PohonPanen        *string      `bun:"," json:"pohon_panen"`
	TglPanen          string       `bun:",notnull" json:"tgl_panen"`
	
	LotID     *string  `bun:",nullzero" json:"lot_id,omitempty"`
	BlokID    *string  `bun:",nullzero" json:"blok_id,omitempty"`
	Berat     float64  `bun:",default:0" json:"berat"`

	CreatedAt         time.Time    `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time    `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt         *time.Time   `bun:"," json:"deleted_at,omitempty"`
	
	JenisDurianDetail *JenisDurian `bun:"rel:belongs-to,join:jenis_durian=id" json:"jenis_durian_detail,omitempty"`
	PohonPanenDetail  *Pohon       `bun:"rel:belongs-to,join:pohon_panen=id" json:"pohon_panen_detail,omitempty"`
	Lot               *StokLot     `bun:"rel:belongs-to,join:lot_id=id" json:"lot,omitempty"`
	Blok              *Blok        `bun:"rel:belongs-to,join:blok_id=id" json:"blok,omitempty"`
}

func (m *BuahRaw) BeforeAppendModel(_ context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if m.ID == "" {
			m.ID = ksuid.New().String()
		}
	case *bun.UpdateQuery:
		m.UpdatedAt = time.Now()
	}
	return nil
}
