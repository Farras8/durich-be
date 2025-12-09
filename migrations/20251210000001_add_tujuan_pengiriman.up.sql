-- up
CREATE TABLE tb_tujuan_pengiriman (
    id VARCHAR(27) PRIMARY KEY,
    nama TEXT NOT NULL,
    tipe TEXT NOT NULL,
    alamat TEXT,
    kontak TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

ALTER TABLE tb_pengiriman ADD COLUMN tujuan_id VARCHAR(27);
ALTER TABLE tb_pengiriman ADD CONSTRAINT fk_pengiriman_tujuan FOREIGN KEY (tujuan_id) REFERENCES tb_tujuan_pengiriman(id);
