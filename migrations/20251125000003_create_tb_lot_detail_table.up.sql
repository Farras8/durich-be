CREATE TABLE tb_lot_detail (
    lot_id VARCHAR(27) NOT NULL,
    buah_raw_id VARCHAR(27) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (lot_id, buah_raw_id),
    CONSTRAINT fk_lot FOREIGN KEY (lot_id) REFERENCES tb_stok_lot(id) ON DELETE CASCADE,
    CONSTRAINT fk_buah FOREIGN KEY (buah_raw_id) REFERENCES tb_buah_raw(id) ON DELETE CASCADE
);