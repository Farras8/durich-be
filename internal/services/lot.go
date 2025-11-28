package services

import (
	"context"
	"durich-be/internal/constants"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"errors"
	"fmt"
)

type LotService interface {
	Create(ctx context.Context, req requests.LotCreateRequest) (*response.LotResponse, error)
	GetList(ctx context.Context, status, jenisDurian, kondisi string) ([]response.LotResponse, error)
	GetDetail(ctx context.Context, id string) (*response.LotDetailResponse, error)
	AddItems(ctx context.Context, lotID string, req requests.LotAddItemsRequest) (*response.LotAddItemsResponse, error)
	RemoveItem(ctx context.Context, lotID string, req requests.LotRemoveItemRequest) error
	Finalize(ctx context.Context, lotID string, req requests.LotFinalizeRequest) (*response.LotFinalizeResponse, error)
}

type lotService struct {
	lotRepo     repository.LotRepository
	buahRawRepo repository.BuahRawRepository
}

func NewLotService(lotRepo repository.LotRepository, buahRawRepo repository.BuahRawRepository) LotService {
	return &lotService{
		lotRepo:     lotRepo,
		buahRawRepo: buahRawRepo,
	}
}

func (s *lotService) Create(ctx context.Context, req requests.LotCreateRequest) (*response.LotResponse, error) {
	lot := &domain.StokLot{
		JenisDurian: req.JenisDurianID,
		KondisiBuah: req.KondisiBuah,
		Status:      constants.LotStatusDraft,
	}

	err := s.lotRepo.Create(ctx, lot)
	if err != nil {
		return nil, err
	}

	// Fetch nama jenis durian untuk response
	// Agar simple & clean, kita panggil GetByID dari lotRepo yang sudah ada Relation
	createdLot, err := s.lotRepo.GetByID(ctx, lot.ID)
	namaJenis := ""
	if err == nil && createdLot.JenisDurianDetail != nil {
		namaJenis = createdLot.JenisDurianDetail.NamaJenis
	}

	return &response.LotResponse{
		ID:              lot.ID,
		JenisDurianID:   lot.JenisDurian,
		JenisDurianNama: namaJenis,
		KondisiBuah:     lot.KondisiBuah,
		BeratAwal:       lot.BeratAwal,
		QtyAwal:         lot.QtyAwal,
		BeratSisa:       lot.BeratSisa,
		QtySisa:         lot.QtySisa,
		Status:          lot.Status,
		CreatedAt:       lot.CreatedAt,
	}, nil
}

func (s *lotService) GetList(ctx context.Context, status, jenisDurian, kondisi string) ([]response.LotResponse, error) {
	lots, err := s.lotRepo.GetList(ctx, status, jenisDurian, kondisi)
	if err != nil {
		return nil, err
	}

	result := make([]response.LotResponse, len(lots))
	for i, lot := range lots {
		namaJenis := ""
		if lot.JenisDurianDetail != nil {
			namaJenis = lot.JenisDurianDetail.NamaJenis
		}
		
		result[i] = response.LotResponse{
			ID:              lot.ID,
			JenisDurianID:   lot.JenisDurian,
			JenisDurianNama: namaJenis,
			KondisiBuah:     lot.KondisiBuah,
			BeratAwal:       lot.BeratAwal,
			QtyAwal:         lot.QtyAwal,
			BeratSisa:       lot.BeratSisa,
			QtySisa:         lot.QtySisa,
			Status:          lot.Status,
			CreatedAt:       lot.CreatedAt,
		}
	}

	return result, nil
}

