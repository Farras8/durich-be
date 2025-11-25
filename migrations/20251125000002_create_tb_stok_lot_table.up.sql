CREATE TABLE tb_stok_lot (
    id VARCHAR(27) PRIMARY KEY,
    jenis_durian TEXT NOT NULL,
    kondisi_buah TEXT NOT NULL,
    berat_awal DECIMAL(10,2) DEFAULT 0,
    qty_awal INT DEFAULT 0,
    berat_sisa DECIMAL(10,2) DEFAULT 0,
    qty_sisa INT DEFAULT 0,
    status TEXT DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);