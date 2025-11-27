package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
	"fmt"
)

type BuahRawRepository interface {
	BulkCreate(ctx context.Context, data []domain.BuahRaw) error
	Create(ctx context.Context, data *domain.BuahRaw) error
	GetLastKodeByJenis(ctx context.Context, kodeJenis string) (string, error)
	GetJenisDurianByID(ctx context.Context, id string) (domain.JenisDurian, error)
	GetBlokFullDetail(ctx context.Context, blokID string) (domain.Blok, domain.Divisi, domain.Estate, domain.Company, error)
	GetList(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.BuahRaw, int, error)
	GetByID(ctx context.Context, id string) (domain.BuahRaw, error)
	Update(ctx context.Context, data *domain.BuahRaw) error
	Delete(ctx context.Context, id string) error
	GetLotDetails(ctx context.Context, lotID string) ([]domain.BuahRaw, error)
}

type buahRawRepository struct {
	db *database.Database
}

func NewBuahRawRepository(db *database.Database) BuahRawRepository {
	return &buahRawRepository{db: db}
}

func (r *buahRawRepository) BulkCreate(ctx context.Context, data []domain.BuahRaw) error {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewInsert().Model(&data).Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *buahRawRepository) GetList(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.BuahRaw, int, error) {
	var list []domain.BuahRaw
	q := r.db.InitQuery(ctx).NewSelect().Model(&list).
		Relation("JenisDurianDetail").
		Relation("BlokPanenDetail").
		Relation("BlokPanenDetail.Divisi").
		Relation("BlokPanenDetail.Divisi.Estate").
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Relation("PohonPanenDetail").
		Where("buah_raw.deleted_at IS NULL").
		Order("buah_raw.created_at DESC")

	if val, ok := filter["tgl_panen"].(string); ok && val != "" {
		q.Where("buah_raw.tgl_panen = ?", val)
	}
	if val, ok := filter["blok_panen_id"].(string); ok && val != "" {
		q.Where("buah_raw.blok_panen = ?", val)
	}
	if val, ok := filter["jenis_durian_id"].(string); ok && val != "" {
		q.Where("buah_raw.jenis_durian = ?", val)
	}
	if val, ok := filter["is_sorted"]; ok {
		q.Where("buah_raw.is_sorted = ?", val)
	}

	count, err := q.Limit(limit).Offset(offset).ScanAndCount(ctx)
	return list, count, err
}

func (r *buahRawRepository) GetByID(ctx context.Context, id string) (domain.BuahRaw, error) {
	var data domain.BuahRaw
	err := r.db.InitQuery(ctx).NewSelect().Model(&data).
		Relation("JenisDurianDetail").
		Relation("BlokPanenDetail").
		Relation("BlokPanenDetail.Divisi").
		Relation("BlokPanenDetail.Divisi.Estate").
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Relation("PohonPanenDetail").
		Where("buah_raw.id = ?", id).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)
	return data, err
}

func (r *buahRawRepository) Update(ctx context.Context, data *domain.BuahRaw) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().Model(data).WherePK().Exec(ctx)
	return err
}

func (r *buahRawRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.InitQuery(ctx).NewUpdate().
		Model((*domain.BuahRaw)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *buahRawRepository) Create(ctx context.Context, data *domain.BuahRaw) error {
	_, err := r.db.InitQuery(ctx).NewInsert().Model(data).Exec(ctx)
	return err
}

func (r *buahRawRepository) GetLastKodeByJenis(ctx context.Context, kodeJenis string) (string, error) {
	var buah domain.BuahRaw
	// Cari kode buah terakhir yang diawali dengan kodeJenis (misal "MK-")
	// Order by kode_buah desc limit 1
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&buah).
		Column("kode_buah").
		Where("kode_buah LIKE ?", fmt.Sprintf("%s-%%", kodeJenis)).
		Order("kode_buah DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return "", err // Bisa jadi belum ada data, handle di service
	}
	return buah.KodeBuah, nil
}

func (r *buahRawRepository) GetJenisDurianByID(ctx context.Context, id string) (domain.JenisDurian, error) {
	var jenis domain.JenisDurian
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&jenis).
		Where("id = ?", id).
		Scan(ctx)
	return jenis, err
}

func (r *buahRawRepository) GetBlokFullDetail(ctx context.Context, blokID string) (domain.Blok, domain.Divisi, domain.Estate, domain.Company, error) {
	var blok domain.Blok
	var divisi domain.Divisi
	var estate domain.Estate
	var company domain.Company

	// Manual Join karena saya belum tahu setup relation Bun di project ini
	// Cara aman: Query bertingkat atau Join manual

	err := r.db.InitQuery(ctx).NewSelect().
		Model(&blok).
		Where("id = ?", blokID).
		Scan(ctx)
	if err != nil {
		return blok, divisi, estate, company, err
	}

	err = r.db.InitQuery(ctx).NewSelect().Model(&divisi).Where("id = ?", blok.DivisiID).Scan(ctx)
	if err != nil {
		return blok, divisi, estate, company, err
	}

	err = r.db.InitQuery(ctx).NewSelect().Model(&estate).Where("id = ?", divisi.EstateID).Scan(ctx)
	if err != nil {
		return blok, divisi, estate, company, err
	}

	err = r.db.InitQuery(ctx).NewSelect().Model(&company).Where("id = ?", estate.CompanyID).Scan(ctx)

	return blok, divisi, estate, company, err
}

func (r *buahRawRepository) GetLotDetails(ctx context.Context, lotID string) ([]domain.BuahRaw, error) {
	var buahList []domain.BuahRaw

	err := r.db.InitQuery(ctx).NewSelect().
		Model(&buahList).
		Join("INNER JOIN tb_lot_detail ON tb_lot_detail.buah_raw_id = buah_raw.id").
		Relation("BlokPanenDetail").
		Relation("BlokPanenDetail.Divisi").
		Relation("BlokPanenDetail.Divisi.Estate").
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Where("tb_lot_detail.lot_id = ?", lotID).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)

	return buahList, err
}
