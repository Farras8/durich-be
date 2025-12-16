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
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type LotService interface {
	Create(ctx context.Context, req requests.LotCreateRequest, locationID string) (*response.LotResponse, error)
	GetList(ctx context.Context, status, jenisDurian, kondisi, locationID, scope, createdAt string) ([]response.LotResponse, error)
	GetDetail(ctx context.Context, id string) (*response.LotDetailResponse, error)
	AddItems(ctx context.Context, lotID string, req requests.LotAddItemsRequest, locationID string) (*response.LotAddItemsResponse, error)
	RemoveItem(ctx context.Context, lotID string, req requests.LotRemoveItemRequest, locationID string) error
	Finalize(ctx context.Context, lotID string, req requests.LotFinalizeRequest, locationID string) (*response.LotFinalizeResponse, error)
}

type lotService struct {
	db          *bun.DB
	lotRepo     repository.LotRepository
	buahRawRepo repository.BuahRawRepository
}

func NewLotService(db *bun.DB, lotRepo repository.LotRepository, buahRawRepo repository.BuahRawRepository) LotService {
	return &lotService{
		db:          db,
		lotRepo:     lotRepo,
		buahRawRepo: buahRawRepo,
	}
}

func (s *lotService) Create(ctx context.Context, req requests.LotCreateRequest, locationID string) (*response.LotResponse, error) {
	if locationID != "" {
		return nil, errors.ValidationError("akses ditolak: hanya pusat yang dapat membuat lot baru (grading)")
	}

	var result *response.LotResponse

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		jenis, err := s.buahRawRepo.GetJenisDurianByID(ctx, req.JenisDurianID)
		if err != nil {
			return fmt.Errorf("jenis durian tidak ditemukan: %v", err)
		}

		dateStr := time.Now().Format("020106")

		kode, err := s.lotRepo.GetNextLotSequence(ctx, tx, dateStr, jenis.Kode, req.KondisiBuah)
		if err != nil {
			return fmt.Errorf("gagal generate kode lot: %v", err)
		}

		lot := &domain.StokLot{
			Kode:          kode,
			JenisDurianID: req.JenisDurianID,
			KondisiBuah:   req.KondisiBuah,
			Status:        constants.LotStatusDraft,
		}

		err = s.lotRepo.Create(ctx, tx, lot)
		if err != nil {
			return err
		}

		lot.JenisDurianDetail = &jenis

		result = &response.LotResponse{
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
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *lotService) GetList(ctx context.Context, status, jenisDurianID, kondisi, locationID, scope, createdAt string) ([]response.LotResponse, error) {
	lots, err := s.lotRepo.GetList(ctx, s.db, status, jenisDurianID, kondisi, locationID, scope, createdAt)
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
			CurrentBerat:    lot.CurrentBerat,
			Status:          lot.Status,
			CreatedAt:       lot.CreatedAt,
		}
	}

	return result, nil
}

func (s *lotService) GetDetail(ctx context.Context, id string) (*response.LotDetailResponse, error) {
	lot, err := s.lotRepo.GetByID(ctx, s.db, id)
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

			jenisDurianInfo := ""
			if buah.JenisDurianDetail != nil {
				jenisDurianInfo = fmt.Sprintf("%s - %s", buah.JenisDurianDetail.Kode, buah.JenisDurianDetail.NamaJenis)
			}

			items = append(items, response.LotItemResponse{
				ID:          buah.ID,
				KodeBuah:    buah.KodeBuah,
				TglPanen:    buah.TglPanen,
				AsalBlok:    asalBlok,
				Berat:       buah.Berat,
				JenisDurian: jenisDurianInfo,
			})
		}
	}

	currentQty := len(items)

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
			CurrentBerat:    lot.CurrentBerat,
			Status:          lot.Status,
			CreatedAt:       lot.CreatedAt,
		},
		Items: items,
	}, nil
}

