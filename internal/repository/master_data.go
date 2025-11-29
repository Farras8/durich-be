package repository

import (
	"context"
	"database/sql"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
)

type MasterDataRepository interface {
	CreateCompany(ctx context.Context, company *domain.Company) error
	GetCompanies(ctx context.Context) ([]domain.Company, error)
	GetCompanyByID(ctx context.Context, id string) (*domain.Company, error)
	UpdateCompany(ctx context.Context, id string, company *domain.Company) error
	DeleteCompany(ctx context.Context, id string) error

	CreateEstate(ctx context.Context, estate *domain.Estate) error
	GetEstates(ctx context.Context, companyID string) ([]domain.Estate, error)
	GetEstateByID(ctx context.Context, id string) (*domain.Estate, error)
	UpdateEstate(ctx context.Context, id string, estate *domain.Estate) error
	DeleteEstate(ctx context.Context, id string) error

	CreateDivisi(ctx context.Context, divisi *domain.Divisi) error
	GetDivisiList(ctx context.Context, estateID string) ([]domain.Divisi, error)
	GetDivisiByID(ctx context.Context, id string) (*domain.Divisi, error)
	UpdateDivisi(ctx context.Context, id string, divisi *domain.Divisi) error
	DeleteDivisi(ctx context.Context, id string) error

	CreateBlok(ctx context.Context, blok *domain.Blok) error
	GetBloks(ctx context.Context, divisiID string) ([]domain.Blok, error)
	GetBlokByID(ctx context.Context, id string) (*domain.Blok, error)
	UpdateBlok(ctx context.Context, id string, blok *domain.Blok) error
	DeleteBlok(ctx context.Context, id string) error

	CreateJenisDurian(ctx context.Context, jenis *domain.JenisDurian) error
	GetJenisDurianList(ctx context.Context) ([]domain.JenisDurian, error)
	GetJenisDurianByID(ctx context.Context, id string) (*domain.JenisDurian, error)
	UpdateJenisDurian(ctx context.Context, id string, jenis *domain.JenisDurian) error
	DeleteJenisDurian(ctx context.Context, id string) error

	CreatePohon(ctx context.Context, pohon *domain.Pohon) error
	GetPohonList(ctx context.Context) ([]domain.Pohon, error)
	GetPohonByID(ctx context.Context, id string) (*domain.Pohon, error)
	UpdatePohon(ctx context.Context, id string, pohon *domain.Pohon) error
	DeletePohon(ctx context.Context, id string) error
}

type masterDataRepository struct {
	db *database.Database
}

func NewMasterDataRepository(db *database.Database) MasterDataRepository {
	return &masterDataRepository{db: db}
}

func (r *masterDataRepository) CreateCompany(ctx context.Context, company *domain.Company) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(company).Exec(ctx)
	return err
}

func (r *masterDataRepository) GetCompanies(ctx context.Context) ([]domain.Company, error) {
	var companies []domain.Company
	err := r.db.InitQuery(ctx).NewSelect().Model(&companies).Where("deleted_at IS NULL").Scan(ctx)
	return companies, err
}

func (r *masterDataRepository) GetCompanyByID(ctx context.Context, id string) (*domain.Company, error) {
	company := &domain.Company{}
	err := r.db.InitQuery(ctx).NewSelect().Model(company).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return company, err
}

func (r *masterDataRepository) UpdateCompany(ctx context.Context, id string, company *domain.Company) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(company).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *masterDataRepository) DeleteCompany(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Company)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *masterDataRepository) CreateEstate(ctx context.Context, estate *domain.Estate) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(estate).Exec(ctx)
	return err
}

func (r *masterDataRepository) GetEstates(ctx context.Context, companyID string) ([]domain.Estate, error) {
	var estates []domain.Estate
	query := r.db.InitQuery(ctx).NewSelect().Model(&estates).Where("estate.deleted_at IS NULL")
	if companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}
	err := query.Relation("Company").Scan(ctx)
	return estates, err
}

func (r *masterDataRepository) GetEstateByID(ctx context.Context, id string) (*domain.Estate, error) {
	estate := &domain.Estate{}
	err := r.db.InitQuery(ctx).NewSelect().Model(estate).Where("estate.id = ? AND estate.deleted_at IS NULL", id).Relation("Company").Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return estate, err
}

func (r *masterDataRepository) UpdateEstate(ctx context.Context, id string, estate *domain.Estate) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(estate).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *masterDataRepository) DeleteEstate(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Estate)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *masterDataRepository) CreateDivisi(ctx context.Context, divisi *domain.Divisi) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(divisi).Exec(ctx)
	return err
}

