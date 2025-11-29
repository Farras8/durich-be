package response

type DashboardStokResponse struct {
	Summary     StokSummary           `json:"summary"`
	StokByJenis []StokByJenis         `json:"stok_by_jenis"`
	Throughput  ThroughputSummary     `json:"throughput"`
	Trend7Hari  []ThroughputTrendItem `json:"trend_7_hari"`
}

type StokSummary struct {
	TotalBuahMentah    int `json:"total_buah_mentah"`
	BuahBelumDisortir  int `json:"buah_belum_disortir"`
	TotalLotAktif      int `json:"total_lot_aktif"`
	LotReadyToShip     int `json:"lot_ready_to_ship"`
	LotEmpty           int `json:"lot_empty"`
}

type StokByJenis struct {
	JenisDurian           string            `json:"jenis_durian"`
	TotalQty              int               `json:"total_qty"`
	TotalBerat            float64           `json:"total_berat"`
	LotCount              int               `json:"lot_count"`
	AvgGradeDistribution  map[string]string `json:"avg_grade_distribution"`
}

type ThroughputSummary struct {
	BuahMasukHarian  float64 `json:"buah_masuk_harian"`
	BuahKeluarHarian float64 `json:"buah_keluar_harian"`
	LotSelesaiHarian float64 `json:"lot_selesai_harian"`
	PengirimanHarian float64 `json:"pengiriman_harian"`
}

type ThroughputTrendItem struct {
	Tanggal      string `json:"tanggal"`
	BuahMasuk    int    `json:"buah_masuk"`
	BuahKeluar   int    `json:"buah_keluar"`
	LotCreated   int    `json:"lot_created"`
	ShipmentSent int    `json:"shipment_sent"`
}

type DashboardSalesResponse struct {
	Summary          SalesSummary           `json:"summary"`
	BreakdownByJenis []SalesBreakdownJenis  `json:"breakdown_by_jenis"`
	BreakdownByTipe  []SalesBreakdownTipe   `json:"breakdown_by_tipe"`
	TrendHarga       []SalesTrendHarga      `json:"trend_harga"`
	TopBuyers        []SalesTopBuyer        `json:"top_buyers"`
}

type SalesSummary struct {
	TotalOmzet           float64 `json:"total_omzet"`
	TotalTransaksi       int     `json:"total_transaksi"`
	RataHargaPerKg       float64 `json:"rata_rata_harga_per_kg"`
	TotalBeratTerjual    float64 `json:"total_berat_terjual"`
	GrowthVsBulanLalu    string  `json:"growth_vs_bulan_lalu"`
}

type SalesBreakdownJenis struct {
	JenisDurian      string  `json:"jenis_durian"`
	Omzet            float64 `json:"omzet"`
	BeratTerjual     float64 `json:"berat_terjual"`
	RataHargaPerKg   float64 `json:"rata_harga_per_kg"`
	ShareOmzet       string  `json:"share_omzet"`
	TransaksiCount   int     `json:"transaksi_count"`
}

type SalesBreakdownTipe struct {
	TipeJual           string  `json:"tipe_jual"`
	Omzet              float64 `json:"omzet"`
	TransaksiCount     int     `json:"transaksi_count"`
	RataNilaiTransaksi float64 `json:"rata_nilai_transaksi"`
}

type SalesTrendHarga struct {
	Tanggal        string  `json:"tanggal"`
	JenisDurian    string  `json:"jenis_durian"`
	RataHargaPerKg float64 `json:"rata_harga_per_kg"`
}

type SalesTopBuyer struct {
	Tujuan             string  `json:"tujuan"`
	TotalPembelian     float64 `json:"total_pembelian"`
	Frekuensi          int     `json:"frekuensi"`
	RataPerTransaksi   float64 `json:"rata_per_transaksi"`
}

type WarehouseDataResponse struct {
	TotalBuahRawToday int `json:"total_buah_raw_today"`
	TotalLotReady     int `json:"total_lot_ready"`
	TotalLotSent      int `json:"total_lot_sent"`
}
