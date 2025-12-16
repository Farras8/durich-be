package services

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"durich-be/pkg/database"
	"durich-be/pkg/errors"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type ShipmentService interface {
	Create(ctx context.Context, req requests.ShipmentCreateRequest, userID string) (*response.ShipmentResponse, error)
	GetList(ctx context.Context, tujuan, status, locationID, listType, tujuanType string, page, limit int) ([]response.ShipmentResponse, int64, error)
	GetByID(ctx context.Context, id string) (*response.ShipmentDetailResponse, error)
	AddItem(ctx context.Context, shipmentID string, req requests.ShipmentAddItemRequest, locationID string) error
	RemoveItem(ctx context.Context, shipmentID string, detailID string) error
	UpdateStatus(ctx context.Context, shipmentID string, req requests.ShipmentUpdateStatusRequest, userID string) error
	Finalize(ctx context.Context, id string) error
	Receive(ctx context.Context, id string, req requests.ShipmentReceiveRequest) error
}

type shipmentService struct {
	db         *database.Database
	repo       repository.ShipmentRepository
	tujuanRepo repository.TujuanPengirimanRepository
}

func NewShipmentService(
	db *database.Database,
	repo repository.ShipmentRepository,
	tujuanRepo repository.TujuanPengirimanRepository,
) ShipmentService {
	return &shipmentService{
		db:         db,
		repo:       repo,
		tujuanRepo: tujuanRepo,
	}
}

func (s *shipmentService) Create(ctx context.Context, req requests.ShipmentCreateRequest, userID string) (*response.ShipmentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tujuanDetail, err := s.tujuanRepo.GetByID(ctx, req.TujuanID)
	if err != nil {
		return nil, errors.ValidationError("invalid tujuan_id")
	}
	if tujuanDetail == nil {
		return nil, errors.ValidationError("tujuan pengiriman not found")
	}

	var createdShipment *domain.Pengiriman

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		kode, err := s.repo.GetNextShipmentKode(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to generate shipment code: %w", err)
		}

		tglKirim := req.TglKirim
		if tglKirim.IsZero() {
			tglKirim = time.Now()
		}

		shipment := &domain.Pengiriman{
			Kode:      kode,
			Tujuan:    tujuanDetail.Nama,
			TujuanID:  req.TujuanID,
			TglKirim:  tglKirim,
			Status:    constants.ShipmentStatusDraft,
			CreatedBy: userID,
		}

		err = s.repo.Create(ctx, tx, shipment)
		if err != nil {
			return fmt.Errorf("failed to create shipment: %w", err)
		}

		createdShipment = shipment
		return nil
	})

	if err != nil {
		return nil, err
	}

	resp := response.NewShipmentResponse(createdShipment)
	return &resp, nil
}

func (s *shipmentService) GetList(ctx context.Context, tujuan, status, locationID, listType, tujuanType string, page, limit int) ([]response.ShipmentResponse, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	shipments, total, err := s.repo.GetList(ctx, s.db.DB, tujuan, status, locationID, listType, tujuanType, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get shipment list: %w", err)
	}

	var resps []response.ShipmentResponse
	for _, p := range shipments {
		resps = append(resps, response.NewShipmentResponse(&p))
	}
	return resps, total, nil
}

func (s *shipmentService) GetByID(ctx context.Context, id string) (*response.ShipmentDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if id == "" {
		return nil, errors.ValidationError("shipment id is required")
	}

	p, err := s.repo.GetByID(ctx, s.db.DB, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}
	if p == nil {
		return nil, errors.NotFoundError("shipment not found")
	}

	header := response.NewShipmentResponse(p)
	items := make([]response.ShipmentItemResponse, 0, len(p.Details))

	for _, d := range p.Details {
		item := response.ShipmentItemResponse{
			ID:         d.ID,
			LotID:      d.LotSumberID,
			QtyAmbil:   d.QtyAmbil,
			BeratAmbil: d.BeratAmbil,
		}
		if d.Lot != nil {
			item.KodeLot = d.Lot.Kode
			item.Grade = d.Lot.KondisiBuah
			if d.Lot.JenisDurianDetail != nil {
				item.JenisDurian = d.Lot.JenisDurianDetail.NamaJenis
			}
		}
		items = append(items, item)
	}

	return &response.ShipmentDetailResponse{
		Header: header,
		Items:  items,
	}, nil
}

