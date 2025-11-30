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
	"time"
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
	// Get jenis durian detail first to get code
	jenis, err := s.buahRawRepo.GetJenisDurianByID(ctx, req.JenisDurianID)
	if err != nil {
		return nil, fmt.Errorf("jenis durian tidak ditemukan: %v", err)
	}

	dateStr := time.Now().Format("020106") // DDMMYY

	kode, err := s.lotRepo.GetNextLotSequence(ctx, dateStr, jenis.Kode, req.KondisiBuah)
	if err != nil {
		return nil, fmt.Errorf("gagal generate kode lot: %v", err)
	}

	lot := &domain.StokLot{
		Kode:          kode,
		JenisDurianID: req.JenisDurianID,
		KondisiBuah:   req.KondisiBuah,
		Status:        constants.LotStatusDraft,
	}

	err = s.lotRepo.Create(ctx, lot)
	if err != nil {
		return nil, err
	}

	// Attach relation manual for response
	lot.JenisDurianDetail = &jenis

	return &response.LotResponse{
		ID:              lot.ID,
		Kode:            lot.Kode,
		JenisDurianID:   lot.JenisDurianID,
		JenisDurianNama: jenis.NamaJenis,
		KondisiBuah:     lot.KondisiBuah,
		BeratAwal:       lot.BeratAwal,
		QtyAwal:         lot.QtyAwal,
		BeratSisa:       lot.BeratSisa,
		QtySisa:         lot.QtySisa,
		Status:          lot.Status,
		CreatedAt:       lot.CreatedAt,
	}, nil
}

func (s *lotService) GetList(ctx context.Context, status, jenisDurianID, kondisi string) ([]response.LotResponse, error) {
	lots, err := s.lotRepo.GetList(ctx, status, jenisDurianID, kondisi)
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
			Kode:            lot.Kode,
			JenisDurianID:   lot.JenisDurianID,
			JenisDurianNama: namaJenis,
			KondisiBuah:     lot.KondisiBuah,
			BeratAwal:       lot.BeratAwal,
			QtyAwal:         lot.QtyAwal,
			BeratSisa:       lot.BeratSisa,
			QtySisa:         lot.QtySisa,
			CurrentQty:      lot.CurrentQty,
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
			if buah.PohonPanenDetail != nil && buah.PohonPanenDetail.Blok != nil {
				blok := buah.PohonPanenDetail.Blok
				if blok.Divisi != nil && blok.Divisi.Estate != nil &&
					blok.Divisi.Estate.Company != nil {
					asalBlok = fmt.Sprintf("%s%s%s%s",
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

	// Use the length of items as CurrentQty for consistency in GetDetail
	currentQty := len(items)
	// Or rely on lot.CurrentQty if GetByID already fetches it correctly (which we updated it to do)
	// But len(items) is more "real-time" if GetLotDetails is the source of truth for items list.
	// However, GetByID also queries the count. Let's stick to len(items) if available, or fallback.
	// Actually, since we display items, len(items) IS the current qty being displayed.

	return &response.LotDetailResponse{
		Header: response.LotResponse{
			ID:              lot.ID,
			Kode:            lot.Kode,
			JenisDurianID:   lot.JenisDurianID,
			JenisDurianNama: namaJenis,
			KondisiBuah:     lot.KondisiBuah,
			BeratAwal:       lot.BeratAwal,
			QtyAwal:         lot.QtyAwal,
			BeratSisa:       lot.BeratSisa,
			QtySisa:         lot.QtySisa,
			CurrentQty:      currentQty,
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

		if buah.JenisDurian != lot.JenisDurianID {
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
