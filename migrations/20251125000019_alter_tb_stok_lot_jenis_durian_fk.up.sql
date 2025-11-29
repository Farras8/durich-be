ALTER TABLE tb_stok_lot RENAME COLUMN jenis_durian TO jenis_durian_id;
ALTER TABLE tb_stok_lot ALTER COLUMN jenis_durian_id TYPE VARCHAR(27);
ALTER TABLE tb_stok_lot ADD CONSTRAINT fk_stok_lot_jenis_durian FOREIGN KEY (jenis_durian_id) REFERENCES jenis_durian(id);