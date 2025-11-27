package services

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"durich-be/pkg/errors"
	"time"
)

type ShipmentService interface {
	Create(ctx context.Context, req requests.ShipmentCreateRequest, userID string) (*response.ShipmentResponse, error)
	GetList(ctx context.Context, tujuan, status string) ([]response.ShipmentResponse, error)
	GetByID(ctx context.Context, id string) (*response.ShipmentDetailResponse, error)
	AddItem(ctx context.Context, shipmentID string, req requests.ShipmentAddItemRequest) error
	RemoveItem(ctx context.Context, shipmentID string, detailID string) error
	Finalize(ctx context.Context, id string) error
}

type shipmentService struct {
	repo repository.ShipmentRepository
}

func NewShipmentService(repo repository.ShipmentRepository) ShipmentService {
	return &shipmentService{repo: repo}
}

func (s *shipmentService) Create(ctx context.Context, req requests.ShipmentCreateRequest, userID string) (*response.ShipmentResponse, error) {
	tglKirim := req.TglKirim
	if tglKirim.IsZero() {
		tglKirim = time.Now()
	}

	shipment := &domain.Pengiriman{
		Tujuan:    req.Tujuan,
		TglKirim:  tglKirim,
		Status:    "DRAFT",
		CreatedBy: userID,
	}

	if err := s.repo.Create(ctx, shipment); err != nil {
		return nil, err
	}

	resp := response.NewShipmentResponse(shipment)
	return &resp, nil
}

func (s *shipmentService) GetList(ctx context.Context, tujuan, status string) ([]response.ShipmentResponse, error) {
	shipments, err := s.repo.GetList(ctx, tujuan, status)
	if err != nil {
		return nil, err
	}

	var resps []response.ShipmentResponse
	for _, p := range shipments {
		resps = append(resps, response.NewShipmentResponse(&p))
	}
	return resps, nil
}

func (s *shipmentService) GetByID(ctx context.Context, id string) (*response.ShipmentDetailResponse, error) {
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
		if d.Lot != nil && d.Lot.JenisDurianDetail != nil {
			item.JenisDurian = d.Lot.JenisDurianDetail.NamaJenis
		}
		items = append(items, item)
	}

	return &response.ShipmentDetailResponse{
		Header: header,
		Items:  items,
	}, nil
}

func (s *shipmentService) AddItem(ctx context.Context, shipmentID string, req requests.ShipmentAddItemRequest) error {
	shipment, err := s.repo.GetByID(ctx, shipmentID)
	if err != nil {
		return err
	}
	if shipment.Status != "DRAFT" {
		return errors.ValidationError("shipment must be DRAFT to add items")
	}

	detail := &domain.PengirimanDetail{
		PengirimanID: shipmentID,
		LotSumberID:  req.LotID,
		QtyAmbil:     req.Qty,
		BeratAmbil:   req.Berat,
	}

	return s.repo.AddItem(ctx, detail)
}

func (s *shipmentService) RemoveItem(ctx context.Context, shipmentID string, detailID string) error {
	shipment, err := s.repo.GetByID(ctx, shipmentID)
	if err != nil {
		return err
	}
	if shipment.Status != "DRAFT" {
		return errors.ValidationError("shipment must be DRAFT to remove items")
	}

	
	detail, err := s.repo.GetDetailByID(ctx, detailID)
	if err != nil {
		return err
	}
	if detail.PengirimanID != shipmentID {
		return errors.ValidationError("detail does not belong to this shipment")
	}

	return s.repo.RemoveItem(ctx, detailID)
}

func (s *shipmentService) Finalize(ctx context.Context, id string) error {
	shipment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if shipment.Status != "DRAFT" {
		return errors.ValidationError("shipment must be DRAFT to finalize")
	}
	if len(shipment.Details) == 0 {
		return errors.ValidationError("shipment cannot be empty")
	}

	return s.repo.Finalize(ctx, id)
}