func (s *lotService) AddItems(ctx context.Context, lotID string, req requests.LotAddItemsRequest, locationID string) (*response.LotAddItemsResponse, error) {
	if locationID != "" {
		return nil, errors.ValidationError("akses ditolak: hanya pusat yang dapat memodifikasi lot")
	}

	var result *response.LotAddItemsResponse

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		lot, err := s.lotRepo.GetByIDForUpdate(ctx, tx, lotID)
		if err != nil {
			return err
		}

		if lot.Status != constants.LotStatusDraft {
			return std_errors.New("hanya lot dengan status DRAFT yang bisa ditambahkan item")
		}

		pohon, err := s.lotRepo.GetPohonByKode(ctx, tx, req.PohonKode, req.BlokID)
		if err != nil {
			return fmt.Errorf("pohon dengan kode %s tidak ditemukan di blok yang dipilih", req.PohonKode)
		}

		prefix := s.buildLocationPrefix(pohon)
		tglPanen := time.Now().Format("2006-01-02")

		sequence, err := s.buahRawRepo.GetNextSequenceWithLock(ctx, prefix, tglPanen)
		if err != nil {
			return fmt.Errorf("gagal generate sequence buah: %v", err)
		}

		kodeBuah := fmt.Sprintf("%s-F%05d", prefix, sequence)

		buah := &domain.BuahRaw{
			KodeBuah:    kodeBuah,
			JenisDurian: lot.JenisDurianID,
			PohonPanen:  &pohon.ID,
			TglPanen:    tglPanen,
			LotID:       &lotID,
			BlokID:      &req.BlokID,
			Berat:       req.Berat,
		}

		err = s.lotRepo.AddBuah(ctx, tx, buah)
		if err != nil {
			return err
		}

		count, err := s.lotRepo.GetItemCount(ctx, tx, lotID)
		if err != nil {
			return err
		}

		result = &response.LotAddItemsResponse{
			CurrentQty: count,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *lotService) buildLocationPrefix(pohon *domain.Pohon) string {
	if pohon == nil || pohon.Blok == nil {
		return ""
	}

	blok := pohon.Blok
	if blok.Divisi == nil || blok.Divisi.Estate == nil || blok.Divisi.Estate.Company == nil {
		return ""
	}

	return fmt.Sprintf("%s%s%s%s%s",
		blok.Divisi.Estate.Company.Kode,
		blok.Divisi.Estate.Kode,
		blok.Divisi.Kode,
		blok.Kode,
		pohon.Kode,
	)
}

func (s *lotService) RemoveItem(ctx context.Context, lotID string, req requests.LotRemoveItemRequest, locationID string) error {
	if locationID != "" {
		return errors.ValidationError("akses ditolak: hanya pusat yang dapat memodifikasi lot")
	}

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		lot, err := s.lotRepo.GetByIDForUpdate(ctx, tx, lotID)
		if err != nil {
			return err
		}

		if lot.Status != constants.LotStatusDraft {
			return std_errors.New("hanya lot dengan status DRAFT yang bisa dikurangi item")
		}

		return s.lotRepo.RemoveItem(ctx, tx, lotID, req.BuahRawID)
	})

	return err
}

func (s *lotService) Finalize(ctx context.Context, lotID string, req requests.LotFinalizeRequest, locationID string) (*response.LotFinalizeResponse, error) {
	if locationID != "" {
		return nil, errors.ValidationError("akses ditolak: hanya pusat yang dapat memfinalisasi lot")
	}

	var result *response.LotFinalizeResponse

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		lot, err := s.lotRepo.GetByIDForUpdate(ctx, tx, lotID)
		if err != nil {
			return err
		}

		if lot.Status != constants.LotStatusDraft {
			return std_errors.New("hanya lot dengan status DRAFT yang bisa difinalisasi")
		}

		count, err := s.lotRepo.GetItemCount(ctx, tx, lotID)
		if err != nil {
			return err
		}

		if count == 0 {
			return std_errors.New("lot harus memiliki minimal 1 item sebelum difinalisasi")
		}

		totalWeight, err := s.lotRepo.GetTotalWeight(ctx, tx, lotID)
		if err != nil {
			return fmt.Errorf("gagal menghitung total berat: %v", err)
		}

		lot.BeratAwal = totalWeight
		lot.QtyAwal = count
		lot.BeratSisa = totalWeight
		lot.QtySisa = count
		lot.Status = constants.LotStatusReady

		err = s.lotRepo.Update(ctx, tx, lot)
		if err != nil {
			return err
		}

		result = &response.LotFinalizeResponse{
			ID:         lot.ID,
			QtyTotal:   lot.QtyAwal,
			BeratTotal: lot.BeratAwal,
			Status:     lot.Status,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
