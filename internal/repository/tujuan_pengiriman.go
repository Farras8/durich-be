package repository

import (
	"context"
	"database/sql"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
)

type TujuanPengirimanRepository interface {
	Create(ctx context.Context, tujuan *domain.TujuanPengiriman) error
	GetAll(ctx context.Context, tipe string) ([]domain.TujuanPengiriman, error)
	GetByID(ctx context.Context, id string) (*domain.TujuanPengiriman, error)
	Update(ctx context.Context, id string, tujuan *domain.TujuanPengiriman) error
	Delete(ctx context.Context, id string) error
}

type tujuanPengirimanRepository struct {
	db *database.Database
}

func NewTujuanPengirimanRepository(db *database.Database) TujuanPengirimanRepository {
	return &tujuanPengirimanRepository{db: db}
}

func (r *tujuanPengirimanRepository) Create(ctx context.Context, tujuan *domain.TujuanPengiriman) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(tujuan).Exec(ctx)
	return err
}

func (r *tujuanPengirimanRepository) GetAll(ctx context.Context, tipe string) ([]domain.TujuanPengiriman, error) {
	var tujuans []domain.TujuanPengiriman
	query := r.db.InitQuery(ctx).NewSelect().Model(&tujuans).Where("deleted_at IS NULL")

	if tipe != "" {
		query = query.Where("tipe = ?", tipe)
	}

	err := query.Scan(ctx)
	return tujuans, err
}

func (r *tujuanPengirimanRepository) GetByID(ctx context.Context, id string) (*domain.TujuanPengiriman, error) {
	tujuan := &domain.TujuanPengiriman{}
	err := r.db.InitQuery(ctx).NewSelect().Model(tujuan).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return tujuan, err
}

func (r *tujuanPengirimanRepository) Update(ctx context.Context, id string, tujuan *domain.TujuanPengiriman) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(tujuan).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *tujuanPengirimanRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.TujuanPengiriman)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}
