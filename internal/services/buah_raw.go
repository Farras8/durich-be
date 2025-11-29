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
	Create(ctx context.Context, req requests.BuahRawCreateRequest) (response.BuahRawResponse, error)
	BulkCreate(ctx context.Context, req requests.BuahRawBulkCreateRequest) ([]response.BuahRawResponse, error)
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

func (s *buahRawService) Create(ctx context.Context, req requests.BuahRawCreateRequest) (response.BuahRawResponse, error) {
	tglPanen := req.TglPanen
	if tglPanen == "" {
		tglPanen = time.Now().Format("2006-01-02")
	}

	// Default pohon ID logic
	pohonPanenID := req.PohonPanenID
	defaultPohonID := "6SRlQ8zX9vJ2mN5P6Q7R8S9T001"
	if pohonPanenID == nil || *pohonPanenID == "" {
		pohonPanenID = &defaultPohonID
	}

	// Get pohon with full hierarchy (Company -> Estate -> Divisi -> Blok -> Pohon)
	pohon, err := s.repo.GetPohonWithFullHierarchy(ctx, *pohonPanenID)
	if err != nil {
		return response.BuahRawResponse{}, fmt.Errorf("pohon tidak ditemukan: %v", err)
	}

	// Build prefix from location hierarchy
	prefix := s.buildLocationPrefix(pohon)
	if prefix == "" {
		return response.BuahRawResponse{}, fmt.Errorf("gagal membuat prefix lokasi: data hierarki tidak lengkap")
	}

	// Get next sequence for this location prefix
	sequence, err := s.repo.GetNextSequenceWithLock(ctx, prefix, tglPanen)
	if err != nil {
		return response.BuahRawResponse{}, fmt.Errorf("gagal generate sequence: %v", err)
	}

	newKodeBuah := fmt.Sprintf("%s-F%05d", prefix, sequence)
	newID := ksuid.New().String()
	now := time.Now()

	jenisDurian, err := s.getJenisDurianCached(ctx, req.JenisDurianID)
	if err != nil {
		return response.BuahRawResponse{}, fmt.Errorf("jenis durian tidak ditemukan: %v", err)
	}

	buah := domain.BuahRaw{
		ID:          newID,
		KodeBuah:    newKodeBuah,
		JenisDurian: req.JenisDurianID,
		PohonPanen:  pohonPanenID,
		TglPanen:    tglPanen,
		IsSorted:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.repo.Create(ctx, &buah)
	if err != nil {
		return response.BuahRawResponse{}, err
	}

	// Manually attach relations to return full response without re-querying
	buah.JenisDurianDetail = &jenisDurian
	buah.PohonPanenDetail = pohon

	return s.mapToResponse(buah), nil
}

func (s *buahRawService) BulkCreate(ctx context.Context, req requests.BuahRawBulkCreateRequest) ([]response.BuahRawResponse, error) {
	tglPanen := req.TglPanen
	if tglPanen == "" {
		tglPanen = time.Now().Format("2006-01-02")
	}

	// Extract unique pohon IDs and fetch their hierarchies
	pohonIDs := s.extractUniquePohonIDs(req.Items)
	pohonMap, err := s.getPohonBatchWithHierarchy(ctx, pohonIDs)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data pohon: %v", err)
	}

	// Build prefix map from pohon
	prefixMap := s.buildPrefixMap(pohonMap)

	// Get sequences for all prefixes
	sequenceMap := make(map[string]int)
	s.mu.Lock()
	for prefix := range s.collectUniquePrefixes(prefixMap) {
		sequence, err := s.repo.GetNextSequenceWithLock(ctx, prefix, tglPanen)
		if err != nil {
			s.mu.Unlock()
			return nil, fmt.Errorf("gagal generate sequence untuk %s: %v", prefix, err)
		}
		sequenceMap[prefix] = sequence
	}
	s.mu.Unlock()

	// Get jenis durian details for response mapping
	jenisIDs := make([]string, 0)
	for _, item := range req.Items {
		jenisIDs = append(jenisIDs, item.JenisDurianID)
	}
	jenisMap, err := s.getJenisDurianBatch(ctx, jenisIDs)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data jenis durian: %v", err)
	}

	buahToInsert, insertedIDs := s.buildBuahRawListFromLocation(req, prefixMap, sequenceMap, tglPanen)

	if len(buahToInsert) > 0 {
		err := s.repo.BulkCreate(ctx, buahToInsert)
		if err != nil {
			return nil, err
		}
	}

	// Map inserted items to full response
	result := make([]response.BuahRawResponse, 0, len(buahToInsert))
	for i, b := range buahToInsert {
		// Attach loaded relations
		if b.PohonPanen != nil {
			b.PohonPanenDetail = pohonMap[*b.PohonPanen]
		}
		if jenis, ok := jenisMap[b.JenisDurian]; ok {
			// create copy to assign address
			j := jenis
			b.JenisDurianDetail = &j
		}
		// Re-assign ID since it was generated in buildBuahRawListFromLocation but not in original struct passed to repo if repo modifies it (ksuid is string so ok)
		// Wait, buahToInsert has IDs already generated.
		
		// Ensure we use the ID from the inserted list if needed, but buildBuahRawListFromLocation sets it.
		b.ID = insertedIDs[i] 
		
		result = append(result, s.mapToResponse(b))
	}

	return result, nil
}

