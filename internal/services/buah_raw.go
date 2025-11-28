package services

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

type BuahRawService interface {
	Create(ctx context.Context, req requests.BuahRawCreateRequest) (string, error)
	BulkCreate(ctx context.Context, req requests.BuahRawBulkCreateRequest) ([]string, error)
	GetList(ctx context.Context, filter map[string]interface{}, limit, page int) (response.PaginationResponse, error)
	GetUnsorted(ctx context.Context, filter map[string]interface{}, limit, page int) (response.PaginationResponse, error)
	GetDetail(ctx context.Context, id string) (response.BuahRawResponse, error)
	Update(ctx context.Context, id string, req requests.BuahRawUpdateRequest) error
	Delete(ctx context.Context, id string) error
	ClearJenisCache()
}

type buahRawService struct {
	repo       repository.BuahRawRepository
	jenisCache sync.Map
	mu         sync.Mutex
}

func NewBuahRawService(repo repository.BuahRawRepository) BuahRawService {
	return &buahRawService{
		repo: repo,
	}
}

func (s *buahRawService) Create(ctx context.Context, req requests.BuahRawCreateRequest) (string, error) {
	tglPanen := req.TglPanen
	if tglPanen == "" {
		tglPanen = time.Now().Format("2006-01-02")
	}

	jenis, err := s.getJenisDurianCached(ctx, req.JenisDurianID)
	if err != nil {
		return "", fmt.Errorf("jenis durian tidak ditemukan: %v", err)
	}

	sequence, err := s.repo.GetNextSequenceWithLock(ctx, jenis.Kode)
	if err != nil {
		return "", fmt.Errorf("gagal generate sequence: %v", err)
	}

	newKodeBuah := fmt.Sprintf("%s-%05d", jenis.Kode, sequence)
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
	tglPanen := req.TglPanen
	if tglPanen == "" {
		tglPanen = time.Now().Format("2006-01-02")
	}

	jenisIDs := s.extractUniqueJenisIDs(req.Items)
	jenisMap, err := s.getJenisDurianBatch(ctx, jenisIDs)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data jenis durian: %v", err)
	}

	sequenceMap := make(map[string]int)
	s.mu.Lock()
	for kode := range s.collectJenisCodes(jenisMap) {
		sequence, err := s.repo.GetNextSequenceWithLock(ctx, kode)
		if err != nil {
			s.mu.Unlock()
			return nil, fmt.Errorf("gagal generate sequence untuk %s: %v", kode, err)
		}
		sequenceMap[kode] = sequence
	}
	s.mu.Unlock()

	buahToInsert, insertedIDs := s.buildBuahRawList(req, jenisMap, sequenceMap, tglPanen)

	if len(buahToInsert) > 0 {
		err := s.repo.BulkCreate(ctx, buahToInsert)
		if err != nil {
			return nil, err
		}
	}

	return insertedIDs, nil
}

func (s *buahRawService) extractUniqueJenisIDs(items []requests.BuahRawBulkCreateItem) []string {
	uniqueMap := make(map[string]bool)
	for _, item := range items {
		uniqueMap[item.JenisDurianID] = true
	}

	jenisIDs := make([]string, 0, len(uniqueMap))
	for id := range uniqueMap {
		jenisIDs = append(jenisIDs, id)
	}
	return jenisIDs
}

func (s *buahRawService) getJenisDurianBatch(ctx context.Context, ids []string) (map[string]domain.JenisDurian, error) {
	uncachedIDs := make([]string, 0)
	result := make(map[string]domain.JenisDurian)

	for _, id := range ids {
		if cached, ok := s.jenisCache.Load(id); ok {
			result[id] = cached.(domain.JenisDurian)
		} else {
			uncachedIDs = append(uncachedIDs, id)
		}
	}

	if len(uncachedIDs) > 0 {
		fetched, err := s.repo.GetJenisDurianByIDs(ctx, uncachedIDs)
		if err != nil {
			return nil, err
		}

		for id, jenis := range fetched {
			s.jenisCache.Store(id, jenis)
			result[id] = jenis
		}
	}

	if len(result) != len(ids) {
		return nil, fmt.Errorf("beberapa jenis durian tidak ditemukan")
	}

	return result, nil
}

func (s *buahRawService) collectJenisCodes(jenisMap map[string]domain.JenisDurian) map[string]bool {
	codes := make(map[string]bool)
	for _, jenis := range jenisMap {
		codes[jenis.Kode] = true
	}
	return codes
}