func (s *shipmentService) AddItem(ctx context.Context, shipmentID string, req requests.ShipmentAddItemRequest, locationID string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if shipmentID == "" {
		return errors.ValidationError("shipment id is required")
	}
	if req.LotID == "" {
		return errors.ValidationError("lot id is required")
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		shipment, err := s.repo.GetByID(ctx, tx, shipmentID)
		if err != nil {
			return fmt.Errorf("failed to get shipment: %w", err)
		}
		if shipment == nil {
			return errors.NotFoundError("shipment not found")
		}
		if shipment.Status != constants.ShipmentStatusDraft {
			return errors.ValidationError("shipment must be DRAFT to add items")
		}

		for _, d := range shipment.Details {
			if d.LotSumberID == req.LotID {
				return errors.ValidationError("lot already added to this shipment")
			}
		}

		detail := &domain.PengirimanDetail{
			PengirimanID: shipmentID,
			LotSumberID:  req.LotID,
		}

		err = s.repo.AddItem(ctx, tx, detail, locationID)
		if err != nil {
			return fmt.Errorf("failed to add item: %w", err)
		}

		return nil
	})
}

func (s *shipmentService) RemoveItem(ctx context.Context, shipmentID string, detailID string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if shipmentID == "" {
		return errors.ValidationError("shipment id is required")
	}
	if detailID == "" {
		return errors.ValidationError("detail id is required")
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		err := s.repo.RemoveItem(ctx, tx, shipmentID, detailID)
		if err != nil {
			return fmt.Errorf("failed to remove item: %w", err)
		}
		return nil
	})
}

func (s *shipmentService) UpdateStatus(ctx context.Context, shipmentID string, req requests.ShipmentUpdateStatusRequest, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if shipmentID == "" {
		return errors.ValidationError("shipment id is required")
	}
	if req.Status == "" {
		return errors.ValidationError("status is required")
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		shipment, err := s.repo.GetByID(ctx, tx, shipmentID)
		if err != nil {
			return fmt.Errorf("failed to get shipment: %w", err)
		}
		if shipment == nil {
			return errors.NotFoundError("shipment not found")
		}

		currentStatus := shipment.Status
		newStatus := req.Status

		isValidTransition := s.isValidStatusTransition(currentStatus, newStatus)
		if !isValidTransition {
			return errors.ValidationError(fmt.Sprintf("invalid status transition from %s to %s", currentStatus, newStatus))
		}

		if newStatus == constants.ShipmentStatusReceived {
			if len(shipment.Details) == 0 {
				return errors.ValidationError("cannot mark as received: shipment has no items")
			}
		}

		err = s.repo.UpdateStatus(ctx, tx, shipmentID, newStatus, req.Notes, userID)
		if err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}

		return nil
	})
}

func (s *shipmentService) isValidStatusTransition(current, target string) bool {
	validTransitions := map[string][]string{
		constants.ShipmentStatusDraft:    {constants.ShipmentStatusSending},
		constants.ShipmentStatusSending:  {constants.ShipmentStatusReceived},
		constants.ShipmentStatusReceived: {constants.ShipmentStatusCompleted},
	}

	allowedTargets, exists := validTransitions[current]
	if !exists {
		return false
	}

	for _, allowed := range allowedTargets {
		if allowed == target {
			return true
		}
	}
	return false
}

func (s *shipmentService) Finalize(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if id == "" {
		return errors.ValidationError("shipment id is required")
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		shipment, err := s.repo.GetByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get shipment: %w", err)
		}
		if shipment == nil {
			return errors.NotFoundError("shipment not found")
		}

		if shipment.Status != constants.ShipmentStatusDraft {
			return errors.ValidationError("shipment must be DRAFT to finalize")
		}

		if len(shipment.Details) == 0 {
			return errors.ValidationError("cannot finalize empty shipment")
		}

		tujuan, err := s.tujuanRepo.GetByID(ctx, shipment.TujuanID)
		if err != nil {
			return fmt.Errorf("failed to get tujuan: %w", err)
		}
		if tujuan == nil {
			return errors.ValidationError("tujuan pengiriman not found")
		}

		details, err := s.repo.GetDetailsByShipmentID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get details: %w", err)
		}

		var lotIDs []string
		for _, d := range details {
			lotIDs = append(lotIDs, d.LotSumberID)
		}

		var bookedCount int
		bookedCount, err = tx.NewSelect().
			Model((*domain.StokLot)(nil)).
			Where("id IN (?)", bun.In(lotIDs)).
			Where("status = ?", constants.LotStatusBooked).
			Count(ctx)
		if err != nil {
			return err
		}
		if bookedCount != len(details) {
			return errors.ValidationError("all lots must be in BOOKED status")
		}

		err = s.repo.UpdateShipmentToSending(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to update shipment status: %w", err)
		}

		err = s.repo.UpdateLotsToShipped(ctx, tx, lotIDs)
		if err != nil {
			return fmt.Errorf("failed to update lot status: %w", err)
		}

		// ✅ TODO: Add audit log here if needed
		// s.auditRepo.Log(ctx, tx, ...)

		return nil
	})
}