func (s *buahRawService) extractUniquePohonIDs(items []requests.BuahRawBulkCreateItem) []string {
	uniqueMap := make(map[string]bool)
	defaultPohonID := "6SRlQ8zX9vJ2mN5P6Q7R8S9T001"
	
	for _, item := range items {
		pohonID := defaultPohonID
		if item.PohonPanenID != nil && *item.PohonPanenID != "" {
			pohonID = *item.PohonPanenID
		}
		uniqueMap[pohonID] = true
	}

	pohonIDs := make([]string, 0, len(uniqueMap))
	for id := range uniqueMap {
		pohonIDs = append(pohonIDs, id)
	}
	return pohonIDs
}

func (s *buahRawService) getPohonBatchWithHierarchy(ctx context.Context, ids []string) (map[string]*domain.Pohon, error) {
	result := make(map[string]*domain.Pohon)

	for _, id := range ids {
		pohon, err := s.repo.GetPohonWithFullHierarchy(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("pohon %s tidak ditemukan: %v", id, err)
		}
		result[id] = pohon
	}

	return result, nil
}

func (s *buahRawService) buildPrefixMap(pohonMap map[string]*domain.Pohon) map[string]string {
	prefixMap := make(map[string]string)
	for pohonID, pohon := range pohonMap {
		prefix := s.buildLocationPrefix(pohon)
		prefixMap[pohonID] = prefix
	}
	return prefixMap
}

func (s *buahRawService) collectUniquePrefixes(prefixMap map[string]string) map[string]bool {
	prefixes := make(map[string]bool)
	for _, prefix := range prefixMap {
		if prefix != "" {
			prefixes[prefix] = true
		}
	}
	return prefixes
}

// Keep for backward compatibility (still used in Update method)
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

func (s *buahRawService) buildBuahRawListFromLocation(
	req requests.BuahRawBulkCreateRequest,
	prefixMap map[string]string,
	sequenceMap map[string]int,
	tglPanen string,
) ([]domain.BuahRaw, []string) {
	var buahToInsert []domain.BuahRaw
	var insertedIDs []string
	now := time.Now()
	defaultPohonID := "6SRlQ8zX9vJ2mN5P6Q7R8S9T001"

	for _, item := range req.Items {
		// Determine pohon ID for this item
		pohonID := defaultPohonID
		if item.PohonPanenID != nil && *item.PohonPanenID != "" {
			pohonID = *item.PohonPanenID
		}

		prefix := prefixMap[pohonID]
		if prefix == "" {
			continue // Skip items with invalid prefix
		}

		currentSeq := sequenceMap[prefix]

		for i := 0; i < item.Jumlah; i++ {
			newKodeBuah := fmt.Sprintf("%s-F%05d", prefix, currentSeq)
			newID := ksuid.New().String()

			buah := domain.BuahRaw{
				ID:          newID,
				KodeBuah:    newKodeBuah,
				JenisDurian: item.JenisDurianID,
				PohonPanen:  &pohonID,
				TglPanen:    tglPanen,
				IsSorted:    false,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			buahToInsert = append(buahToInsert, buah)
			insertedIDs = append(insertedIDs, newID)
			currentSeq++
		}

		sequenceMap[prefix] = currentSeq
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
	if req.PohonPanenID != nil {
		// If empty string provided, use default. If valid string, use it.
		// If user sends explicit empty string, they might mean "reset to default" or "no tree".
		// Requirement says: "kalau pohon idnya ga diisi default id pohonnya itu 6SRlQ8zX9vJ2mN5P6Q7R8S9T001"
		// So if user updates and sends "", set to default.
		val := *req.PohonPanenID
		if val == "" {
			defaultID := "6SRlQ8zX9vJ2mN5P6Q7R8S9T001"
			item.PohonPanen = &defaultID
		} else {
			item.PohonPanen = req.PohonPanenID
		}
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
		LokasiPanen: s.buildLokasiPanenFromPohon(item.PohonPanenDetail),
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

func (s *buahRawService) buildLokasiPanenFromPohon(pohon *domain.Pohon) response.LokasiPanen {
	lokasi := response.LokasiPanen{}

	if pohon == nil || pohon.Blok == nil {
		return lokasi
	}

	blok := pohon.Blok
	lokasi.BlokID = blok.ID
	lokasi.BlokNama = blok.NamaBlok

	if blok.Divisi != nil {
		lokasi.DivisiNama = blok.Divisi.Nama

		if blok.Divisi.Estate != nil {
			lokasi.EstateNama = blok.Divisi.Estate.Nama

			if blok.Divisi.Estate.Company != nil {
				lokasi.CompanyNama = blok.Divisi.Estate.Company.Nama
				lokasi.KodeLengkap = fmt.Sprintf("%s%s%s%s",
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

// buildLocationPrefix creates prefix from full hierarchy: Company+Estate+Divisi+Blok+Pohon
// Example: IPSRES0101A010000 (IPS + RES + 01 + 01A01 + 0000)
func (s *buahRawService) buildLocationPrefix(pohon *domain.Pohon) string {
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