func (s *buahRawService) buildBuahRawList(
	req requests.BuahRawBulkCreateRequest,
	jenisMap map[string]domain.JenisDurian,
	sequenceMap map[string]int,
	tglPanen string,
) ([]domain.BuahRaw, []string) {
	var buahToInsert []domain.BuahRaw
	var insertedIDs []string
	now := time.Now()

	for _, item := range req.Items {
		jenis := jenisMap[item.JenisDurianID]
		currentSeq := sequenceMap[jenis.Kode]

		for i := 0; i < item.Jumlah; i++ {
			newKodeBuah := fmt.Sprintf("%s-%05d", jenis.Kode, currentSeq)
			newID := ksuid.New().String()

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
			currentSeq++
		}

		sequenceMap[jenis.Kode] = currentSeq
	}

	return buahToInsert, insertedIDs
}

func (s *buahRawService) getJenisDurianCached(ctx context.Context, id string) (domain.JenisDurian, error) {
	if cached, ok := s.jenisCache.Load(id); ok {
		return cached.(domain.JenisDurian), nil
	}

	jenis, err := s.repo.GetJenisDurianByID(ctx, id)
	if err != nil {
		return jenis, err
	}

	s.jenisCache.Store(id, jenis)
	return jenis, nil
}

func (s *buahRawService) GetList(ctx context.Context, filter map[string]interface{}, limit, page int) (response.PaginationResponse, error) {
	offset := (page - 1) * limit
	list, count, err := s.repo.GetList(ctx, filter, limit, offset)
	if err != nil {
		return response.PaginationResponse{}, err
	}

	data := make([]response.BuahRawResponse, 0, len(list))
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

func (s *buahRawService) GetUnsorted(ctx context.Context, filter map[string]interface{}, limit, page int) (response.PaginationResponse, error) {
	offset := (page - 1) * limit
	list, count, err := s.repo.GetUnsorted(ctx, filter, limit, offset)
	if err != nil {
		return response.PaginationResponse{}, err
	}

	data := make([]response.BuahRawResponse, 0, len(list))
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

	item.UpdatedAt = time.Now()

	return s.repo.Update(ctx, &item)
}

func (s *buahRawService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *buahRawService) ClearJenisCache() {
	s.jenisCache = sync.Map{}
}

func (s *buahRawService) mapToResponse(item domain.BuahRaw) response.BuahRawResponse {
	resp := response.BuahRawResponse{
		ID:          item.ID,
		KodeBuah:    item.KodeBuah,
		JenisDurian: s.buildJenisDetail(item.JenisDurianDetail),
		LokasiPanen: s.buildLokasiPanen(item.BlokPanenDetail),
		PohonPanen:  s.buildPohonKode(item.PohonPanenDetail, item.PohonPanen),
		TglPanen:    item.TglPanen,
		IsSorted:    item.IsSorted,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
	}

	return resp
}

func (s *buahRawService) buildJenisDetail(detail *domain.JenisDurian) response.JenisDurianDetail {
	if detail == nil {
		return response.JenisDurianDetail{}
	}

	return response.JenisDurianDetail{
		ID:   detail.ID,
		Kode: detail.Kode,
		Nama: detail.NamaJenis,
	}
}

func (s *buahRawService) buildLokasiPanen(blok *domain.Blok) response.LokasiPanen {
	lokasi := response.LokasiPanen{}

	if blok == nil {
		return lokasi
	}

	lokasi.BlokID = blok.ID
	lokasi.BlokNama = blok.NamaBlok

	if blok.Divisi != nil {
		lokasi.DivisiNama = blok.Divisi.Nama

		if blok.Divisi.Estate != nil {
			lokasi.EstateNama = blok.Divisi.Estate.Nama

			if blok.Divisi.Estate.Company != nil {
				lokasi.CompanyNama = blok.Divisi.Estate.Company.Nama
				lokasi.KodeLengkap = fmt.Sprintf("%s-%s-%s-%s",
					blok.Divisi.Estate.Company.Kode,
					blok.Divisi.Estate.Kode,
					blok.Divisi.Kode,
					blok.Kode,
				)
			}
		}
	}

	return lokasi
}

func (s *buahRawService) buildPohonKode(detail *domain.Pohon, pohonPanen *string) *string {
	if detail != nil {
		k := detail.Kode
		return &k
	}
	if pohonPanen != nil {
		return nil
	}
	return nil
}
