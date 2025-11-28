package repository

import (
	"context"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
	"fmt"

	"github.com/uptrace/bun"
)

type BuahRawRepository interface {
	BulkCreate(ctx context.Context, data []domain.BuahRaw) error
	Create(ctx context.Context, data *domain.BuahRaw) error
	GetLastKodeByJenis(ctx context.Context, kodeJenis string) (string, error)
	GetJenisDurianByID(ctx context.Context, id string) (domain.JenisDurian, error)
	GetJenisDurianByIDs(ctx context.Context, ids []string) (map[string]domain.JenisDurian, error)
	GetBlokFullDetail(ctx context.Context, blokID string) (domain.Blok, domain.Divisi, domain.Estate, domain.Company, error)
	GetList(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.BuahRaw, int, error)
	GetUnsorted(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.BuahRaw, int, error)
	GetByID(ctx context.Context, id string) (domain.BuahRaw, error)
	Update(ctx context.Context, data *domain.BuahRaw) error
	Delete(ctx context.Context, id string) error
	GetLotDetails(ctx context.Context, lotID string) ([]domain.BuahRaw, error)
	GetNextSequenceWithLock(ctx context.Context, kodeJenis string) (int, error)
}

type buahRawRepository struct {
	db *database.Database
}

func NewBuahRawRepository(db *database.Database) BuahRawRepository {
	return &buahRawRepository{db: db}
}

func (r *buahRawRepository) BulkCreate(ctx context.Context, data []domain.BuahRaw) error {
	const batchSize = 1000

	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		_, err = tx.NewInsert().Model(&batch).Exec(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *buahRawRepository) GetList(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.BuahRaw, int, error) {
	var list []domain.BuahRaw
	q := r.db.InitQuery(ctx).NewSelect().Model(&list)

	q = r.applyRelations(q, filter)
	q = r.applyFilters(q, filter)
	q = q.Where("buah_raw.deleted_at IS NULL").
		Order("buah_raw.created_at DESC")

	count, err := q.Limit(limit).Offset(offset).ScanAndCount(ctx)
	return list, count, err
}

func (r *buahRawRepository) GetUnsorted(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]domain.BuahRaw, int, error) {
	var list []domain.BuahRaw
	q := r.db.InitQuery(ctx).NewSelect().Model(&list)

	q = r.applyRelations(q, filter)
	q = r.applyFilters(q, filter)
	q = q.Where("buah_raw.deleted_at IS NULL").
		Where("buah_raw.is_sorted = false").
		Order("buah_raw.created_at DESC")

	count, err := q.Limit(limit).Offset(offset).ScanAndCount(ctx)
	return list, count, err
}

func (r *buahRawRepository) applyRelations(q *bun.SelectQuery, filter map[string]interface{}) *bun.SelectQuery {
	includeRelations, ok := filter["include_relations"].(map[string]bool)
	if !ok || len(includeRelations) == 0 {
		return q.
			Relation("JenisDurianDetail").
			Relation("BlokPanenDetail").
			Relation("BlokPanenDetail.Divisi").
			Relation("BlokPanenDetail.Divisi.Estate").
			Relation("BlokPanenDetail.Divisi.Estate.Company").
			Relation("PohonPanenDetail")
	}

	if includeRelations["jenis"] {
		q = q.Relation("JenisDurianDetail")
	}
	if includeRelations["blok"] {
		q = q.Relation("BlokPanenDetail").
			Relation("BlokPanenDetail.Divisi").
			Relation("BlokPanenDetail.Divisi.Estate").
			Relation("BlokPanenDetail.Divisi.Estate.Company")
	}
	if includeRelations["pohon"] {
		q = q.Relation("PohonPanenDetail")
	}

	return q
}

func (r *buahRawRepository) applyFilters(q *bun.SelectQuery, filter map[string]interface{}) *bun.SelectQuery {
	if val, ok := filter["tgl_panen"].(string); ok && val != "" {
		q = q.Where("buah_raw.tgl_panen = ?", val)
	}
	if val, ok := filter["blok_panen_id"].(string); ok && val != "" {
		q = q.Where("buah_raw.blok_panen = ?", val)
	}
	if val, ok := filter["jenis_durian_id"].(string); ok && val != "" {
		q = q.Where("buah_raw.jenis_durian = ?", val)
	}
	if val, ok := filter["is_sorted"]; ok {
		q = q.Where("buah_raw.is_sorted = ?", val)
	}

	return q
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
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&buah).
		Column("kode_buah").
		Where("kode_buah LIKE ?", fmt.Sprintf("%s-%%", kodeJenis)).
		Order("kode_buah DESC").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return "", err
	}
	return buah.KodeBuah, nil
}

func (r *buahRawRepository) GetNextSequenceWithLock(ctx context.Context, kodeJenis string) (int, error) {
	tx, err := r.db.InitQuery(ctx).Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var buah domain.BuahRaw
	err = tx.NewSelect().
		Model(&buah).
		Column("kode_buah").
		Where("kode_buah LIKE ?", fmt.Sprintf("%s-%%", kodeJenis)).
		Order("kode_buah DESC").
		Limit(1).
		For("UPDATE").
		Scan(ctx)

	sequence := 1
	if err == nil && buah.KodeBuah != "" {
		var lastSeq int
		_, err = fmt.Sscanf(buah.KodeBuah, kodeJenis+"-%d", &lastSeq)
		if err == nil {
			sequence = lastSeq + 1
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return sequence, nil
}

func (r *buahRawRepository) GetJenisDurianByID(ctx context.Context, id string) (domain.JenisDurian, error) {
	var jenis domain.JenisDurian
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&jenis).
		Where("id = ?", id).
		Scan(ctx)
	return jenis, err
}

func (r *buahRawRepository) GetJenisDurianByIDs(ctx context.Context, ids []string) (map[string]domain.JenisDurian, error) {
	var jenisList []domain.JenisDurian
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&jenisList).
		Where("id IN (?)", bun.In(ids)).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	result := make(map[string]domain.JenisDurian, len(jenisList))
	for _, j := range jenisList {
		result[j.ID] = j
	}
	return result, nil
}

func (r *buahRawRepository) GetBlokFullDetail(ctx context.Context, blokID string) (domain.Blok, domain.Divisi, domain.Estate, domain.Company, error) {
	var blok domain.Blok
	var divisi domain.Divisi
	var estate domain.Estate
	var company domain.Company

	err := r.db.InitQuery(ctx).NewSelect().
		Model(&blok).
		Where("id = ?", blokID).
		Scan(ctx)
	if err != nil {
		return blok, divisi, estate, company, err
	}

	err = r.db.InitQuery(ctx).NewSelect().
		Model(&divisi).
		Where("id = ?", blok.DivisiID).
		Scan(ctx)
	if err != nil {
		return blok, divisi, estate, company, err
	}

	err = r.db.InitQuery(ctx).NewSelect().
		Model(&estate).
		Where("id = ?", divisi.EstateID).
		Scan(ctx)
	if err != nil {
		return blok, divisi, estate, company, err
	}

	err = r.db.InitQuery(ctx).NewSelect().
		Model(&company).
		Where("id = ?", estate.CompanyID).
		Scan(ctx)

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