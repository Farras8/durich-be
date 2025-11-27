package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
)

type SalesRepository interface {
	Create(ctx context.Context, sales *domain.Penjualan) error
	GetList(ctx context.Context, startDate, endDate, tipeJual string) ([]domain.Penjualan, error)
	GetByID(ctx context.Context, id string) (*domain.Penjualan, error)
	Update(ctx context.Context, sales *domain.Penjualan) error
	Delete(ctx context.Context, id string) error
	GetPengirimanByID(ctx context.Context, id string) (*domain.Pengiriman, error)
	UpdatePengirimanStatus(ctx context.Context, id, status string) error
	CheckSalesExistByShipmentID(ctx context.Context, shipmentID string) (bool, error)
}

type salesRepository struct {
	db *database.Database
}

func NewSalesRepository(db *database.Database) SalesRepository {
	return &salesRepository{db: db}
}

func (r *salesRepository) Create(ctx context.Context, sales *domain.Penjualan) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewInsert().Model(sales).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = tx.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", "COMPLETED").
		Where("id = ?", sales.PengirimanID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *salesRepository) GetList(ctx context.Context, startDate, endDate, tipeJual string) ([]domain.Penjualan, error) {
	var sales []domain.Penjualan
	query := r.db.InitQuery(ctx).NewSelect().
		Model(&sales).
		Relation("Pengiriman").
		Where("penjualan.deleted_at IS NULL")

	if startDate != "" {
		query = query.Where("penjualan.created_at >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("penjualan.created_at <= ?", endDate)
	}
	if tipeJual != "" {
		query = query.Where("penjualan.tipe_jual = ?", tipeJual)
	}

	query = query.Order("penjualan.created_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *salesRepository) GetByID(ctx context.Context, id string) (*domain.Penjualan, error) {
	sales := new(domain.Penjualan)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(sales).
		Relation("Pengiriman").
		Relation("Pengiriman.Details").
		Relation("Pengiriman.Details.Lot").
		Relation("Pengiriman.Details.Lot.JenisDurianDetail").
		Where("penjualan.id = ?", id).
		Where("penjualan.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return sales, nil
}

func (r *salesRepository) Update(ctx context.Context, sales *domain.Penjualan) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model(sales).
		Column("harga_total", "tipe_jual", "updated_at").
		WherePK().
		Exec(ctx)
	return err
}

func (r *salesRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get sales to know shipment ID
	sales := new(domain.Penjualan)
	err = tx.NewSelect().Model(sales).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	// Soft Delete Sales
	_, err = tx.NewUpdate().
		Model((*domain.Penjualan)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Restore Shipment Status to OTW
	_, err = tx.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", "OTW").
		Where("id = ?", sales.PengirimanID).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *salesRepository) GetPengirimanByID(ctx context.Context, id string) (*domain.Pengiriman, error) {
	shipment := new(domain.Pengiriman)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(shipment).
		Relation("Details").
		Where("p.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

func (r *salesRepository) UpdatePengirimanStatus(ctx context.Context, id, status string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", status).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *salesRepository) CheckSalesExistByShipmentID(ctx context.Context, shipmentID string) (bool, error) {
	count, err := r.db.InitQuery(ctx).NewSelect().
		Model((*domain.Penjualan)(nil)).
		Where("pengiriman_id = ?", shipmentID).
		Where("deleted_at IS NULL").
		Count(ctx)
	return count > 0, err
}
