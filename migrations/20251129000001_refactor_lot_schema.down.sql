ALTER TABLE tb_buah_raw
DROP COLUMN IF EXISTS lot_id,
DROP COLUMN IF EXISTS berat,
DROP COLUMN IF EXISTS blok_id;

ALTER TABLE tb_buah_raw
ADD COLUMN is_sorted BOOLEAN DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS tb_lot_detail (
    lot_id VARCHAR(27) NOT NULL,
    buah_raw_id VARCHAR(27) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (lot_id, buah_raw_id),
    FOREIGN KEY (lot_id) REFERENCES tb_stok_lot(id) ON DELETE CASCADE,
    FOREIGN KEY (buah_raw_id) REFERENCES tb_buah_raw(id) ON DELETE CASCADE
);