func (r *masterDataRepository) GetDivisiList(ctx context.Context, estateID string) ([]domain.Divisi, error) {
	var divisiList []domain.Divisi
	query := r.db.InitQuery(ctx).NewSelect().Model(&divisiList).Where("divisi.deleted_at IS NULL")
	if estateID != "" {
		query = query.Where("estate_id = ?", estateID)
	}
	err := query.Relation("Estate").Scan(ctx)
	return divisiList, err
}

func (r *masterDataRepository) GetDivisiByID(ctx context.Context, id string) (*domain.Divisi, error) {
	divisi := &domain.Divisi{}
	err := r.db.InitQuery(ctx).NewSelect().Model(divisi).Where("divisi.id = ? AND divisi.deleted_at IS NULL", id).Relation("Estate").Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return divisi, err
}

func (r *masterDataRepository) UpdateDivisi(ctx context.Context, id string, divisi *domain.Divisi) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(divisi).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *masterDataRepository) DeleteDivisi(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Divisi)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *masterDataRepository) CreateBlok(ctx context.Context, blok *domain.Blok) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(blok).Exec(ctx)
	return err
}

func (r *masterDataRepository) GetBloks(ctx context.Context, divisiID string) ([]domain.Blok, error) {
	var bloks []domain.Blok
	query := r.db.InitQuery(ctx).NewSelect().Model(&bloks).Where("blok.deleted_at IS NULL")
	if divisiID != "" {
		query = query.Where("divisi_id = ?", divisiID)
	}
	err := query.Relation("Divisi").
		Relation("Divisi.Estate").
		Relation("Divisi.Estate.Company").
		Scan(ctx)
	return bloks, err
}

func (r *masterDataRepository) GetBlokByID(ctx context.Context, id string) (*domain.Blok, error) {
	blok := &domain.Blok{}
	err := r.db.InitQuery(ctx).NewSelect().Model(blok).
		Where("blok.id = ? AND blok.deleted_at IS NULL", id).
		Relation("Divisi").
		Relation("Divisi.Estate").
		Relation("Divisi.Estate.Company").
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return blok, err
}

func (r *masterDataRepository) UpdateBlok(ctx context.Context, id string, blok *domain.Blok) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(blok).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *masterDataRepository) DeleteBlok(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Blok)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *masterDataRepository) CreateJenisDurian(ctx context.Context, jenis *domain.JenisDurian) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(jenis).Exec(ctx)
	return err
}

func (r *masterDataRepository) GetJenisDurianList(ctx context.Context) ([]domain.JenisDurian, error) {
	var jenisList []domain.JenisDurian
	err := r.db.InitQuery(ctx).NewSelect().Model(&jenisList).Where("deleted_at IS NULL").Scan(ctx)
	return jenisList, err
}

func (r *masterDataRepository) GetJenisDurianByID(ctx context.Context, id string) (*domain.JenisDurian, error) {
	jenis := &domain.JenisDurian{}
	err := r.db.InitQuery(ctx).NewSelect().Model(jenis).Where("id = ? AND deleted_at IS NULL", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return jenis, err
}

func (r *masterDataRepository) UpdateJenisDurian(ctx context.Context, id string, jenis *domain.JenisDurian) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(jenis).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *masterDataRepository) DeleteJenisDurian(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.JenisDurian)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *masterDataRepository) CreatePohon(ctx context.Context, pohon *domain.Pohon) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(pohon).Exec(ctx)
	return err
}

func (r *masterDataRepository) GetPohonList(ctx context.Context) ([]domain.Pohon, error) {
	var pohonList []domain.Pohon
	err := r.db.InitQuery(ctx).NewSelect().Model(&pohonList).
		Relation("Blok").
		Relation("Blok.Divisi").
		Relation("Blok.Divisi.Estate").
		Relation("Blok.Divisi.Estate.Company").
		Where("pohon.deleted_at IS NULL").
		Scan(ctx)
	return pohonList, err
}

func (r *masterDataRepository) GetPohonByID(ctx context.Context, id string) (*domain.Pohon, error) {
	pohon := &domain.Pohon{}
	err := r.db.InitQuery(ctx).NewSelect().Model(pohon).
		Relation("Blok").
		Relation("Blok.Divisi").
		Relation("Blok.Divisi.Estate").
		Relation("Blok.Divisi.Estate.Company").
		Where("pohon.id = ? AND pohon.deleted_at IS NULL", id).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return pohon, err
}

func (r *masterDataRepository) UpdatePohon(ctx context.Context, id string, pohon *domain.Pohon) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(pohon).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *masterDataRepository) DeletePohon(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.Pohon)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}
