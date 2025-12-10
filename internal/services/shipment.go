package services

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"durich-be/pkg/errors"
	"time"
)

type ShipmentService interface {
	Create(ctx context.Context, req requests.ShipmentCreateRequest, userID string) (*response.ShipmentResponse, error)
	GetList(ctx context.Context, tujuan, status, locationID, listType string, page, limit int) ([]response.ShipmentResponse, int64, error)
	GetByID(ctx context.Context, id string) (*response.ShipmentDetailResponse, error)
	AddItem(ctx context.Context, shipmentID string, req requests.ShipmentAddItemRequest, locationID string) error
	RemoveItem(ctx context.Context, shipmentID string, detailID string) error
	UpdateStatus(ctx context.Context, shipmentID string, req requests.ShipmentUpdateStatusRequest, userID string) error
	Finalize(ctx context.Context, id string) error
	Receive(ctx context.Context, id string, req requests.ShipmentReceiveRequest) error
}

type shipmentService struct {
	repo       repository.ShipmentRepository
	tujuanRepo repository.TujuanPengirimanRepository
}

func NewShipmentService(repo repository.ShipmentRepository, tujuanRepo repository.TujuanPengirimanRepository) ShipmentService {
	return &shipmentService{
		repo:       repo,
		tujuanRepo: tujuanRepo,
	}
}

func (s *shipmentService) Create(ctx context.Context, req requests.ShipmentCreateRequest, userID string) (*response.ShipmentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	kode, err := s.repo.GetNextShipmentKode(ctx)
	if err != nil {
		return nil, err
	}

	tujuanDetail, err := s.tujuanRepo.GetByID(ctx, req.TujuanID)
	if err != nil {
		return nil, errors.ValidationError("invalid tujuan_id")
	}
	if tujuanDetail == nil {
		return nil, errors.ValidationError("tujuan pengiriman not found")
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

	err = s.repo.Create(ctx, shipment)
	if err != nil {
		return nil, err
	}

	resp := response.NewShipmentResponse(shipment)
	return &resp, nil
}

func (s *shipmentService) GetList(ctx context.Context, tujuan, status, locationID, listType string, page, limit int) ([]response.ShipmentResponse, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	shipments, total, err := s.repo.GetList(ctx, tujuan, status, locationID, listType, page, limit)
	if err != nil {
		return nil, 0, err
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

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	header := response.NewShipmentResponse(p)
	items := make([]response.ShipmentItemResponse, 0)

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

	detail := &domain.PengirimanDetail{
		PengirimanID: shipmentID,
		LotSumberID:  req.LotID,
	}

	return s.repo.AddItem(ctx, detail, locationID)
}

func (s *shipmentService) RemoveItem(ctx context.Context, shipmentID string, detailID string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	return s.repo.RemoveItem(ctx, shipmentID, detailID)
}

func (s *shipmentService) UpdateStatus(ctx context.Context, shipmentID string, req requests.ShipmentUpdateStatusRequest, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	shipment, err := s.repo.GetByID(ctx, shipmentID)
	if err != nil {
		return err
	}

	currentStatus := shipment.Status
	newStatus := req.Status

	isValidTransition := false
	switch currentStatus {
	case constants.ShipmentStatusDraft:
		if newStatus == constants.ShipmentStatusSending {
			isValidTransition = true
		}
	case constants.ShipmentStatusSending:
		if newStatus == constants.ShipmentStatusReceived {
			isValidTransition = true
		}
	case constants.ShipmentStatusReceived:
		if newStatus == constants.ShipmentStatusCompleted {
			isValidTransition = true
		}
	}

	if !isValidTransition {
		return errors.ValidationError("invalid status transition from " + currentStatus + " to " + newStatus)
	}

	return s.repo.UpdateStatus(ctx, shipmentID, newStatus, req.Notes, userID)
}

func (s *shipmentService) Finalize(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if shipment.Status != constants.ShipmentStatusDraft {
		return errors.ValidationError("shipment must be DRAFT to finalize")
	}
	if len(shipment.Details) == 0 {
		return errors.ValidationError("shipment cannot be empty")
	}

	return s.repo.Finalize(ctx, id)
}

func (s *shipmentService) Receive(ctx context.Context, id string, req requests.ShipmentReceiveRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 1. Get Shipment
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if shipment == nil {
		return errors.ValidationError("shipment not found")
	}

	// 2. Validate Status
	if shipment.Status != constants.ShipmentStatusSending && shipment.Status != constants.ShipmentStatusSending {
		// Note: Constants for sending might be SENDING or SHIPPED, checking for sending phase
		// Assuming 'SENDING' is the status after finalize
		if shipment.Status != constants.ShipmentStatusSending {
			return errors.ValidationError("shipment must be in SENDING status to receive")
		}
	}

	// 3. Validate Tujuan Type (Must be INTERNAL)
	tujuan, err := s.tujuanRepo.GetByID(ctx, shipment.TujuanID)
	if err != nil {
		return err
	}
	if tujuan.Tipe != "internal" {
		return errors.ValidationError("only internal shipments can be received via this endpoint")
	}

	// 4. Validate Items and Prepare Updates
	updates := make(map[string]repository.ShipmentReceiveItem)
	existingLots := make(map[string]domain.PengirimanDetail)

	for _, detail := range shipment.Details {
		existingLots[detail.LotSumberID] = detail
	}

	for _, item := range req.Details {
		detail, exists := existingLots[item.LotID]
		if !exists {
			return errors.ValidationError("lot id " + item.LotID + " is not part of this shipment")
		}

		finalQty := detail.QtyAmbil
		if item.QtyDiterima != nil {
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

	// 5. Execute Updates
	return s.repo.Receive(ctx, id, updates, shipment.TujuanID, req.ReceivedDate)
}
