package repository

import (
	"context"
	"durich-be/internal/dto/response"
	"durich-be/internal/domain"
	"durich-be/pkg/database"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type TraceabilityRepository interface {
	TraceLot(ctx context.Context, lotID string) (*response.TraceLotResponse, error)
	TraceFruit(ctx context.Context, fruitID string) (*response.TraceFruitResponse, error)
	TraceShipment(ctx context.Context, shipmentID string) (*response.TraceShipmentResponse, error)
}

type traceabilityRepository struct {
	db *database.Database
}

func NewTraceabilityRepository(db *database.Database) TraceabilityRepository {
	return &traceabilityRepository{db: db}
}

func (r *traceabilityRepository) TraceLot(ctx context.Context, lotID string) (*response.TraceLotResponse, error) {
	lot := new(domain.StokLot)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(lot).
		Relation("JenisDurianDetail").
		Where("stok_lot.id = ?", lotID).
		Where("stok_lot.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var fruits []domain.BuahRaw
	err = r.db.InitQuery(ctx).NewSelect().
		Model(&fruits).
		Join("INNER JOIN tb_lot_detail ON tb_lot_detail.buah_raw_id = buah_raw.id").
		Relation("BlokPanenDetail").
		Relation("BlokPanenDetail.Divisi").
		Relation("BlokPanenDetail.Divisi.Estate").
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Relation("PohonPanenDetail").
		Where("tb_lot_detail.lot_id = ?", lotID).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	var lotDetails []domain.LotDetail
	err = r.db.InitQuery(ctx).NewSelect().
		Model(&lotDetails).
		Where("lot_id = ?", lotID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	lotDetailMap := make(map[string]time.Time, len(lotDetails))
	for _, detail := range lotDetails {
		lotDetailMap[detail.BuahRawID] = detail.CreatedAt
	}

	jenisDurian := ""
	if lot.JenisDurianDetail != nil {
		jenisDurian = lot.JenisDurianDetail.NamaJenis
	}

	lotInfo := response.LotTraceInfo{
		ID:          lot.ID,
		JenisDurian: jenisDurian,
		KondisiBuah: lot.KondisiBuah,
		QtyTotal:    lot.QtyAwal,
		BeratTotal:  lot.BeratAwal,
		Status:      lot.Status,
		CreatedAt:   lot.CreatedAt,
	}

	fruitInfos := make([]response.FruitTraceInfo, 0, len(fruits))
	for _, fruit := range fruits {
		addedAt := lot.CreatedAt
		if detailTime, exists := lotDetailMap[fruit.ID]; exists {
			addedAt = detailTime
		}

		lokasi := r.buildLokasiInfo(fruit)
		pohon := ""
		if fruit.PohonPanenDetail != nil {
			pohon = fruit.PohonPanenDetail.Kode
		}
		lokasi.Pohon = pohon

		fruitInfos = append(fruitInfos, response.FruitTraceInfo{
			BuahRawID:    fruit.ID,
			KodeBuah:     fruit.KodeBuah,
			TglPanen:     fruit.TglPanen,
			LokasiPanen:  lokasi,
			AddedToLotAt: addedAt,
		})
	}

	return &response.TraceLotResponse{
		LotInfo: lotInfo,
		Fruits:  fruitInfos,
	}, nil
}

func (r *traceabilityRepository) TraceFruit(ctx context.Context, fruitID string) (*response.TraceFruitResponse, error) {
	fruit := new(domain.BuahRaw)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(fruit).
		Relation("JenisDurianDetail").
		Relation("BlokPanenDetail").
		Relation("BlokPanenDetail.Divisi").
		Relation("BlokPanenDetail.Divisi.Estate").
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Relation("PohonPanenDetail").
		Where("buah_raw.id = ?", fruitID).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	jenisDurian := ""
	if fruit.JenisDurianDetail != nil {
		jenisDurian = fruit.JenisDurianDetail.NamaJenis
	}

	lokasi := r.buildLokasiInfo(*fruit)
	pohon := ""
	if fruit.PohonPanenDetail != nil {
		pohon = fruit.PohonPanenDetail.Kode
	}
	lokasi.Pohon = pohon

	fruitInfo := response.FruitDetailInfo{
		ID:          fruit.ID,
		KodeBuah:    fruit.KodeBuah,
		JenisDurian: jenisDurian,
		TglPanen:    fruit.TglPanen,
		LokasiPanen: lokasi,
	}

	journey := r.buildFruitJourney(ctx, fruitID)

	return &response.TraceFruitResponse{
		FruitInfo: fruitInfo,
		Journey:   journey,
	}, nil
}

func (r *traceabilityRepository) TraceShipment(ctx context.Context, shipmentID string) (*response.TraceShipmentResponse, error) {
	shipment := new(domain.Pengiriman)
	err := r.db.InitQuery(ctx).NewSelect().
		Model(shipment).
		Relation("Details").
		Where("p.id = ?", shipmentID).
		Where("p.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	totalQty := 0
	totalBerat := float64(0)
	lotIDs := make([]string, 0, len(shipment.Details))
	
	for _, detail := range shipment.Details {
		totalQty += detail.QtyAmbil
		totalBerat += detail.BeratAmbil
		lotIDs = append(lotIDs, detail.LotSumberID)
	}

	shipmentInfo := response.ShipmentTraceInfo{
		ID:         shipment.ID,
		Tujuan:     shipment.Tujuan,
		TglKirim:   shipment.TglKirim,
		Status:     shipment.Status,
		TotalQty:   totalQty,
		TotalBerat: totalBerat,
	}

	if len(lotIDs) == 0 {
		return &response.TraceShipmentResponse{
			ShipmentInfo:      shipmentInfo,
			BreakdownByLokasi: []response.BreakdownByLokasi{},
			DetailedFruits:    []response.DetailedFruitInfo{},
		}, nil
	}

	var fruits []domain.BuahRaw
	err = r.db.InitQuery(ctx).NewSelect().
		Model(&fruits).
		Join("INNER JOIN tb_lot_detail ON tb_lot_detail.buah_raw_id = buah_raw.id").
		Relation("BlokPanenDetail").
		Relation("BlokPanenDetail.Divisi").
		Relation("BlokPanenDetail.Divisi.Estate").
		Relation("BlokPanenDetail.Divisi.Estate.Company").
		Relation("PohonPanenDetail").
		Where("tb_lot_detail.lot_id IN (?)", bun.In(lotIDs)).
		Where("buah_raw.deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	detailedFruits := make([]response.DetailedFruitInfo, 0, len(fruits))
	for _, fruit := range fruits {
		lokasi := r.buildLokasiString(fruit)
		pohon := ""
		if fruit.PohonPanenDetail != nil {
			pohon = fruit.PohonPanenDetail.Kode
		}

		detailedFruits = append(detailedFruits, response.DetailedFruitInfo{
			BuahRawID: fruit.ID,
			KodeBuah:  fruit.KodeBuah,
			Lokasi:    lokasi,
			Pohon:     pohon,
			TglPanen:  fruit.TglPanen,
		})
	}

	return &response.TraceShipmentResponse{
		ShipmentInfo:      shipmentInfo,
		BreakdownByLokasi: []response.BreakdownByLokasi{},
		DetailedFruits:    detailedFruits,
	}, nil
}

func (r *traceabilityRepository) buildLokasiInfo(fruit domain.BuahRaw) response.LokasiTraceInfo {
	lokasi := response.LokasiTraceInfo{}
	
	if fruit.BlokPanenDetail == nil {
		return lokasi
	}

	lokasi.Blok = fruit.BlokPanenDetail.Kode
	
	if fruit.BlokPanenDetail.Divisi == nil {
		return lokasi
	}

	lokasi.Divisi = fruit.BlokPanenDetail.Divisi.Kode
	
	if fruit.BlokPanenDetail.Divisi.Estate == nil {
		return lokasi
	}

	lokasi.Estate = fruit.BlokPanenDetail.Divisi.Estate.Kode
	
	if fruit.BlokPanenDetail.Divisi.Estate.Company != nil {
		lokasi.Company = fruit.BlokPanenDetail.Divisi.Estate.Company.Kode
	}

	return lokasi
}

func (r *traceabilityRepository) buildLokasiString(fruit domain.BuahRaw) string {
	if fruit.BlokPanenDetail == nil ||
		fruit.BlokPanenDetail.Divisi == nil ||
		fruit.BlokPanenDetail.Divisi.Estate == nil ||
		fruit.BlokPanenDetail.Divisi.Estate.Company == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s-%s-%s",
		fruit.BlokPanenDetail.Divisi.Estate.Company.Kode,
		fruit.BlokPanenDetail.Divisi.Estate.Kode,
		fruit.BlokPanenDetail.Divisi.Kode,
		fruit.BlokPanenDetail.Kode)
}

func (r *traceabilityRepository) buildFruitJourney(ctx context.Context, fruitID string) response.FruitJourney {
	journey := response.FruitJourney{}

	var lotDetail domain.LotDetail
	err := r.db.InitQuery(ctx).NewSelect().
		Model(&lotDetail).
		Relation("Lot").
		Where("buah_raw_id = ?", fruitID).
		Scan(ctx)
	
	if err != nil || lotDetail.LotID == "" {
		return journey
	}

	lot := lotDetail.Lot
	if lot == nil {
		return journey
	}

	journey.Lot = &response.LotJourneyInfo{
		ID:          lot.ID,
		KondisiBuah: lot.KondisiBuah,
		JoinedAt:    lotDetail.CreatedAt,
		Status:      lot.Status,
	}

	var pengirimanDetail domain.PengirimanDetail
	err = r.db.InitQuery(ctx).NewSelect().
		Model(&pengirimanDetail).
		Relation("Pengiriman").
		Join("INNER JOIN tb_pengiriman p ON p.id = pd.pengiriman_id").
		Where("pd.lot_sumber_id = ?", lot.ID).
		Scan(ctx)
	
	if err != nil || pengirimanDetail.ID == "" {
		return journey
	}

	pengiriman := pengirimanDetail.Pengiriman
	if pengiriman == nil {
		return journey
	}

	journey.Shipment = &response.ShipmentJourneyInfo{
		ID:       pengiriman.ID,
		Tujuan:   pengiriman.Tujuan,
		TglKirim: pengiriman.TglKirim,
		Status:   pengiriman.Status,
	}

	var penjualan domain.Penjualan
	err = r.db.InitQuery(ctx).NewSelect().
		Model(&penjualan).
		Where("pengiriman_id = ?", pengiriman.ID).
		Where("deleted_at IS NULL").
		Scan(ctx)
	
	if err == nil && penjualan.ID != "" {
		journey.Sales = &response.SalesJourneyInfo{
			ID:           penjualan.ID,
			HargaTotal:   penjualan.HargaTotal,
			TipeJual:     penjualan.TipeJual,
			TglTransaksi: penjualan.CreatedAt,
		}
	}

	return journey
}