package repository

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *domain.Pengiriman) error
	GetByID(ctx context.Context, id string) (*domain.Pengiriman, error)
	GetList(ctx context.Context, tujuan, status string, page, limit int) ([]domain.Pengiriman, int64, error)
	AddItem(ctx context.Context, detail *domain.PengirimanDetail) error
	RemoveItem(ctx context.Context, shipmentID, detailID string) error
	UpdateStatus(ctx context.Context, id, status, notes, userID string) error
	Finalize(ctx context.Context, id string) error
	GetDetailByID(ctx context.Context, id string) (*domain.PengirimanDetail, error)
	GetNextShipmentKode(ctx context.Context) (string, error)
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

func (r *shipmentRepository) GetList(ctx context.Context, tujuan, status string, page, limit int) ([]domain.Pengiriman, int64, error) {
	var shipments []domain.Pengiriman

	query := r.db.InitQuery(ctx).NewSelect().
		Model(&shipments).
		Relation("Details").
		Relation("Creator").
		Where("p.deleted_at IS NULL")

	if tujuan != "" {
		query = query.Where("p.tujuan ILIKE ?", "%"+tujuan+"%")
	}
	if status != "" {
		query = query.Where("p.status = ?", status)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query = query.Order("p.created_at DESC").Limit(limit).Offset(offset)

	err = query.Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	return shipments, int64(total), nil
}

func (r *shipmentRepository) AddItem(ctx context.Context, detail *domain.PengirimanDetail) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var shipmentStatus string
	err = tx.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		Where("id = ?", detail.PengirimanID).
		Where("deleted_at IS NULL").
		Scan(ctx, &shipmentStatus)
	if err != nil {
		return errors.New("shipment not found")
	}
	if shipmentStatus != constants.ShipmentStatusDraft {
		return errors.New("shipment must be DRAFT to add items")
	}

	lot := new(domain.StokLot)
	err = tx.NewSelect().
		Model(lot).
		Where("id = ?", detail.LotSumberID).
		Where("status = ?", constants.LotStatusReady).
		Where("qty_sisa >= ?", detail.QtyAmbil).
		Where("berat_sisa >= ?", detail.BeratAmbil).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		return errors.New("lot not found, not in READY status, or insufficient quantity/weight")
	}

	_, err = tx.NewInsert().Model(detail).Exec(ctx)
	if err != nil {
		return err
	}

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

func (r *shipmentRepository) RemoveItem(ctx context.Context, shipmentID, detailID string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var shipmentStatus string
	err = tx.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		Where("id = ?", shipmentID).
		Where("deleted_at IS NULL").
		Scan(ctx, &shipmentStatus)
	if err != nil {
		return errors.New("shipment not found")
	}
	if shipmentStatus != constants.ShipmentStatusDraft {
		return errors.New("shipment must be DRAFT to remove items")
	}

	detail := new(domain.PengirimanDetail)
	err = tx.NewSelect().Model(detail).Where("id = ?", detailID).Scan(ctx)
	if err != nil {
		return err
	}

	if detail.PengirimanID != shipmentID {
		return errors.New("detail does not belong to this shipment")
	}

	lot := new(domain.StokLot)
	err = tx.NewSelect().Model(lot).Where("id = ?", detail.LotSumberID).For("UPDATE").Scan(ctx)
	if err != nil {
		return err
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

	_, err = tx.NewDelete().Model(detail).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *shipmentRepository) UpdateStatus(ctx context.Context, id, status, notes, userID string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", status).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *shipmentRepository) Finalize(ctx context.Context, id string) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", constants.ShipmentStatusSending).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	var details []domain.PengirimanDetail
	err = tx.NewSelect().Model(&details).Where("pengiriman_id = ?", id).Scan(ctx)
	if err != nil {
		return err
	}

	// Check for empty lots
	var emptyLotIDs []string
	for _, d := range details {
		lot := new(domain.StokLot)
		err = tx.NewSelect().Model(lot).Where("id = ?", d.LotSumberID).Scan(ctx)
		if err != nil {
			return err
		}

		if lot.QtySisa <= 0 {
			emptyLotIDs = append(emptyLotIDs, d.LotSumberID)
		}
	}

	if len(emptyLotIDs) > 0 {
		_, err = tx.NewUpdate().
			Model((*domain.StokLot)(nil)).
			Set("status = ?", constants.LotStatusEmpty).
			Where("id IN (?)", bun.In(emptyLotIDs)).
			Exec(ctx)
		if err != nil {
			return err
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

func (r *shipmentRepository) GetNextShipmentKode(ctx context.Context) (string, error) {
	dateStr := time.Now().Format("060102") // YYMMDD
	prefix := fmt.Sprintf("SHP-%s", dateStr)

	var lastCode string
	err := r.db.InitQuery(ctx).NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("kode").
		Where("kode LIKE ?", prefix+"-%").
		Order("kode DESC").
		Limit(1).
		Scan(ctx, &lastCode)

	seq := 1
	if err == nil && lastCode != "" {
		var lastSeq int
		_, err := fmt.Sscanf(lastCode, prefix+"-%d", &lastSeq)
		if err == nil {
			seq = lastSeq + 1
		}
	}

	return fmt.Sprintf("%s-%03d", prefix, seq), nil
}
