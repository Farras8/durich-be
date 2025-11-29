package repository

import (
	"context"
	"durich-be/internal/dto/response"
	"durich-be/pkg/database"
	"sync"
	"time"
)

type DashboardRepository interface {
	GetStokDashboard(ctx context.Context, dateFrom, dateTo time.Time) (*response.DashboardStokResponse, error)
	GetSalesDashboard(ctx context.Context, dateFrom, dateTo time.Time) (*response.DashboardSalesResponse, error)
	GetWarehouseData(ctx context.Context) (*response.WarehouseDataResponse, error)
}

type dashboardRepository struct {
	db *database.Database
}

func NewDashboardRepository(db *database.Database) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetStokDashboard(ctx context.Context, dateFrom, dateTo time.Time) (*response.DashboardStokResponse, error) {
	var (
		summary     response.StokSummary
		stokByJenis []response.StokByJenis
		throughput  response.ThroughputSummary
		trend       []response.ThroughputTrendItem
	)

	var wg sync.WaitGroup
	errC := make(chan error, 4)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s, err := r.getStokSummary(ctx)
		if err != nil {
			errC <- err
			return
		}
		summary = s
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		j, err := r.getStokByJenis(ctx)
		if err != nil {
			errC <- err
			return
		}
		stokByJenis = j
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t, err := r.getThroughputSummary(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		throughput = t
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tr, err := r.getThroughputTrend(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		trend = tr
	}()

	wg.Wait()
	close(errC)

	for err := range errC {
		if err != nil {
			return nil, err
		}
	}

	return &response.DashboardStokResponse{
		Summary:     summary,
		StokByJenis: stokByJenis,
		Throughput:  throughput,
		Trend7Hari:  trend,
	}, nil
}

func (r *dashboardRepository) getStokSummary(ctx context.Context) (response.StokSummary, error) {
	var summary response.StokSummary

	count, err := r.db.NewSelect().
		Model((*struct{ TB_BUAH_RAW int })(nil)).
		ColumnExpr("COUNT(*) as total").
		Table("tb_buah_raw").
		Where("deleted_at IS NULL").
		Count(ctx)
	if err != nil {
		return summary, err
	}
	summary.TotalBuahMentah = count

	count, err = r.db.NewSelect().
		ColumnExpr("COUNT(*) as total").
		Table("tb_buah_raw").
		Where("deleted_at IS NULL AND is_sorted = ?", false).
		Count(ctx)
	if err != nil {
		return summary, err
	}
	summary.BuahBelumDisortir = count

	count, err = r.db.NewSelect().
		ColumnExpr("COUNT(*) as total").
		Table("tb_stok_lot").
		Where("deleted_at IS NULL AND status != ?", "EMPTY").
		Count(ctx)
	if err != nil {
		return summary, err
	}
	summary.TotalLotAktif = count

	count, err = r.db.NewSelect().
		ColumnExpr("COUNT(*) as total").
		Table("tb_stok_lot").
		Where("deleted_at IS NULL AND status = ?", "READY").
		Count(ctx)
	if err != nil {
		return summary, err
	}
	summary.LotReadyToShip = count

	count, err = r.db.NewSelect().
		ColumnExpr("COUNT(*) as total").
		Table("tb_stok_lot").
		Where("deleted_at IS NULL AND status = ?", "EMPTY").
		Count(ctx)
	if err != nil {
		return summary, err
	}
	summary.LotEmpty = count

	return summary, nil
}

func (r *dashboardRepository) getStokByJenis(ctx context.Context) ([]response.StokByJenis, error) {
	type queryResult struct {
		JenisDurian string  `bun:"nama_jenis"`
		TotalQty    int     `bun:"total_qty"`
		TotalBerat  float64 `bun:"total_berat"`
		LotCount    int     `bun:"lot_count"`
	}

	var results []queryResult

	err := r.db.NewSelect().
		ColumnExpr("jd.nama_jenis").
		ColumnExpr("SUM(sl.qty_sisa) as total_qty").
		ColumnExpr("SUM(sl.berat_sisa) as total_berat").
		ColumnExpr("COUNT(sl.id) as lot_count").
		TableExpr("tb_stok_lot AS sl").
		Join("LEFT JOIN jenis_durian AS jd ON sl.jenis_durian = jd.id").
		Where("sl.deleted_at IS NULL").
		Where("sl.status != ?", "EMPTY").
		Group("jd.nama_jenis").
		Order("total_qty DESC").
		Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	stokByJenis := make([]response.StokByJenis, 0, len(results))
	for _, r := range results {
		stokByJenis = append(stokByJenis, response.StokByJenis{
			JenisDurian:          r.JenisDurian,
			TotalQty:             r.TotalQty,
			TotalBerat:           r.TotalBerat,
			LotCount:             r.LotCount,
			AvgGradeDistribution: map[string]string{},
		})
	}

	return stokByJenis, nil
}

func (r *dashboardRepository) getThroughputSummary(ctx context.Context, dateFrom, dateTo time.Time) (response.ThroughputSummary, error) {
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	if days <= 0 {
		days = 1
	}

	type queryResult struct {
		BuahMasuk  int `bun:"buah_masuk"`
		LotSelesai int `bun:"lot_selesai"`
		Pengiriman int `bun:"pengiriman"`
	}

	var result queryResult

	subquery1 := r.db.NewSelect().
		ColumnExpr("COUNT(br.id)").
		TableExpr("tb_buah_raw AS br").
		Where("br.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("br.deleted_at IS NULL")

	subquery2 := r.db.NewSelect().
		ColumnExpr("COUNT(sl.id)").
		TableExpr("tb_stok_lot AS sl").
		Where("sl.updated_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("sl.status = ?", "READY").
		Where("sl.deleted_at IS NULL")

	subquery3 := r.db.NewSelect().
		ColumnExpr("COUNT(p.id)").
		TableExpr("tb_pengiriman AS p").
		Where("p.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("p.deleted_at IS NULL")

	err := r.db.NewSelect().
		ColumnExpr("(?) as buah_masuk", subquery1).
		ColumnExpr("(?) as lot_selesai", subquery2).
		ColumnExpr("(?) as pengiriman", subquery3).
		Scan(ctx, &result)
	if err != nil {
		return response.ThroughputSummary{}, err
	}

	return response.ThroughputSummary{
		BuahMasukHarian:  float64(result.BuahMasuk) / float64(days),
		BuahKeluarHarian: float64(result.BuahMasuk) / float64(days),
		LotSelesaiHarian: float64(result.LotSelesai) / float64(days),
		PengirimanHarian: float64(result.Pengiriman) / float64(days),
	}, nil
}

func (r *dashboardRepository) getThroughputTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]response.ThroughputTrendItem, error) {
	type queryResult struct {
		Tanggal      string `bun:"tanggal"`
		BuahMasuk    int    `bun:"buah_masuk"`
		LotCreated   int    `bun:"lot_created"`
		ShipmentSent int    `bun:"shipment_sent"`
	}

	var results []queryResult

	err := r.db.NewSelect().
		ColumnExpr("TO_CHAR(created_at, 'YYYY-MM-DD') as tanggal").
		ColumnExpr("COUNT(*) as buah_masuk").
		ColumnExpr("0 as lot_created").
		ColumnExpr("0 as shipment_sent").
		Table("tb_buah_raw").
		Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("deleted_at IS NULL").
		GroupExpr("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("tanggal DESC").
		Limit(7).
		Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	trends := make([]response.ThroughputTrendItem, 0, len(results))
	for _, r := range results {
		trends = append(trends, response.ThroughputTrendItem{
			Tanggal:      r.Tanggal,
			BuahMasuk:    r.BuahMasuk,
			BuahKeluar:   r.BuahMasuk,
			LotCreated:   r.LotCreated,
			ShipmentSent: r.ShipmentSent,
		})
	}

	return trends, nil
}

func (r *dashboardRepository) GetSalesDashboard(ctx context.Context, dateFrom, dateTo time.Time) (*response.DashboardSalesResponse, error) {
	var (
		summary       response.SalesSummary
		breakdownJens []response.SalesBreakdownJenis
		breakdownTipe []response.SalesBreakdownTipe
		trendHarga    []response.SalesTrendHarga
		topBuyers     []response.SalesTopBuyer
	)

	var wg sync.WaitGroup
	errC := make(chan error, 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		s, err := r.getSalesSummary(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		summary = s
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		j, err := r.getSalesBreakdownJenis(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		breakdownJens = j
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t, err := r.getSalesBreakdownTipe(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		breakdownTipe = t
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		h, err := r.getSalesTrendHarga(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		trendHarga = h
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		b, err := r.getSalesTopBuyers(ctx, dateFrom, dateTo)
		if err != nil {
			errC <- err
			return
		}
		topBuyers = b
	}()

	wg.Wait()
	close(errC)

	for err := range errC {
		if err != nil {
			return nil, err
		}
	}

	return &response.DashboardSalesResponse{
		Summary:          summary,
		BreakdownByJenis: breakdownJens,
		BreakdownByTipe:  breakdownTipe,
		TrendHarga:       trendHarga,
		TopBuyers:        topBuyers,
	}, nil
}

func (r *dashboardRepository) getSalesSummary(ctx context.Context, dateFrom, dateTo time.Time) (response.SalesSummary, error) {
	type queryResult struct {
		TotalOmzet        float64 `bun:"total_omzet"`
		TotalTransaksi    int     `bun:"total_transaksi"`
		TotalBeratTerjual float64 `bun:"total_berat_terjual"`
	}

	var result queryResult

	err := r.db.NewSelect().
		ColumnExpr("COALESCE(SUM(harga_total), 0) as total_omzet").
		ColumnExpr("COUNT(*) as total_transaksi").
		ColumnExpr("COALESCE(SUM(berat_terjual), 0) as total_berat_terjual").
		Table("tb_penjualan").
		Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("deleted_at IS NULL").
		Scan(ctx, &result)
	if err != nil {
		return response.SalesSummary{}, err
	}

	rataHargaPerKg := float64(0)
	if result.TotalBeratTerjual > 0 {
		rataHargaPerKg = result.TotalOmzet / result.TotalBeratTerjual
	}

	return response.SalesSummary{
		TotalOmzet:        result.TotalOmzet,
		TotalTransaksi:    result.TotalTransaksi,
		RataHargaPerKg:    rataHargaPerKg,
		TotalBeratTerjual: result.TotalBeratTerjual,
		GrowthVsBulanLalu: "+0%",
	}, nil
}

func (r *dashboardRepository) getSalesBreakdownJenis(ctx context.Context, dateFrom, dateTo time.Time) ([]response.SalesBreakdownJenis, error) {
	type queryResult struct {
		JenisDurian    string  `bun:"jenis_durian"`
		Omzet          float64 `bun:"omzet"`
		BeratTerjual   float64 `bun:"berat_terjual"`
		TransaksiCount int     `bun:"transaksi_count"`
	}

	var results []queryResult

	err := r.db.NewSelect().
		ColumnExpr("COALESCE(jd.nama_jenis, 'Unknown') as jenis_durian").
		ColumnExpr("SUM(p.harga_total) as omzet").
		ColumnExpr("SUM(p.berat_terjual) as berat_terjual").
		ColumnExpr("COUNT(p.id) as transaksi_count").
		TableExpr("tb_penjualan AS p").
		Join("LEFT JOIN tb_pengiriman AS pr ON p.pengiriman_id = pr.id").
		Join("LEFT JOIN tb_pengiriman_detail AS pd ON pr.id = pd.pengiriman_id").
		Join("LEFT JOIN tb_stok_lot AS sl ON pd.lot_sumber_id = sl.id").
		Join("LEFT JOIN jenis_durian AS jd ON sl.jenis_durian = jd.id").
		Where("p.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("p.deleted_at IS NULL").
		Group("jd.nama_jenis").
		Order("omzet DESC").
		Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	breakdown := make([]response.SalesBreakdownJenis, 0, len(results))
	for _, r := range results {
		rataHargaPerKg := float64(0)
		if r.BeratTerjual > 0 {
			rataHargaPerKg = r.Omzet / r.BeratTerjual
		}

		breakdown = append(breakdown, response.SalesBreakdownJenis{
			JenisDurian:    r.JenisDurian,
			Omzet:          r.Omzet,
			BeratTerjual:   r.BeratTerjual,
			RataHargaPerKg: rataHargaPerKg,
			ShareOmzet:     "0%",
			TransaksiCount: r.TransaksiCount,
		})
	}

	return breakdown, nil
}

func (r *dashboardRepository) getSalesBreakdownTipe(ctx context.Context, dateFrom, dateTo time.Time) ([]response.SalesBreakdownTipe, error) {
	type queryResult struct {
		TipeJual       string  `bun:"tipe_jual"`
		Omzet          float64 `bun:"omzet"`
		TransaksiCount int     `bun:"transaksi_count"`
	}

	var results []queryResult

	err := r.db.NewSelect().
		Column("tipe_jual").
		ColumnExpr("SUM(harga_total) as omzet").
		ColumnExpr("COUNT(*) as transaksi_count").
		Table("tb_penjualan").
		Where("created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("deleted_at IS NULL").
		Group("tipe_jual").
		Order("omzet DESC").
		Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	breakdown := make([]response.SalesBreakdownTipe, 0, len(results))
	for _, r := range results {
		rataNilaiTransaksi := float64(0)
		if r.TransaksiCount > 0 {
			rataNilaiTransaksi = r.Omzet / float64(r.TransaksiCount)
		}

		breakdown = append(breakdown, response.SalesBreakdownTipe{
			TipeJual:           r.TipeJual,
			Omzet:              r.Omzet,
			TransaksiCount:     r.TransaksiCount,
			RataNilaiTransaksi: rataNilaiTransaksi,
		})
	}

	return breakdown, nil
}

func (r *dashboardRepository) getSalesTrendHarga(ctx context.Context, dateFrom, dateTo time.Time) ([]response.SalesTrendHarga, error) {
	return []response.SalesTrendHarga{}, nil
}

func (r *dashboardRepository) getSalesTopBuyers(ctx context.Context, dateFrom, dateTo time.Time) ([]response.SalesTopBuyer, error) {
	type queryResult struct {
		Tujuan         string  `bun:"tujuan"`
		TotalPembelian float64 `bun:"total_pembelian"`
		Frekuensi      int     `bun:"frekuensi"`
	}

	var results []queryResult

	err := r.db.NewSelect().
		Column("pr.tujuan").
		ColumnExpr("SUM(p.harga_total) as total_pembelian").
		ColumnExpr("COUNT(p.id) as frekuensi").
		TableExpr("tb_penjualan AS p").
		Join("LEFT JOIN tb_pengiriman AS pr ON p.pengiriman_id = pr.id").
		Where("p.created_at BETWEEN ? AND ?", dateFrom, dateTo).
		Where("p.deleted_at IS NULL").
		Group("pr.tujuan").
		Order("total_pembelian DESC").
		Limit(10).
		Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	topBuyers := make([]response.SalesTopBuyer, 0, len(results))
	for _, r := range results {
		rataPerTransaksi := float64(0)
		if r.Frekuensi > 0 {
			rataPerTransaksi = r.TotalPembelian / float64(r.Frekuensi)
		}

		topBuyers = append(topBuyers, response.SalesTopBuyer{
			Tujuan:           r.Tujuan,
			TotalPembelian:   r.TotalPembelian,
			Frekuensi:        r.Frekuensi,
			RataPerTransaksi: rataPerTransaksi,
		})
	}

	return topBuyers, nil
}

func (r *dashboardRepository) GetWarehouseData(ctx context.Context) (*response.WarehouseDataResponse, error) {
	var (
		totalBuahRawToday int
		totalLotReady     int
		totalLotSent      int
	)

	today := time.Now().Format("2006-01-02")

	var wg sync.WaitGroup
	errC := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := r.db.NewSelect().
			Table("tb_buah_raw").
			Where("tgl_panen = ?", today).
			Where("deleted_at IS NULL").
			Count(ctx)
		if err != nil {
			errC <- err
			return
		}
		totalBuahRawToday = count
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := r.db.NewSelect().
			Table("tb_stok_lot").
			Where("status = ?", "READY").
			Where("deleted_at IS NULL").
			Count(ctx)
		if err != nil {
			errC <- err
			return
		}
		totalLotReady = count
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Count unique lots that are in shipments created today (or overall sent, depends on req)
		// Request says "total lot yang statusnya sent" in the context of "warehouse-data" which usually implies current snapshot or today's activity.
		// "total buah raw HARI INI", so likely "total lot sent HARI INI".
		// But the field name is total_lot_sent, while total_buah_raw_today has "today".
		// Let's assume "Today" for consistency with the first metric, OR "Total Sent ever"?
		// Usually dashboard widgets show:
		// 1. Buah Raw Today (Daily input)
		// 2. Stock Ready (Current Inventory Snapshot)
		// 3. Stock Sent (Today's Output) -> This makes the most sense for a daily dashboard flow.
		
		// Query: Count distinct lot_sumber_id from tb_pengiriman_detail 
		// JOIN tb_pengiriman ON detail.pengiriman_id = pengiriman.id
		// WHERE pengiriman.created_at = today AND pengiriman.deleted_at IS NULL
		
		count, err := r.db.NewSelect().
			ColumnExpr("COUNT(DISTINCT pd.lot_sumber_id)").
			TableExpr("tb_pengiriman_detail AS pd").
			Join("JOIN tb_pengiriman AS p ON pd.pengiriman_id = p.id").
			Where("DATE(p.created_at) = ?", today).
			Where("p.deleted_at IS NULL").
			Count(ctx)
			
		if err != nil {
			errC <- err
			return
		}
		totalLotSent = count
	}()

	wg.Wait()
	close(errC)

	for err := range errC {
		if err != nil {
			return nil, err
		}
	}

	return &response.WarehouseDataResponse{
		TotalBuahRawToday: totalBuahRawToday,
		TotalLotReady:     totalLotReady,
		TotalLotSent:      totalLotSent,
	}, nil
}
