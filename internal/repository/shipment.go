package repository

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type ShipmentReceiveItem struct {
	Berat float64
	Qty   int
}

type ShipmentRepository interface {
	Create(ctx context.Context, db bun.IDB, shipment *domain.Pengiriman) error
	GetByID(ctx context.Context, db bun.IDB, id string) (*domain.Pengiriman, error)
	GetList(ctx context.Context, db bun.IDB, tujuan, status, locationID, listType, tujuanType string, page, limit int) ([]domain.Pengiriman, int64, error)
	AddItem(ctx context.Context, db bun.IDB, detail *domain.PengirimanDetail, locationID string) error
	RemoveItem(ctx context.Context, db bun.IDB, shipmentID, detailID string) error
	UpdateStatus(ctx context.Context, db bun.IDB, id, status, notes, userID string) error
	UpdateShipmentToSending(ctx context.Context, db bun.IDB, id string) error
	UpdateLotsToShipped(ctx context.Context, db bun.IDB, lotIDs []string) error
	GetDetailByID(ctx context.Context, db bun.IDB, id string) (*domain.PengirimanDetail, error)
	GetNextShipmentKode(ctx context.Context, db bun.IDB) (string, error)
	UpdateShipmentToReceived(ctx context.Context, db bun.IDB, id string, receivedDate time.Time) error
	UpdateLotsAfterReceive(ctx context.Context, db bun.IDB, lotID, tujuanID string, berat float64, qty int, receivedDate time.Time) error
	GetDetailsByShipmentID(ctx context.Context, db bun.IDB, shipmentID string) ([]domain.PengirimanDetail, error)
}

type shipmentRepository struct{}

func NewShipmentRepository() ShipmentRepository {
	return &shipmentRepository{}
}

func (r *shipmentRepository) Create(ctx context.Context, db bun.IDB, shipment *domain.Pengiriman) error {
	_, err := db.NewInsert().Model(shipment).Exec(ctx)
	return err
}

