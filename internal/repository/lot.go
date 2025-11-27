package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
)

type LotRepository interface {
	Create(ctx context.Context, lot *domain.StokLot) error
	GetByID(ctx context.Context, id string) (*domain.StokLot, error)
	GetList(ctx context.Context, status, jenisDurian, kondisi string) ([]domain.StokLot, error)
	Update(ctx context.Context, lot *domain.StokLot) error
	AddItems(ctx context.Context, lotID string, buahRawIDs []string) error
	RemoveItem(ctx context.Context, lotID, buahRawID string) error
	GetItemCount(ctx context.Context, lotID string) (int, error)
	GetBuahRawByID(ctx context.Context, id string) (*domain.BuahRaw, error)
	UpdateBuahRawSorted(ctx context.Context, id string, isSorted bool) error
}

type lotRepository struct {
	db *database.Database
}

func NewLotRepository(db *database.Database) LotRepository {
	return &lotRepository{db: db}
}

func (r *lotRepository) Create(ctx context.Context, lot *domain.StokLot) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(lot).Exec(ctx)
	return err
}

func (r *lotRepository) GetByID(ctx context.Context, id string) (*domain.StokLot, error) {
	lot := new(domain.StokLot)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(lot).
		Relation("JenisDurianDetail").
		Where("stok_lot.id = ?", id).
		Where("stok_lot.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lot, nil
}

func (r *lotRepository) GetList(ctx context.Context, status, jenisDurian, kondisi string) ([]domain.StokLot, error) {
	var lots []domain.StokLot
	query := r.db.InitQuery(ctx).NewSelect().
		Model(&lots).
		Relation("JenisDurianDetail").
		Where("stok_lot.deleted_at IS NULL")

	if status != "" {
		query = query.Where("stok_lot.status = ?", status)
	}
	if jenisDurian != "" {
		query = query.Where("stok_lot.jenis_durian = ?", jenisDurian)
	}
	if kondisi != "" {
		query = query.Where("stok_lot.kondisi_buah = ?", kondisi)
	}

	query = query.Order("stok_lot.created_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return lots, nil
}

func (r *lotRepository) Update(ctx context.Context, lot *domain.StokLot) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model(lot).
		WherePK().
		Exec(ctx)
	return err
}

func (r *lotRepository) AddItems(ctx context.Context, lotID string, buahRawIDs []string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, buahID := range buahRawIDs {
		detail := &domain.LotDetail{
			LotID:     lotID,
			BuahRawID: buahID,
		}
		_, err := tx.NewInsert().Model(detail).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewUpdate().
			Model((*domain.BuahRaw)(nil)).
			Set("is_sorted = ?", true).
			Where("id = ?", buahID).
			Exec(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *lotRepository) RemoveItem(ctx context.Context, lotID, buahRawID string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewDelete().
		Model((*domain.LotDetail)(nil)).
		Where("lot_id = ?", lotID).
		Where("buah_raw_id = ?", buahRawID).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = tx.NewUpdate().
		Model((*domain.BuahRaw)(nil)).
		Set("is_sorted = ?", false).
		Where("id = ?", buahRawID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *lotRepository) GetItemCount(ctx context.Context, lotID string) (int, error) {
	count, err := r.db.InitQuery(ctx).NewSelect().
		Model((*domain.LotDetail)(nil)).
		Where("lot_id = ?", lotID).
		Count(ctx)
	return count, err
}

func (r *lotRepository) GetBuahRawByID(ctx context.Context, id string) (*domain.BuahRaw, error) {
	buah := new(domain.BuahRaw)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(buah).
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Where("buah_raw.id = ?", id).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return buah, nil
}

func (r *lotRepository) UpdateBuahRawSorted(ctx context.Context, id string, isSorted bool) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.BuahRaw)(nil)).
		Set("is_sorted = ?", isSorted).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
