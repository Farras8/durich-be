CREATE TABLE tb_buah_raw (
    id VARCHAR(27) PRIMARY KEY,
    jenis_durian TEXT NOT NULL,
    blok_panen VARCHAR(3) NOT NULL,
    pohon_panen VARCHAR(10),
    tgl_panen DATE NOT NULL,
    is_sorted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ
);