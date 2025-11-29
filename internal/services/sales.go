package services

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"durich-be/pkg/errors"
)

type SalesService interface {
	Create(ctx context.Context, req requests.SalesCreateRequest) (*response.SalesResponse, error)
	GetList(ctx context.Context, startDate, endDate, tipeJual string) ([]response.SalesResponse, error)
	GetByID(ctx context.Context, id string) (*response.SalesDetailResponse, error)
	Update(ctx context.Context, id string, req requests.SalesUpdateRequest) error
	Delete(ctx context.Context, id string) error
}

type salesService struct {
	repo repository.SalesRepository
}

func NewSalesService(repo repository.SalesRepository) SalesService {
	return &salesService{repo: repo}
}

func (s *salesService) Create(ctx context.Context, req requests.SalesCreateRequest) (*response.SalesResponse, error) {

	shipment, err := s.repo.GetPengirimanByID(ctx, req.PengirimanID)
	if err != nil {
		return nil, err
	}
	if shipment.Status != constants.ShipmentStatusSending {
		return nil, errors.ValidationError("shipment status must be SENDING")
	}

	exists, err := s.repo.CheckSalesExistByShipmentID(ctx, req.PengirimanID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ValidationError("invoice already exists for this shipment")
	}

	totalBerat := 0.0
	for _, d := range shipment.Details {
		totalBerat += d.BeratAmbil
	}

	sales := &domain.Penjualan{
		PengirimanID: req.PengirimanID,
		BeratTerjual: totalBerat,
		HargaTotal:   req.HargaTotal,
		TipeJual:     req.TipeJual,
	}

	if err := s.repo.Create(ctx, sales); err != nil {
		return nil, err
	}

	resp := response.NewSalesResponse(sales)
	return &resp, nil
}

func (s *salesService) GetList(ctx context.Context, startDate, endDate, tipeJual string) ([]response.SalesResponse, error) {
	salesList, err := s.repo.GetList(ctx, startDate, endDate, tipeJual)
	if err != nil {
		return nil, err
	}

	var resps []response.SalesResponse
	for _, sales := range salesList {
		resps = append(resps, response.NewSalesResponse(&sales))
	}
	return resps, nil
}

func (s *salesService) GetByID(ctx context.Context, id string) (*response.SalesDetailResponse, error) {
	sales, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items := make([]response.ShipmentItemResponse, 0)
	if sales.Pengiriman != nil {
		for _, d := range sales.Pengiriman.Details {
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
	}

	pengirimanInfo := response.SalesShipmentInfoResponse{
		ID:      sales.PengirimanID,
		Tujuan:  "",
		Status:  "",
		Details: items,
	}

	if sales.Pengiriman != nil {
		pengirimanInfo.Tujuan = sales.Pengiriman.Tujuan
		pengirimanInfo.Status = sales.Pengiriman.Status
	}

	return &response.SalesDetailResponse{
		ID: sales.ID,
		InfoPenjualan: response.SalesInfoResponse{
			HargaTotal:   sales.HargaTotal,
			BeratTerjual: sales.BeratTerjual,
			TipeJual:     sales.TipeJual,
			CreatedAt:    sales.CreatedAt,
		},
		InfoPengiriman: pengirimanInfo,
	}, nil
}

func (s *salesService) Update(ctx context.Context, id string, req requests.SalesUpdateRequest) error {
	sales, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if req.HargaTotal > 0 {
		sales.HargaTotal = req.HargaTotal
	}
	if req.TipeJual != "" {
		sales.TipeJual = req.TipeJual
	}

	return s.repo.Update(ctx, sales)
}

func (s *salesService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
