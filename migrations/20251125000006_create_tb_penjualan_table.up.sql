CREATE TABLE tb_penjualan (
    id VARCHAR(27) PRIMARY KEY,
    pengiriman_id VARCHAR(27) NOT NULL,
    berat_terjual DECIMAL(10,2),
    harga_total DECIMAL(15,2),
    tipe_jual TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_pengiriman FOREIGN KEY (pengiriman_id) REFERENCES tb_pengiriman(id)
);