func (s *shipmentService) Receive(ctx context.Context, id string, req requests.ShipmentReceiveRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if id == "" {
		return errors.ValidationError("shipment id is required")
	}
	if len(req.Details) == 0 {
		return errors.ValidationError("received items cannot be empty")
	}
	if req.ReceivedDate.IsZero() {
		return errors.ValidationError("received date is required")
	}
	if req.ReceivedDate.After(time.Now()) {
		return errors.ValidationError("received date cannot be in the future")
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		shipment, err := s.repo.GetByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get shipment: %w", err)
		}
		if shipment == nil {
			return errors.NotFoundError("shipment not found")
		}

		if shipment.Status != constants.ShipmentStatusSending {
			return errors.ValidationError(fmt.Sprintf("shipment must be in SENDING status to receive, current status: %s", shipment.Status))
		}

		tujuan, err := s.tujuanRepo.GetByID(ctx, shipment.TujuanID)
		if err != nil {
			return fmt.Errorf("failed to get tujuan: %w", err)
		}
		if tujuan == nil {
			return errors.ValidationError("tujuan pengiriman not found")
		}
		if tujuan.Tipe != "internal" {
			return errors.ValidationError("only internal shipments can be received via this endpoint")
		}

		existingLots := make(map[string]domain.PengirimanDetail)
		for _, detail := range shipment.Details {
			existingLots[detail.LotSumberID] = detail
		}

		updates := make(map[string]repository.ShipmentReceiveItem)

		for _, item := range req.Details {
			detail, exists := existingLots[item.LotID]
			if !exists {
				return errors.ValidationError(fmt.Sprintf("lot id %s is not part of this shipment", item.LotID))
			}
			if item.BeratDiterima < 0 {
				return errors.ValidationError(fmt.Sprintf("received weight cannot be negative for lot %s", item.LotID))
			}
			if item.BeratDiterima > detail.BeratAmbil*1.1 {
				return errors.ValidationError(fmt.Sprintf("received weight exceeds sent weight by more than 10%% for lot %s", item.LotID))
			}

			finalQty := detail.QtyAmbil
			if item.QtyDiterima != nil {
				if *item.QtyDiterima < 0 {
					return errors.ValidationError(fmt.Sprintf("received quantity cannot be negative for lot %s", item.LotID))
				}
				if *item.QtyDiterima > detail.QtyAmbil {
					return errors.ValidationError(fmt.Sprintf("received quantity exceeds sent quantity for lot %s", item.LotID))
				}
				finalQty = *item.QtyDiterima
			}

			updates[item.LotID] = repository.ShipmentReceiveItem{
				Berat: item.BeratDiterima,
				Qty:   finalQty,
			}
		}

		if len(updates) != len(shipment.Details) {
			return errors.ValidationError("all items must be received")
		}

		err = s.repo.UpdateShipmentToReceived(ctx, tx, id, req.ReceivedDate)
		if err != nil {
			return fmt.Errorf("failed to update shipment status: %w", err)
		}

		for lotID, item := range updates {
			err := s.repo.UpdateLotsAfterReceive(ctx, tx, lotID, shipment.TujuanID, item.Berat, item.Qty, req.ReceivedDate)
			if err != nil {
				return fmt.Errorf("failed to update lot %s: %w", lotID, err)
			}

			// ✅ TODO: Log discrepancy if received != sent
			// Example: If item.Qty != existingLots[lotID].QtyAmbil
			// s.discrepancyRepo.Log(ctx, tx, ...)
		}

		return nil
	})
}
