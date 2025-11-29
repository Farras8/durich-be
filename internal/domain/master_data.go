package domain

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/uptrace/bun"
)

type Company struct {
	bun.BaseModel `bun:"table:company,alias:company"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string     `bun:",unique,notnull" json:"kode"`
	Nama      string     `bun:",notnull" json:"nama"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"," json:"deleted_at,omitempty"`
}

func (m *Company) BeforeAppendModel(_ context.Context, query bun.Query) error {
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

type Estate struct {
	bun.BaseModel `bun:"table:estate,alias:estate"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string     `bun:",unique,notnull" json:"kode"`
	Nama      string     `bun:",notnull" json:"nama"`
	CompanyID string     `bun:",notnull" json:"company_id"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"," json:"deleted_at,omitempty"`
	Company   *Company   `bun:"rel:belongs-to,join:company_id=id" json:"company,omitempty"`
}

func (m *Estate) BeforeAppendModel(_ context.Context, query bun.Query) error {
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

type Divisi struct {
	bun.BaseModel `bun:"table:divisi,alias:divisi"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string     `bun:",unique,notnull" json:"kode"`
	Nama      string     `bun:",notnull" json:"nama"`
	EstateID  string     `bun:",notnull" json:"estate_id"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"," json:"deleted_at,omitempty"`
	Estate    *Estate    `bun:"rel:belongs-to,join:estate_id=id" json:"estate,omitempty"`
}

func (m *Divisi) BeforeAppendModel(_ context.Context, query bun.Query) error {
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

type Blok struct {
	bun.BaseModel `bun:"table:blok,alias:blok"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string     `bun:",notnull" json:"kode"`
	NamaBlok  string     `bun:",notnull" json:"nama_blok"`
	DivisiID  string     `bun:",notnull" json:"divisi_id"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"," json:"deleted_at,omitempty"`
	Divisi    *Divisi    `bun:"rel:belongs-to,join:divisi_id=id" json:"divisi,omitempty"`
}

func (m *Blok) BeforeAppendModel(_ context.Context, query bun.Query) error {
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

type JenisDurian struct {
	bun.BaseModel `bun:"table:jenis_durian,alias:jenis_durian"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string     `bun:"," json:"kode"`
	NamaJenis string     `bun:",notnull" json:"nama_jenis"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"," json:"deleted_at,omitempty"`
}

func (m *JenisDurian) BeforeAppendModel(_ context.Context, query bun.Query) error {
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

type Pohon struct {
	bun.BaseModel `bun:"table:pohon,alias:pohon"`

	ID        string     `bun:",pk" json:"id"`
	Kode      string     `bun:",unique,notnull" json:"kode"`
	Nama      string     `bun:",notnull" json:"nama"`
	BlokID    *string    `bun:"," json:"blok_id,omitempty"`
	CreatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"," json:"deleted_at,omitempty"`
	Blok      *Blok      `bun:"rel:belongs-to,join:blok_id=id" json:"blok,omitempty"`
}

func (m *Pohon) BeforeAppendModel(_ context.Context, query bun.Query) error {
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
