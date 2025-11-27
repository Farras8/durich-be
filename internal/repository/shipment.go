package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
	"errors"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *domain.Pengiriman) error
	GetByID(ctx context.Context, id string) (*domain.Pengiriman, error)
	GetList(ctx context.Context, tujuan, status string) ([]domain.Pengiriman, error)
	AddItem(ctx context.Context, detail *domain.PengirimanDetail) error
	RemoveItem(ctx context.Context, detailID string) error
	Finalize(ctx context.Context, id string) error
	GetDetailByID(ctx context.Context, id string) (*domain.PengirimanDetail, error)
}

type shipmentRepository struct {
	db *database.Database
}

func NewShipmentRepository(db *database.Database) ShipmentRepository {
	return &shipmentRepository{db: db}
}

func (r *shipmentRepository) Create(ctx context.Context, shipment *domain.Pengiriman) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(shipment).Exec(ctx)
	return err
}

func (r *shipmentRepository) GetByID(ctx context.Context, id string) (*domain.Pengiriman, error) {
	shipment := new(domain.Pengiriman)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(shipment).
		Relation("Details").
		Relation("Details.Lot").
		Relation("Details.Lot.JenisDurianDetail").
		Relation("Creator").
		Where("p.id = ?", id).
		Where("p.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

func (r *shipmentRepository) GetList(ctx context.Context, tujuan, status string) ([]domain.Pengiriman, error) {
	var shipments []domain.Pengiriman
	query := r.db.InitQuery(ctx).NewSelect().
		Model(&shipments).
		Relation("Details").
		Where("p.deleted_at IS NULL")

	if tujuan != "" {
		query = query.Where("p.tujuan ILIKE ?", "%"+tujuan+"%")
	}
	if status != "" {
		query = query.Where("p.status = ?", status)
	}

	query = query.Order("p.created_at DESC")

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return shipments, nil
}

func (r *shipmentRepository) AddItem(ctx context.Context, detail *domain.PengirimanDetail) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check Lot validity and stock
	lot := new(domain.StokLot)
	err = tx.NewSelect().Model(lot).Where("id = ?", detail.LotSumberID).For("UPDATE").Scan(ctx)
	if err != nil {
		return err
	}

	if lot.Status != "READY" {
		return errors.New("lot status must be READY")
	}
	if lot.QtySisa < detail.QtyAmbil {
		return errors.New("insufficient lot quantity")
	}
	if lot.BeratSisa < detail.BeratAmbil {
		return errors.New("insufficient lot weight")
	}

	// Insert Detail
	_, err = tx.NewInsert().Model(detail).Exec(ctx)
	if err != nil {
		return err
	}

	// Update Lot Stock
	lot.QtySisa -= detail.QtyAmbil
	lot.BeratSisa -= detail.BeratAmbil
	
	_, err = tx.NewUpdate().
		Model(lot).
		Column("qty_sisa", "berat_sisa").
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *shipmentRepository) RemoveItem(ctx context.Context, detailID string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get Detail
	detail := new(domain.PengirimanDetail)
	err = tx.NewSelect().Model(detail).Where("id = ?", detailID).Scan(ctx)
	if err != nil {
		return err
	}

	// Restore Lot Stock
	lot := new(domain.StokLot)
	err = tx.NewSelect().Model(lot).Where("id = ?", detail.LotSumberID).For("UPDATE").Scan(ctx)
	if err != nil {
		return err // If lot deleted, this might fail, handle gracefully? Assuming soft delete only.
	}

	lot.QtySisa += detail.QtyAmbil
	lot.BeratSisa += detail.BeratAmbil

	_, err = tx.NewUpdate().
		Model(lot).
		Column("qty_sisa", "berat_sisa").
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	// Delete Detail
	_, err = tx.NewDelete().Model(detail).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *shipmentRepository) Finalize(ctx context.Context, id string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Shipment Status
	_, err = tx.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", "OTW").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	// Check if any lots are now empty (QtySisa == 0) and update their status to EMPTY
	// Retrieve all lots involved in this shipment
	var details []domain.PengirimanDetail
	err = tx.NewSelect().Model(&details).Where("pengiriman_id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	for _, d := range details {
		lot := new(domain.StokLot)
		err = tx.NewSelect().Model(lot).Where("id = ?", d.LotSumberID).Scan(ctx)
		if err != nil {
			return err
		}
		
		if lot.QtySisa <= 0 {
			_, err = tx.NewUpdate().
				Model((*domain.StokLot)(nil)).
				Set("status = ?", "EMPTY").
				Where("id = ?", d.LotSumberID).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *shipmentRepository) GetDetailByID(ctx context.Context, id string) (*domain.PengirimanDetail, error) {
	detail := new(domain.PengirimanDetail)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(detail).
		Relation("Lot").
		Where("pd.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return detail, nil
}
