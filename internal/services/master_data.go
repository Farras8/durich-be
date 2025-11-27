package services

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/internal/dto/requests"
	"durich-be/internal/dto/response"
	"durich-be/internal/repository"
	"errors"
	"fmt"
)

type MasterDataService interface {
	CreateCompany(ctx context.Context, req requests.CompanyCreateRequest) (*response.CompanyResponse, error)
	GetCompanies(ctx context.Context) ([]response.CompanyResponse, error)
	GetCompanyByID(ctx context.Context, id string) (*response.CompanyResponse, error)
	UpdateCompany(ctx context.Context, id string, req requests.CompanyUpdateRequest) (*response.CompanyResponse, error)
	DeleteCompany(ctx context.Context, id string) error

	CreateEstate(ctx context.Context, req requests.EstateCreateRequest) (*response.EstateResponse, error)
	GetEstates(ctx context.Context, companyID string) ([]response.EstateResponse, error)
	GetEstateByID(ctx context.Context, id string) (*response.EstateResponse, error)
	UpdateEstate(ctx context.Context, id string, req requests.EstateUpdateRequest) (*response.EstateResponse, error)
	DeleteEstate(ctx context.Context, id string) error

	CreateDivisi(ctx context.Context, req requests.DivisiCreateRequest) (*response.DivisiResponse, error)
	GetDivisiList(ctx context.Context, estateID string) ([]response.DivisiResponse, error)
	GetDivisiByID(ctx context.Context, id string) (*response.DivisiResponse, error)
	UpdateDivisi(ctx context.Context, id string, req requests.DivisiUpdateRequest) (*response.DivisiResponse, error)
	DeleteDivisi(ctx context.Context, id string) error

	CreateBlok(ctx context.Context, req requests.BlokCreateRequest) (*response.BlokResponse, error)
	GetBloks(ctx context.Context, divisiID string) ([]response.BlokResponse, error)
	GetBlokByID(ctx context.Context, id string) (*response.BlokResponse, error)
	UpdateBlok(ctx context.Context, id string, req requests.BlokUpdateRequest) (*response.BlokResponse, error)
	DeleteBlok(ctx context.Context, id string) error

	CreateJenisDurian(ctx context.Context, req requests.JenisDurianCreateRequest) (*response.JenisDurianResponse, error)
	GetJenisDurianList(ctx context.Context) ([]response.JenisDurianResponse, error)
	GetJenisDurianByID(ctx context.Context, id string) (*response.JenisDurianResponse, error)
	UpdateJenisDurian(ctx context.Context, id string, req requests.JenisDurianUpdateRequest) (*response.JenisDurianResponse, error)
	DeleteJenisDurian(ctx context.Context, id string) error

	CreatePohon(ctx context.Context, req requests.PohonCreateRequest) (*response.PohonResponse, error)
	GetPohonList(ctx context.Context) ([]response.PohonResponse, error)
	GetPohonByID(ctx context.Context, id string) (*response.PohonResponse, error)
	UpdatePohon(ctx context.Context, id string, req requests.PohonUpdateRequest) (*response.PohonResponse, error)
	DeletePohon(ctx context.Context, id string) error
}

type masterDataService struct {
	repo repository.MasterDataRepository
}

func NewMasterDataService(repo repository.MasterDataRepository) MasterDataService {
	return &masterDataService{repo: repo}
}

func (s *masterDataService) CreateCompany(ctx context.Context, req requests.CompanyCreateRequest) (*response.CompanyResponse, error) {
	company := &domain.Company{
		Kode: req.Kode,
		Nama: req.Nama,
	}
	err := s.repo.CreateCompany(ctx, company)
	if err != nil {
		return nil, err
	}
	return &response.CompanyResponse{
		ID:        company.ID,
		Kode:      company.Kode,
		Nama:      company.Nama,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}, nil
}

