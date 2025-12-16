package repository

import (
	"context"
	"durich-be/internal/domain"
	"fmt"

	"github.com/uptrace/bun"
)

type LotRepository interface {
	Create(ctx context.Context, db bun.IDB, lot *domain.StokLot) error
	GetByID(ctx context.Context, db bun.IDB, id string) (*domain.StokLot, error)
	GetByIDForUpdate(ctx context.Context, db bun.IDB, id string) (*domain.StokLot, error)
	GetList(ctx context.Context, db bun.IDB, status, jenisDurianID, kondisi, locationID, scope, createdAt string) ([]domain.StokLot, error)
	Update(ctx context.Context, db bun.IDB, lot *domain.StokLot) error
	AddBuah(ctx context.Context, db bun.IDB, buah *domain.BuahRaw) error
	RemoveItem(ctx context.Context, db bun.IDB, lotID, buahRawID string) error
	GetItemCount(ctx context.Context, db bun.IDB, lotID string) (int, error)
	GetBuahRawByID(ctx context.Context, db bun.IDB, id string) (*domain.BuahRaw, error)
	GetNextLotKode(ctx context.Context, db bun.IDB) (string, error)
	GetNextLotSequence(ctx context.Context, db bun.IDB, dateStr, jenisKode, grade string) (string, error)
	GetPohonByKode(ctx context.Context, db bun.IDB, kode string, blokID string) (*domain.Pohon, error)
	GetTotalWeight(ctx context.Context, db bun.IDB, lotID string) (float64, error)
}

type lotRepository struct {
	db *bun.DB
}

func NewLotRepository(db *bun.DB) LotRepository {
	return &lotRepository{db: db}
}

func (r *lotRepository) Create(ctx context.Context, db bun.IDB, lot *domain.StokLot) error {
	_, err := db.NewInsert().Model(lot).Exec(ctx)
	return err
}

