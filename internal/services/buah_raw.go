package services

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

type BuahRawService interface {
	Create(ctx context.Context, req requests.BuahRawCreateRequest) (string, error)
	BulkCreate(ctx context.Context, req requests.BuahRawBulkCreateRequest) ([]string, error)
	GetList(ctx context.Context, filter map[string]interface{}, limit, page int) (response.PaginationResponse, error)
	GetDetail(ctx context.Context, id string) (response.BuahRawResponse, error)
	Update(ctx context.Context, id string, req requests.BuahRawUpdateRequest) error
	Delete(ctx context.Context, id string) error
}

type buahRawService struct {
	repo repository.BuahRawRepository
}

func NewBuahRawService(repo repository.BuahRawRepository) BuahRawService {
	return &buahRawService{repo: repo}
}

func (s *buahRawService) Create(ctx context.Context, req requests.BuahRawCreateRequest) (string, error) {
	tglPanen := req.TglPanen
	if tglPanen == "" {
		tglPanen = time.Now().Format("2006-01-02")
	}

	jenis, err := s.repo.GetJenisDurianByID(ctx, req.JenisDurianID)
	if err != nil {
		return "", fmt.Errorf("jenis durian tidak ditemukan: %v", err)
	}

	lastKode, _ := s.repo.GetLastKodeByJenis(ctx, jenis.Kode)

	currentSequence := 0
	if lastKode != "" {
		parts := strings.Split(lastKode, "-")
		if len(parts) == 2 {
			currentSequence, _ = strconv.Atoi(parts[1])
		}
	}

	currentSequence++
	newKodeBuah := fmt.Sprintf("%s-%05d", jenis.Kode, currentSequence)

	newID := ksuid.New().String()
	now := time.Now()

	buah := domain.BuahRaw{
		ID:          newID,
		KodeBuah:    newKodeBuah,
		JenisDurian: req.JenisDurianID,
		BlokPanen:   req.BlokPanenID,
		PohonPanen:  req.PohonPanenID,
		TglPanen:    tglPanen,
		IsSorted:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.repo.Create(ctx, &buah)
	if err != nil {
		return "", err
	}

	return newID, nil
}

func (s *buahRawService) BulkCreate(ctx context.Context, req requests.BuahRawBulkCreateRequest) ([]string, error) {
	var buahToInsert []domain.BuahRaw
	var insertedIDs []string
	tglPanen := req.TglPanen
	if tglPanen == "" {
		tglPanen = time.Now().Format("2006-01-02")
	}

	for _, item := range req.Items {
		jenis, err := s.repo.GetJenisDurianByID(ctx, item.JenisDurianID)
		if err != nil {
			return nil, fmt.Errorf("jenis durian tidak ditemukan: %v", err)
		}


		lastKode, _ := s.repo.GetLastKodeByJenis(ctx, jenis.Kode)
		
		currentSequence := 0
		if lastKode != "" {
			parts := strings.Split(lastKode, "-")
			if len(parts) == 2 {
				currentSequence, _ = strconv.Atoi(parts[1])
			}
		}


		for i := 0; i < item.Jumlah; i++ {
			currentSequence++
			newKodeBuah := fmt.Sprintf("%s-%05d", jenis.Kode, currentSequence)
			
			newID := ksuid.New().String()
			now := time.Now()

			buah := domain.BuahRaw{
				ID:          newID,
				KodeBuah:    newKodeBuah,
				JenisDurian: item.JenisDurianID, 
				BlokPanen:   req.BlokPanenID,    
				PohonPanen:  item.PohonPanenID,
				TglPanen:    tglPanen,
				IsSorted:    false,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			
			buahToInsert = append(buahToInsert, buah)
			insertedIDs = append(insertedIDs, newID)
		}
	}

	if len(buahToInsert) > 0 {
		err := s.repo.BulkCreate(ctx, buahToInsert)
		if err != nil {
			return nil, err
		}
	}

	return insertedIDs, nil
}

func (s *buahRawService) GetList(ctx context.Context, filter map[string]interface{}, limit, page int) (response.PaginationResponse, error) {
	offset := (page - 1) * limit
	list, count, err := s.repo.GetList(ctx, filter, limit, offset)
	if err != nil {
		return response.PaginationResponse{}, err
	}

	var data []response.BuahRawResponse
	for _, item := range list {
		data = append(data, s.mapToResponse(item))
	}

	return response.PaginationResponse{
		Data: data,
		Meta: response.PaginationMeta{
			Page:      page,
			Limit:     limit,
			TotalData: count,
			TotalPage: (count + limit - 1) / limit,
		},
	}, nil
}

func (s *buahRawService) GetDetail(ctx context.Context, id string) (response.BuahRawResponse, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return response.BuahRawResponse{}, err
	}
	if item.ID == "" {
		return response.BuahRawResponse{}, fmt.Errorf("data not found")
	}
	return s.mapToResponse(item), nil
}

func (s *buahRawService) Update(ctx context.Context, id string, req requests.BuahRawUpdateRequest) error {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if item.ID == "" {
		return fmt.Errorf("data not found")
	}

	if req.TglPanen != "" {
		item.TglPanen = req.TglPanen
	}
	if req.BlokPanenID != "" {
		item.BlokPanen = req.BlokPanenID
	}
	if req.PohonPanenID != nil {
		item.PohonPanen = req.PohonPanenID
	}
	if req.JenisDurianID != "" {
		item.JenisDurian = req.JenisDurianID
	}

	return s.repo.Update(ctx, &item)
}

func (s *buahRawService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *buahRawService) mapToResponse(item domain.BuahRaw) response.BuahRawResponse {
	kodeLengkap := ""
	blokID := ""
	blokNama := ""
	divisiNama := ""
	estateNama := ""
	companyNama := ""

	if item.BlokPanenDetail != nil {
		blokID = item.BlokPanenDetail.ID
		blokNama = item.BlokPanenDetail.NamaBlok
		if item.BlokPanenDetail.Divisi != nil {
			divisiNama = item.BlokPanenDetail.Divisi.Nama
			if item.BlokPanenDetail.Divisi.Estate != nil {
				estateNama = item.BlokPanenDetail.Divisi.Estate.Nama
				if item.BlokPanenDetail.Divisi.Estate.Company != nil {
					companyNama = item.BlokPanenDetail.Divisi.Estate.Company.Nama
					kodeLengkap = fmt.Sprintf("%s-%s-%s-%s",
						item.BlokPanenDetail.Divisi.Estate.Company.Kode,
						item.BlokPanenDetail.Divisi.Estate.Kode,
						item.BlokPanenDetail.Divisi.Kode,
						item.BlokPanenDetail.Kode,
					)
				}
			}
		}
	}

	jenisDetail := response.JenisDurianDetail{}
	if item.JenisDurianDetail != nil {
		jenisDetail = response.JenisDurianDetail{
			ID:   item.JenisDurianDetail.ID,
			Kode: item.JenisDurianDetail.Kode,
			Nama: item.JenisDurianDetail.NamaJenis,
		}
	}

	var pohonKode *string
	if item.PohonPanenDetail != nil {
		k := item.PohonPanenDetail.Kode
		pohonKode = &k
	} else if item.PohonPanen != nil {
		// fallback if detail is nil but id is there? 
		// usually means not loaded or not found
		pohonKode = nil
	}

	return response.BuahRawResponse{
		ID:          item.ID,
		KodeBuah:    item.KodeBuah,
		JenisDurian: jenisDetail,
		LokasiPanen: response.LokasiPanen{
			KodeLengkap: kodeLengkap,
			BlokID:      blokID,
			BlokNama:    blokNama,
			DivisiNama:  divisiNama,
			EstateNama:  estateNama,
			CompanyNama: companyNama,
		},
		PohonPanen: pohonKode,
		TglPanen:   item.TglPanen,
		IsSorted:   item.IsSorted,
		CreatedAt:  item.CreatedAt.Format(time.RFC3339),
	}
}