func (s *masterDataService) GetCompanies(ctx context.Context) ([]response.CompanyResponse, error) {
	companies, err := s.repo.GetCompanies(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.CompanyResponse, 0, len(companies))
	for _, c := range companies {
		result = append(result, response.CompanyResponse{
			ID:        c.ID,
			Kode:      c.Kode,
			Nama:      c.Nama,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return result, nil
}

func (s *masterDataService) GetCompanyByID(ctx context.Context, id string) (*response.CompanyResponse, error) {
	company, err := s.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, errors.New("company not found")
	}
	return &response.CompanyResponse{
		ID:        company.ID,
		Kode:      company.Kode,
		Nama:      company.Nama,
		CreatedAt: company.CreatedAt,
		UpdatedAt: company.UpdatedAt,
	}, nil
}

func (s *masterDataService) UpdateCompany(ctx context.Context, id string, req requests.CompanyUpdateRequest) (*response.CompanyResponse, error) {
	existing, err := s.repo.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("company not found")
	}
	existing.Nama = req.Nama
	err = s.repo.UpdateCompany(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return &response.CompanyResponse{
		ID:        existing.ID,
		Kode:      existing.Kode,
		Nama:      existing.Nama,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

func (s *masterDataService) DeleteCompany(ctx context.Context, id string) error {
	return s.repo.DeleteCompany(ctx, id)
}

func (s *masterDataService) CreateEstate(ctx context.Context, req requests.EstateCreateRequest) (*response.EstateResponse, error) {
	estate := &domain.Estate{
		Kode:      req.Kode,
		Nama:      req.Nama,
		CompanyID: req.CompanyID,
	}
	err := s.repo.CreateEstate(ctx, estate)
	if err != nil {
		return nil, err
	}
	return &response.EstateResponse{
		ID:        estate.ID,
		Kode:      estate.Kode,
		Nama:      estate.Nama,
		CompanyID: estate.CompanyID,
		CreatedAt: estate.CreatedAt,
		UpdatedAt: estate.UpdatedAt,
	}, nil
}

func (s *masterDataService) GetEstates(ctx context.Context, companyID string) ([]response.EstateResponse, error) {
	estates, err := s.repo.GetEstates(ctx, companyID)
	if err != nil {
		return nil, err
	}
	result := make([]response.EstateResponse, 0, len(estates))
	for _, e := range estates {
		resp := response.EstateResponse{
			ID:        e.ID,
			Kode:      e.Kode,
			Nama:      e.Nama,
			CompanyID: e.CompanyID,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
		if e.Company != nil {
			resp.Company = &response.CompanyResponse{
				ID:        e.Company.ID,
				Kode:      e.Company.Kode,
				Nama:      e.Company.Nama,
				CreatedAt: e.Company.CreatedAt,
				UpdatedAt: e.Company.UpdatedAt,
			}
		}
		result = append(result, resp)
	}
	return result, nil
}

func (s *masterDataService) GetEstateByID(ctx context.Context, id string) (*response.EstateResponse, error) {
	estate, err := s.repo.GetEstateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if estate == nil {
		return nil, errors.New("estate not found")
	}
	resp := &response.EstateResponse{
		ID:        estate.ID,
		Kode:      estate.Kode,
		Nama:      estate.Nama,
		CompanyID: estate.CompanyID,
		CreatedAt: estate.CreatedAt,
		UpdatedAt: estate.UpdatedAt,
	}
	if estate.Company != nil {
		resp.Company = &response.CompanyResponse{
			ID:        estate.Company.ID,
			Kode:      estate.Company.Kode,
			Nama:      estate.Company.Nama,
			CreatedAt: estate.Company.CreatedAt,
			UpdatedAt: estate.Company.UpdatedAt,
		}
	}
	return resp, nil
}

func (s *masterDataService) UpdateEstate(ctx context.Context, id string, req requests.EstateUpdateRequest) (*response.EstateResponse, error) {
	existing, err := s.repo.GetEstateByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("estate not found")
	}
	existing.Nama = req.Nama
	existing.CompanyID = req.CompanyID
	err = s.repo.UpdateEstate(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return &response.EstateResponse{
		ID:        existing.ID,
		Kode:      existing.Kode,
		Nama:      existing.Nama,
		CompanyID: existing.CompanyID,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

func (s *masterDataService) DeleteEstate(ctx context.Context, id string) error {
	return s.repo.DeleteEstate(ctx, id)
}

func (s *masterDataService) CreateDivisi(ctx context.Context, req requests.DivisiCreateRequest) (*response.DivisiResponse, error) {
	divisi := &domain.Divisi{
		Kode:     req.Kode,
		Nama:     req.Nama,
		EstateID: req.EstateID,
	}
	err := s.repo.CreateDivisi(ctx, divisi)
	if err != nil {
		return nil, err
	}
	return &response.DivisiResponse{
		ID:        divisi.ID,
		Kode:      divisi.Kode,
		Nama:      divisi.Nama,
		EstateID:  divisi.EstateID,
		CreatedAt: divisi.CreatedAt,
		UpdatedAt: divisi.UpdatedAt,
	}, nil
}

func (s *masterDataService) GetDivisiList(ctx context.Context, estateID string) ([]response.DivisiResponse, error) {
	divisiList, err := s.repo.GetDivisiList(ctx, estateID)
	if err != nil {
		return nil, err
	}
	result := make([]response.DivisiResponse, 0, len(divisiList))
	for _, d := range divisiList {
		resp := response.DivisiResponse{
			ID:        d.ID,
			Kode:      d.Kode,
			Nama:      d.Nama,
			EstateID:  d.EstateID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		}
		if d.Estate != nil {
			resp.Estate = &response.EstateResponse{
				ID:        d.Estate.ID,
				Kode:      d.Estate.Kode,
				Nama:      d.Estate.Nama,
				CompanyID: d.Estate.CompanyID,
				CreatedAt: d.Estate.CreatedAt,
				UpdatedAt: d.Estate.UpdatedAt,
			}
		}
		result = append(result, resp)
	}
	return result, nil
}

func (s *masterDataService) GetDivisiByID(ctx context.Context, id string) (*response.DivisiResponse, error) {
	divisi, err := s.repo.GetDivisiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if divisi == nil {
		return nil, errors.New("divisi not found")
	}
	resp := &response.DivisiResponse{
		ID:        divisi.ID,
		Kode:      divisi.Kode,
		Nama:      divisi.Nama,
		EstateID:  divisi.EstateID,
		CreatedAt: divisi.CreatedAt,
		UpdatedAt: divisi.UpdatedAt,
	}
	if divisi.Estate != nil {
		resp.Estate = &response.EstateResponse{
			ID:        divisi.Estate.ID,
			Kode:      divisi.Estate.Kode,
			Nama:      divisi.Estate.Nama,
			CompanyID: divisi.Estate.CompanyID,
			CreatedAt: divisi.Estate.CreatedAt,
			UpdatedAt: divisi.Estate.UpdatedAt,
		}
	}
	return resp, nil
}

func (s *masterDataService) UpdateDivisi(ctx context.Context, id string, req requests.DivisiUpdateRequest) (*response.DivisiResponse, error) {
	existing, err := s.repo.GetDivisiByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("divisi not found")
	}
	existing.Nama = req.Nama
	existing.EstateID = req.EstateID
	err = s.repo.UpdateDivisi(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return &response.DivisiResponse{
		ID:        existing.ID,
		Kode:      existing.Kode,
		Nama:      existing.Nama,
		EstateID:  existing.EstateID,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

func (s *masterDataService) DeleteDivisi(ctx context.Context, id string) error {
	return s.repo.DeleteDivisi(ctx, id)
}

func (s *masterDataService) CreateBlok(ctx context.Context, req requests.BlokCreateRequest) (*response.BlokResponse, error) {
	blok := &domain.Blok{
		Kode:     req.Kode,
		NamaBlok: req.NamaBlok,
		DivisiID: req.DivisiID,
	}
	err := s.repo.CreateBlok(ctx, blok)
	if err != nil {
		return nil, err
	}
	return &response.BlokResponse{
		ID:        blok.ID,
		Kode:      blok.Kode,
		NamaBlok:  blok.NamaBlok,
		DivisiID:  blok.DivisiID,
		CreatedAt: blok.CreatedAt,
		UpdatedAt: blok.UpdatedAt,
	}, nil
}

func (s *masterDataService) GetBloks(ctx context.Context, divisiID string) ([]response.BlokResponse, error) {
	bloks, err := s.repo.GetBloks(ctx, divisiID)
	if err != nil {
		return nil, err
	}
	result := make([]response.BlokResponse, 0, len(bloks))
	for _, b := range bloks {
		kodeLengkap := ""
		if b.Divisi != nil && b.Divisi.Estate != nil && b.Divisi.Estate.Company != nil {
			kodeLengkap = fmt.Sprintf("%s-%s-%s-%s",
				b.Divisi.Estate.Company.Kode,
				b.Divisi.Estate.Kode,
				b.Divisi.Kode,
				b.Kode,
			)
		}

		resp := response.BlokResponse{
			ID:          b.ID,
			Kode:        b.Kode,
			NamaBlok:    b.NamaBlok,
			KodeLengkap: kodeLengkap,
			DivisiID:    b.DivisiID,
			CreatedAt:   b.CreatedAt,
			UpdatedAt:   b.UpdatedAt,
		}
		if b.Divisi != nil {
			resp.Divisi = &response.DivisiResponse{
				ID:        b.Divisi.ID,
				Kode:      b.Divisi.Kode,
				Nama:      b.Divisi.Nama,
				EstateID:  b.Divisi.EstateID,
				CreatedAt: b.Divisi.CreatedAt,
				UpdatedAt: b.Divisi.UpdatedAt,
			}
		}
		result = append(result, resp)
	}
	return result, nil
}

func (s *masterDataService) GetBlokByID(ctx context.Context, id string) (*response.BlokResponse, error) {
	blok, err := s.repo.GetBlokByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if blok == nil {
		return nil, errors.New("blok not found")
	}
	
	kodeLengkap := ""
	if blok.Divisi != nil && blok.Divisi.Estate != nil && blok.Divisi.Estate.Company != nil {
		kodeLengkap = fmt.Sprintf("%s-%s-%s-%s",
			blok.Divisi.Estate.Company.Kode,
			blok.Divisi.Estate.Kode,
			blok.Divisi.Kode,
			blok.Kode,
		)
	}

	resp := &response.BlokResponse{
		ID:          blok.ID,
		Kode:        blok.Kode,
		NamaBlok:    blok.NamaBlok,
		KodeLengkap: kodeLengkap,
		DivisiID:    blok.DivisiID,
		CreatedAt:   blok.CreatedAt,
		UpdatedAt:   blok.UpdatedAt,
	}
	if blok.Divisi != nil {
		resp.Divisi = &response.DivisiResponse{
			ID:        blok.Divisi.ID,
			Kode:      blok.Divisi.Kode,
			Nama:      blok.Divisi.Nama,
			EstateID:  blok.Divisi.EstateID,
			CreatedAt: blok.Divisi.CreatedAt,
			UpdatedAt: blok.Divisi.UpdatedAt,
		}
	}
	return resp, nil
}

func (s *masterDataService) UpdateBlok(ctx context.Context, id string, req requests.BlokUpdateRequest) (*response.BlokResponse, error) {
	existing, err := s.repo.GetBlokByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("blok not found")
	}
	existing.NamaBlok = req.NamaBlok
	existing.DivisiID = req.DivisiID
	err = s.repo.UpdateBlok(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return &response.BlokResponse{
		ID:        existing.ID,
		Kode:      existing.Kode,
		NamaBlok:  existing.NamaBlok,
		DivisiID:  existing.DivisiID,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

func (s *masterDataService) DeleteBlok(ctx context.Context, id string) error {
	return s.repo.DeleteBlok(ctx, id)
}

func (s *masterDataService) CreateJenisDurian(ctx context.Context, req requests.JenisDurianCreateRequest) (*response.JenisDurianResponse, error) {
	jenis := &domain.JenisDurian{
		Kode:      req.Kode,
		NamaJenis: req.NamaJenis,
	}
	err := s.repo.CreateJenisDurian(ctx, jenis)
	if err != nil {
		return nil, err
	}
	return &response.JenisDurianResponse{
		ID:        jenis.ID,
		Kode:      jenis.Kode,
		NamaJenis: jenis.NamaJenis,
		CreatedAt: jenis.CreatedAt,
		UpdatedAt: jenis.UpdatedAt,
	}, nil
}

func (s *masterDataService) GetJenisDurianList(ctx context.Context) ([]response.JenisDurianResponse, error) {
	jenisList, err := s.repo.GetJenisDurianList(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.JenisDurianResponse, 0, len(jenisList))
	for _, j := range jenisList {
		result = append(result, response.JenisDurianResponse{
			ID:        j.ID,
			Kode:      j.Kode,
			NamaJenis: j.NamaJenis,
			CreatedAt: j.CreatedAt,
			UpdatedAt: j.UpdatedAt,
		})
	}
	return result, nil
}

func (s *masterDataService) GetJenisDurianByID(ctx context.Context, id string) (*response.JenisDurianResponse, error) {
	jenis, err := s.repo.GetJenisDurianByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if jenis == nil {
		return nil, errors.New("jenis durian not found")
	}
	return &response.JenisDurianResponse{
		ID:        jenis.ID,
		Kode:      jenis.Kode,
		NamaJenis: jenis.NamaJenis,
		CreatedAt: jenis.CreatedAt,
		UpdatedAt: jenis.UpdatedAt,
	}, nil
}

func (s *masterDataService) UpdateJenisDurian(ctx context.Context, id string, req requests.JenisDurianUpdateRequest) (*response.JenisDurianResponse, error) {
	existing, err := s.repo.GetJenisDurianByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("jenis durian not found")
	}
	existing.NamaJenis = req.NamaJenis
	err = s.repo.UpdateJenisDurian(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return &response.JenisDurianResponse{
		ID:        existing.ID,
		Kode:      existing.Kode,
		NamaJenis: existing.NamaJenis,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

func (s *masterDataService) DeleteJenisDurian(ctx context.Context, id string) error {
	return s.repo.DeleteJenisDurian(ctx, id)
}

func (s *masterDataService) CreatePohon(ctx context.Context, req requests.PohonCreateRequest) (*response.PohonResponse, error) {
	pohon := &domain.Pohon{
		Kode: req.Kode,
		Nama: req.Nama,
	}
	err := s.repo.CreatePohon(ctx, pohon)
	if err != nil {
		return nil, err
	}
	return &response.PohonResponse{
		ID:        pohon.ID,
		Kode:      pohon.Kode,
		Nama:      pohon.Nama,
		CreatedAt: pohon.CreatedAt,
		UpdatedAt: pohon.UpdatedAt,
	}, nil
}

func (s *masterDataService) GetPohonList(ctx context.Context) ([]response.PohonResponse, error) {
	pohonList, err := s.repo.GetPohonList(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]response.PohonResponse, 0, len(pohonList))
	for _, p := range pohonList {
		result = append(result, response.PohonResponse{
			ID:        p.ID,
			Kode:      p.Kode,
			Nama:      p.Nama,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}
	return result, nil
}

func (s *masterDataService) GetPohonByID(ctx context.Context, id string) (*response.PohonResponse, error) {
	pohon, err := s.repo.GetPohonByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if pohon == nil {
		return nil, errors.New("pohon not found")
	}
	return &response.PohonResponse{
		ID:        pohon.ID,
		Kode:      pohon.Kode,
		Nama:      pohon.Nama,
		CreatedAt: pohon.CreatedAt,
		UpdatedAt: pohon.UpdatedAt,
	}, nil
}

func (s *masterDataService) UpdatePohon(ctx context.Context, id string, req requests.PohonUpdateRequest) (*response.PohonResponse, error) {
	existing, err := s.repo.GetPohonByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("pohon not found")
	}
	existing.Nama = req.Nama
	err = s.repo.UpdatePohon(ctx, id, existing)
	if err != nil {
		return nil, err
	}
	return &response.PohonResponse{
		ID:        existing.ID,
		Kode:      existing.Kode,
		Nama:      existing.Nama,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: existing.UpdatedAt,
	}, nil
}

func (s *masterDataService) DeletePohon(ctx context.Context, id string) error {
	return s.repo.DeletePohon(ctx, id)
}
