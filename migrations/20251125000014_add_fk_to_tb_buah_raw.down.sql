ALTER TABLE tb_buah_raw DROP CONSTRAINT fk_jenis_durian;
ALTER TABLE tb_buah_raw DROP CONSTRAINT fk_blok_panen;
ALTER TABLE tb_buah_raw DROP CONSTRAINT fk_pohon_panen;
ALTER TABLE tb_buah_raw ALTER COLUMN jenis_durian TYPE TEXT;
ALTER TABLE tb_buah_raw ALTER COLUMN blok_panen TYPE VARCHAR(3);
ALTER TABLE tb_buah_raw ALTER COLUMN pohon_panen TYPE VARCHAR(10);