package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type StokLot struct {
	bun.BaseModel `bun:"table:tb_stok_lot,alias:stok_lot"`

	ID           string     `bun:",pk" json:"id"`
	JenisDurian  string     `bun:",notnull" json:"jenis_durian"`
	KondisiBuah  string     `bun:",notnull" json:"kondisi_buah"`
	BeratAwal    float64    `bun:",default:0" json:"berat_awal"`
	QtyAwal      int        `bun:",default:0" json:"qty_awal"`
	BeratSisa    float64    `bun:",default:0" json:"berat_sisa"`
	QtySisa      int        `bun:",default:0" json:"qty_sisa"`
	Status       string     `bun:",default:'DRAFT'" json:"status"`
	CreatedAt    time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt    *time.Time `bun:"" json:"deleted_at,omitempty"`

	Items []LotDetail `bun:"rel:has-many,join:id=lot_id" json:"items,omitempty"`

	JenisDurianDetail *JenisDurian `bun:"rel:belongs-to,join:jenis_durian=id" json:"jenis_durian_detail,omitempty"`
}

type LotDetail struct {
	bun.BaseModel `bun:"table:tb_lot_detail,alias:lot_detail"`

	LotID      string    `bun:",pk" json:"lot_id"`
	BuahRawID  string    `bun:",pk" json:"buah_raw_id"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`

	Lot     *StokLot `bun:"rel:belongs-to,join:lot_id=id" json:"lot,omitempty"`
	BuahRaw *BuahRaw `bun:"rel:belongs-to,join:buah_raw_id=id" json:"buah_raw,omitempty"`
}

func (m *StokLot) BeforeAppendModel(_ context.Context, query bun.Query) error {
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