func (r *lotRepository) GetByID(ctx context.Context, db bun.IDB, id string) (*domain.StokLot, error) {
	lot := new(domain.StokLot)
	err := db.NewSelect().
		Model(lot).
		Relation("JenisDurianDetail").
		ColumnExpr("stok_lot.*").
		ColumnExpr("(SELECT COUNT(*) FROM tb_buah_raw WHERE lot_id = stok_lot.id) AS current_qty").
		ColumnExpr("(SELECT COALESCE(SUM(berat), 0) FROM tb_buah_raw WHERE lot_id = stok_lot.id) AS current_berat").
		Where("stok_lot.id = ?", id).
		Where("stok_lot.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lot, nil
}

func (r *lotRepository) GetByIDForUpdate(ctx context.Context, db bun.IDB, id string) (*domain.StokLot, error) {
	lot := new(domain.StokLot)
	err := db.NewSelect().
		Model(lot).
		Relation("JenisDurianDetail").
		ColumnExpr("stok_lot.*").
		ColumnExpr("(SELECT COUNT(*) FROM tb_buah_raw WHERE lot_id = stok_lot.id) AS current_qty").
		ColumnExpr("(SELECT COALESCE(SUM(berat), 0) FROM tb_buah_raw WHERE lot_id = stok_lot.id) AS current_berat").
		Where("stok_lot.id = ?", id).
		Where("stok_lot.deleted_at IS NULL").
		For("UPDATE OF stok_lot").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lot, nil
}

func (r *lotRepository) GetList(ctx context.Context, db bun.IDB, status, jenisDurianID, kondisi, locationID, scope, createdAt string) ([]domain.StokLot, error) {
	var lots []domain.StokLot
	query := db.NewSelect().
		Model(&lots).
		Relation("JenisDurianDetail").
		Relation("Posisi").
		ColumnExpr("stok_lot.*").
		ColumnExpr("(SELECT COUNT(*) FROM tb_buah_raw WHERE lot_id = stok_lot.id) AS current_qty").
		ColumnExpr("(SELECT COALESCE(SUM(berat), 0) FROM tb_buah_raw WHERE lot_id = stok_lot.id) AS current_berat").
		Where("stok_lot.deleted_at IS NULL")

	if locationID != "" {
		query = query.Where("stok_lot.current_location_id = ?", locationID)
	} else if scope == "local" {
		query = query.Where("stok_lot.current_location_id IS NULL")
	}

	if status != "" {
		query = query.Where("stok_lot.status = ?", status)
	}
	if jenisDurianID != "" {
		query = query.Where("stok_lot.jenis_durian_id = ?", jenisDurianID)
	}
	if kondisi != "" {
		query = query.Where("stok_lot.kondisi_buah = ?", kondisi)
	}
	if createdAt != "" {
		query = query.Where("DATE(stok_lot.created_at) = ?", createdAt)
	}

	query = query.Order("stok_lot.created_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lots, nil
}

func (r *lotRepository) Update(ctx context.Context, db bun.IDB, lot *domain.StokLot) error {
	_, err := db.NewUpdate().
		Model(lot).
		WherePK().
		Exec(ctx)
	return err
}

func (r *lotRepository) AddBuah(ctx context.Context, db bun.IDB, buah *domain.BuahRaw) error {
	_, err := db.NewInsert().Model(buah).Exec(ctx)
	return err
}

func (r *lotRepository) GetPohonByKode(ctx context.Context, db bun.IDB, kode string, blokID string) (*domain.Pohon, error) {
	pohon := new(domain.Pohon)
	query := db.NewSelect().
		Model(pohon).
		Relation("Blok").
		Relation("Blok.Divisi").
		Relation("Blok.Divisi.Estate").
		Relation("Blok.Divisi.Estate.Company").
		Where("pohon.kode = ?", kode).
		Where("pohon.blok_id = ?", blokID)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return pohon, nil
}

func (r *lotRepository) RemoveItem(ctx context.Context, db bun.IDB, lotID, buahRawID string) error {
	_, err := db.NewDelete().
		Model((*domain.BuahRaw)(nil)).
		Where("id = ?", buahRawID).
		Where("lot_id = ?", lotID).
		Exec(ctx)
	return err
}

func (r *lotRepository) GetItemCount(ctx context.Context, db bun.IDB, lotID string) (int, error) {
	count, err := db.NewSelect().
		Model((*domain.BuahRaw)(nil)).
		Where("lot_id = ?", lotID).
		Count(ctx)
	return count, err
}

func (r *lotRepository) GetTotalWeight(ctx context.Context, db bun.IDB, lotID string) (float64, error) {
	var totalWeight float64
	err := db.NewSelect().
		Model((*domain.BuahRaw)(nil)).
		ColumnExpr("COALESCE(SUM(berat), 0)").
		Where("lot_id = ?", lotID).
		Scan(ctx, &totalWeight)
	return totalWeight, err
}

func (r *lotRepository) GetBuahRawByID(ctx context.Context, db bun.IDB, id string) (*domain.BuahRaw, error) {
	buah := new(domain.BuahRaw)
	err := db.NewSelect().
		Model(buah).
		Relation("PohonPanenDetail.Blok.Divisi.Estate.Company").
		Where("buah_raw.id = ?", id).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return buah, nil
}

func (r *lotRepository) GetNextLotSequence(ctx context.Context, db bun.IDB, dateStr, jenisKode, grade string) (string, error) {
	prefix := fmt.Sprintf("LOT-%s-%s-%s", jenisKode, grade, dateStr)

	var lot domain.StokLot
	err := db.NewSelect().
		Model(&lot).
		Column("kode").
		Where("kode LIKE ?", prefix+"-%").
		Order("kode DESC").
		Limit(1).
		For("UPDATE").
		Scan(ctx)

	nextSeq := 1
	if err == nil && lot.Kode != "" {
		var seq int
		_, err = fmt.Sscanf(lot.Kode, prefix+"-%d", &seq)
		if err == nil {
			nextSeq = seq + 1
		}
	}

	newKode := fmt.Sprintf("%s-%02d", prefix, nextSeq)
	return newKode, nil
}

func (r *lotRepository) GetNextLotKode(ctx context.Context, db bun.IDB) (string, error) {
	var lot domain.StokLot
	err := db.NewSelect().
		Model(&lot).
		Column("kode").
		Where("kode LIKE ?", "LOT-%").
		Order("kode DESC").
		Limit(1).
		For("UPDATE").
		Scan(ctx)

	nextSeq := 1
	if err == nil && lot.Kode != "" {
		var seq int
		_, err = fmt.Sscanf(lot.Kode, "LOT-%d", &seq)
		if err == nil {
			nextSeq = seq + 1
		}
	}

	newKode := fmt.Sprintf("LOT-%03d", nextSeq)

	// Check if exists
	var existing domain.StokLot
	err = db.NewSelect().
		Model(&existing).
		Where("kode = ?", newKode).
		Scan(ctx)
	if err == nil {
		nextSeq++
		newKode = fmt.Sprintf("LOT-%03d", nextSeq)
	}

	return newKode, nil
}