func (r *shipmentRepository) GetByID(ctx context.Context, db bun.IDB, id string) (*domain.Pengiriman, error) {
	shipment := new(domain.Pengiriman)
	err := db.NewSelect().
		Model(shipment).
		Relation("Details").
		Relation("Details.Lot").
		Relation("Details.Lot.JenisDurianDetail").
		Relation("Creator").
		Relation("TujuanDetail").
		Where("p.id = ?", id).
		Where("p.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

func (r *shipmentRepository) GetList(ctx context.Context, db bun.IDB, tujuan, status, locationID, listType, tujuanType string, page, limit int) ([]domain.Pengiriman, int64, error) {
	var shipments []domain.Pengiriman

	query := db.NewSelect().
		Model(&shipments).
		Relation("Details").
		Relation("Creator").
		Relation("TujuanDetail").
		Where("p.deleted_at IS NULL")

	if tujuanType != "" {
		query = query.Where("tujuan_detail.tipe = ?", tujuanType)
	}

	if locationID != "" {
		if listType == "incoming" {
			query = query.Where("p.tujuan_id = ?", locationID)
		} else if listType == "outgoing" {
			query = query.Where("creator.current_location_id = ?", locationID)
		} else {
			query = query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where("p.tujuan_id = ?", locationID).
					WhereOr("creator.current_location_id = ?", locationID)
			})
		}
	}

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

// AddItem - Pure data access, terima bun.IDB
func (r *shipmentRepository) AddItem(ctx context.Context, db bun.IDB, detail *domain.PengirimanDetail, locationID string) error {
	var shipmentStatus string
	var shipmentCreatorLocation *string
	err := db.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		ColumnExpr("creator.current_location_id").
		Join("JOIN users AS creator ON creator.id = p.created_by").
		Where("p.id = ?", detail.PengirimanID).
		Where("p.deleted_at IS NULL").
		For("UPDATE OF p").
		Scan(ctx, &shipmentStatus, &shipmentCreatorLocation)

	if err != nil {
		return errors.New("shipment not found")
	}
	if shipmentStatus != constants.ShipmentStatusDraft {
		return errors.New("shipment must be DRAFT to add items")
	}

	lot := new(domain.StokLot)
	query := db.NewSelect().
		Model(lot).
		Where("id = ?", detail.LotSumberID).
		Where("status = ?", constants.LotStatusReady).
		For("UPDATE")

	if locationID == "" {
		query = query.Where("current_location_id IS NULL")
	} else {
		query = query.Where("current_location_id = ?", locationID)
	}

	err = query.Scan(ctx)
	if err != nil {
		return errors.New("lot not found, not in READY status, or belongs to another location")
	}

	if lot.QtySisa <= 0 || lot.BeratSisa <= 0 {
		return errors.New("lot has insufficient stock")
	}

	detail.QtyAmbil = lot.QtySisa
	detail.BeratAmbil = lot.BeratSisa

	_, err = db.NewInsert().Model(detail).Exec(ctx)
	if err != nil {
		return err
	}

	lot.Status = constants.LotStatusBooked
	lot.QtySisa = 0
	lot.BeratSisa = 0

	_, err = db.NewUpdate().
		Model(lot).
		Column("status", "qty_sisa", "berat_sisa").
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *shipmentRepository) RemoveItem(ctx context.Context, db bun.IDB, shipmentID, detailID string) error {
	var shipmentStatus string
	err := db.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		Where("id = ?", shipmentID).
		Where("deleted_at IS NULL").
		For("UPDATE").
		Scan(ctx, &shipmentStatus)
	if err != nil {
		return errors.New("shipment not found")
	}
	if shipmentStatus != constants.ShipmentStatusDraft {
		return errors.New("shipment must be DRAFT to remove items")
	}
	detail := new(domain.PengirimanDetail)
	err = db.NewSelect().Model(detail).Where("id = ?", detailID).Scan(ctx)
	if err != nil {
		return errors.New("shipment detail not found")
	}

	if detail.PengirimanID != shipmentID {
		return errors.New("detail does not belong to this shipment")
	}

	lot := new(domain.StokLot)
	err = db.NewSelect().Model(lot).Where("id = ?", detail.LotSumberID).For("UPDATE").Scan(ctx)
	if err != nil {
		return errors.New("lot not found")
	}

	if lot.Status != constants.LotStatusBooked {
		return errors.New("lot is not in BOOKED status, cannot restore")
	}
	lot.QtySisa = detail.QtyAmbil
	lot.BeratSisa = detail.BeratAmbil
	lot.Status = constants.LotStatusReady

	_, err = db.NewUpdate().
		Model(lot).
		Column("qty_sisa", "berat_sisa", "status").
		WherePK().
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewDelete().Model(detail).WherePK().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatus - Simple update, transaction di-handle di service
func (r *shipmentRepository) UpdateStatus(ctx context.Context, db bun.IDB, id, status, notes, userID string) error {
	var currentStatus string
	err := db.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		For("UPDATE").
		Scan(ctx, &currentStatus)
	if err != nil {
		return errors.New("shipment not found")
	}

	// Update status
	_, err = db.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateShipmentToSending - Dipisah untuk composability
func (r *shipmentRepository) UpdateShipmentToSending(ctx context.Context, db bun.IDB, id string) error {
	var currentStatus string
	err := db.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		For("UPDATE").
		Scan(ctx, &currentStatus)
	if err != nil {
		return errors.New("shipment not found")
	}
	if currentStatus != constants.ShipmentStatusDraft {
		return errors.New("shipment must be DRAFT to finalize")
	}

	_, err = db.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", constants.ShipmentStatusSending).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// UpdateLotsToShipped - Dipisah untuk composability
func (r *shipmentRepository) UpdateLotsToShipped(ctx context.Context, db bun.IDB, lotIDs []string) error {
	if len(lotIDs) == 0 {
		return nil
	}

	_, err := db.NewUpdate().
		Model((*domain.StokLot)(nil)).
		Set("status = ?", constants.LotStatusShipped).
		Set("updated_at = ?", time.Now()).
		Where("id IN (?)", bun.In(lotIDs)).
		Where("qty_sisa = 0").
		Where("berat_sisa = 0").
		Exec(ctx)
	return err
}

func (r *shipmentRepository) GetDetailByID(ctx context.Context, db bun.IDB, id string) (*domain.PengirimanDetail, error) {
	detail := new(domain.PengirimanDetail)
	err := db.NewSelect().
		Model(detail).
		Relation("Lot").
		Where("pd.id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

// GetNextShipmentKode - Menggunakan FOR UPDATE untuk prevent race condition
func (r *shipmentRepository) GetNextShipmentKode(ctx context.Context, db bun.IDB) (string, error) {
	dateStr := time.Now().Format("060102") // YYMMDD
	prefix := fmt.Sprintf("SHP-%s", dateStr)

	var lastCode string
	err := db.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("kode").
		Where("kode LIKE ?", prefix+"-%").
		Order("kode DESC").
		Limit(1).
		For("UPDATE SKIP LOCKED").
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

// GetDetailsByShipmentID - Helper untuk get details dengan lock
func (r *shipmentRepository) GetDetailsByShipmentID(ctx context.Context, db bun.IDB, shipmentID string) ([]domain.PengirimanDetail, error) {
	var details []domain.PengirimanDetail
	err := db.NewSelect().
		Model(&details).
		Where("pengiriman_id = ?", shipmentID).
		For("UPDATE").
		Scan(ctx)
	return details, err
}

// UpdateShipmentToReceived - Dipisah untuk composability
func (r *shipmentRepository) UpdateShipmentToReceived(ctx context.Context, db bun.IDB, id string, receivedDate time.Time) error {
	var currentStatus string
	err := db.NewSelect().
		Model((*domain.Pengiriman)(nil)).
		Column("status").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		For("UPDATE").
		Scan(ctx, &currentStatus)
	if err != nil {
		return errors.New("shipment not found")
	}

	if currentStatus != constants.ShipmentStatusSending {
		return fmt.Errorf("invalid shipment status: %s, expected: %s", currentStatus, constants.ShipmentStatusSending)
	}

	_, err = db.NewUpdate().
		Model((*domain.Pengiriman)(nil)).
		Set("status = ?", constants.ShipmentStatusReceived).
		Set("received_at = ?", receivedDate).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

// UpdateLotsAfterReceive - Update individual lot
func (r *shipmentRepository) UpdateLotsAfterReceive(ctx context.Context, db bun.IDB, lotID, tujuanID string, berat float64, qty int, receivedDate time.Time) error {
	var existingLot domain.StokLot
	err := db.NewSelect().
		Model(&existingLot).
		Where("id = ?", lotID).
		For("UPDATE").
		Scan(ctx)
	if err != nil {
		return fmt.Errorf("lot %s not found", lotID)
	}

	if berat < 0 || qty < 0 {
		return fmt.Errorf("received quantity/weight cannot be negative for lot %s", lotID)
	}

	_, err = db.NewUpdate().
		Model((*domain.StokLot)(nil)).
		Set("current_location_id = ?", tujuanID).
		Set("berat_sisa = ?", berat).
		Set("qty_sisa = ?", qty).
		Set("status = ?", constants.LotStatusReady).
		Set("arrived_at = ?", receivedDate).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", lotID).
		Exec(ctx)
	return err
}