func (s *lotService) GetDetail(ctx context.Context, id string) (*response.LotDetailResponse, error) {
	lot, err := s.lotRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	namaJenis := ""
	if lot.JenisDurianDetail != nil {
		namaJenis = lot.JenisDurianDetail.NamaJenis
	}

	items := []response.LotItemResponse{}
	
	detailList, err := s.buahRawRepo.GetLotDetails(ctx, id)
	if err == nil {
		for _, buah := range detailList {
			asalBlok := ""
			if buah.BlokPanenDetail != nil {
				blok := buah.BlokPanenDetail
				if blok.Divisi != nil && blok.Divisi.Estate != nil && 
				   blok.Divisi.Estate.Company != nil {
					asalBlok = fmt.Sprintf("%s-%s-%s-%s",
						blok.Divisi.Estate.Company.Kode,
						blok.Divisi.Estate.Kode,
						blok.Divisi.Kode,
						blok.Kode,
					)
				}
			}

			items = append(items, response.LotItemResponse{
				ID:       buah.ID,
				KodeBuah: buah.KodeBuah,
				TglPanen: buah.TglPanen,
				AsalBlok: asalBlok,
			})
		}
	}

	return &response.LotDetailResponse{
		Header: response.LotResponse{
			ID:              lot.ID,
			JenisDurianID:   lot.JenisDurian,
			JenisDurianNama: namaJenis,
			KondisiBuah:     lot.KondisiBuah,
			BeratAwal:       lot.BeratAwal,
			QtyAwal:         lot.QtyAwal,
			BeratSisa:       lot.BeratSisa,
			QtySisa:         lot.QtySisa,
			Status:          lot.Status,
			CreatedAt:       lot.CreatedAt,
		},
		Items: items,
	}, nil
}

func (s *lotService) AddItems(ctx context.Context, lotID string, req requests.LotAddItemsRequest) (*response.LotAddItemsResponse, error) {
	lot, err := s.lotRepo.GetByID(ctx, lotID)
	if err != nil {
		return nil, err
	}

	if lot.Status != constants.LotStatusDraft {
		return nil, errors.New("hanya lot dengan status DRAFT yang bisa ditambahkan item")
	}

	for _, buahID := range req.BuahRawIDs {
		buah, err := s.lotRepo.GetBuahRawByID(ctx, buahID)
		if err != nil {
			return nil, fmt.Errorf("buah dengan ID %s tidak ditemukan", buahID)
		}

		if buah.IsSorted {
			return nil, fmt.Errorf("buah dengan ID %s sudah masuk lot lain", buahID)
		}

		if buah.JenisDurian != lot.JenisDurian {
			return nil, fmt.Errorf("buah %s memiliki jenis durian yang berbeda dengan lot ini", buah.KodeBuah)
		}
	}

	err = s.lotRepo.AddItems(ctx, lotID, req.BuahRawIDs)
	if err != nil {
		return nil, err
	}

	count, err := s.lotRepo.GetItemCount(ctx, lotID)
	if err != nil {
		return nil, err
	}

	return &response.LotAddItemsResponse{
		CurrentQty: count,
	}, nil
}

func (s *lotService) RemoveItem(ctx context.Context, lotID string, req requests.LotRemoveItemRequest) error {
	lot, err := s.lotRepo.GetByID(ctx, lotID)
	if err != nil {
		return err
	}

	if lot.Status != constants.LotStatusDraft {
		return errors.New("hanya lot dengan status DRAFT yang bisa dikurangi item")
	}

	return s.lotRepo.RemoveItem(ctx, lotID, req.BuahRawID)
}

func (s *lotService) Finalize(ctx context.Context, lotID string, req requests.LotFinalizeRequest) (*response.LotFinalizeResponse, error) {
	lot, err := s.lotRepo.GetByID(ctx, lotID)
	if err != nil {
		return nil, err
	}

	if lot.Status != constants.LotStatusDraft {
		return nil, errors.New("hanya lot dengan status DRAFT yang bisa difinalisasi")
	}

	count, err := s.lotRepo.GetItemCount(ctx, lotID)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("lot harus memiliki minimal 1 item sebelum difinalisasi")
	}

	lot.BeratAwal = req.BeratAwal
	lot.QtyAwal = count
	lot.BeratSisa = req.BeratAwal
	lot.QtySisa = count
	lot.Status = constants.LotStatusReady

	err = s.lotRepo.Update(ctx, lot)
	if err != nil {
		return nil, err
	}

	return &response.LotFinalizeResponse{
		ID:         lot.ID,
		QtyTotal:   lot.QtyAwal,
		BeratTotal: lot.BeratAwal,
		Status:     lot.Status,
	}, nil
}
