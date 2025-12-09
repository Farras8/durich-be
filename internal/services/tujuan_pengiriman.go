package services

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"durich-be/pkg/errors"
	std_errors "errors"
)

type TujuanPengirimanService interface {
	Create(ctx context.Context, req requests.CreateTujuanPengirimanRequest, locationID string) (*response.TujuanPengirimanResponse, error)
	GetAll(ctx context.Context) ([]*response.TujuanPengirimanResponse, error)
	GetByID(ctx context.Context, id string) (*response.TujuanPengirimanResponse, error)
	Update(ctx context.Context, id string, req requests.UpdateTujuanPengirimanRequest, locationID string) (*response.TujuanPengirimanResponse, error)
	Delete(ctx context.Context, id string, locationID string) error
}

type tujuanPengirimanService struct {
	repo repository.TujuanPengirimanRepository
}

func NewTujuanPengirimanService(repo repository.TujuanPengirimanRepository) TujuanPengirimanService {
	return &tujuanPengirimanService{repo: repo}
}

func (s *tujuanPengirimanService) Create(ctx context.Context, req requests.CreateTujuanPengirimanRequest, locationID string) (*response.TujuanPengirimanResponse, error) {
	// Validation: Only Central Users can manage Master Data
	if locationID != "" {
		return nil, errors.ValidationError("akses ditolak: hanya pusat yang dapat mengelola master data tujuan pengiriman")
	}

	// Validate Type
	if req.Tipe != constants.TujuanTypeInternal && req.Tipe != constants.TujuanTypeExternal {
		return nil, errors.ValidationError("tipe tujuan tidak valid (harus 'internal' atau 'external')")
	}

	tujuan := &domain.TujuanPengiriman{
		Nama:   req.Nama,
		Tipe:   req.Tipe,
		Alamat: req.Alamat,
		Kontak: req.Kontak,
	}

	if err := s.repo.Create(ctx, tujuan); err != nil {
		return nil, err
	}

	return response.NewTujuanPengirimanResponse(tujuan), nil
}

func (s *tujuanPengirimanService) GetAll(ctx context.Context) ([]*response.TujuanPengirimanResponse, error) {
	tujuans, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return response.NewTujuanPengirimanListResponse(tujuans), nil
}

func (s *tujuanPengirimanService) GetByID(ctx context.Context, id string) (*response.TujuanPengirimanResponse, error) {
	tujuan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tujuan == nil {
		return nil, std_errors.New("tujuan pengiriman not found")
	}

	return response.NewTujuanPengirimanResponse(tujuan), nil
}

func (s *tujuanPengirimanService) Update(ctx context.Context, id string, req requests.UpdateTujuanPengirimanRequest, locationID string) (*response.TujuanPengirimanResponse, error) {
	// Validation: Only Central Users can manage Master Data
	if locationID != "" {
		return nil, errors.ValidationError("akses ditolak: hanya pusat yang dapat mengelola master data tujuan pengiriman")
	}

	// Validate Type
	if req.Tipe != constants.TujuanTypeInternal && req.Tipe != constants.TujuanTypeExternal {
		return nil, errors.ValidationError("tipe tujuan tidak valid (harus 'internal' atau 'external')")
	}

	tujuan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tujuan == nil {
		return nil, std_errors.New("tujuan pengiriman not found")
	}

	tujuan.Nama = req.Nama
	tujuan.Tipe = req.Tipe
	tujuan.Alamat = req.Alamat
	tujuan.Kontak = req.Kontak

	if err := s.repo.Update(ctx, id, tujuan); err != nil {
		return nil, err
	}

	return response.NewTujuanPengirimanResponse(tujuan), nil
}

func (s *tujuanPengirimanService) Delete(ctx context.Context, id string, locationID string) error {
	// Validation: Only Central Users can manage Master Data
	if locationID != "" {
		return errors.ValidationError("akses ditolak: hanya pusat yang dapat mengelola master data tujuan pengiriman")
	}

	return s.repo.Delete(ctx, id)
}
