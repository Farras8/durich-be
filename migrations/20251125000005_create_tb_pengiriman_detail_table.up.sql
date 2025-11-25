CREATE TABLE tb_pengiriman_detail (
    id SERIAL PRIMARY KEY,
    pengiriman_id VARCHAR(27) NOT NULL,
    lot_sumber_id VARCHAR(27) NOT NULL,
    qty_ambil INT NOT NULL,
    berat_ambil DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_pengiriman FOREIGN KEY (pengiriman_id) REFERENCES tb_pengiriman(id) ON DELETE CASCADE,
    CONSTRAINT fk_lot FOREIGN KEY (lot_sumber_id) REFERENCES tb_stok_lot(id) ON DELETE RESTRICT
);