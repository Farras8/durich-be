ALTER TABLE tb_buah_raw ALTER COLUMN jenis_durian TYPE VARCHAR(27);
ALTER TABLE tb_buah_raw ALTER COLUMN blok_panen TYPE VARCHAR(27);
ALTER TABLE tb_buah_raw ALTER COLUMN pohon_panen TYPE VARCHAR(27);
ALTER TABLE tb_buah_raw ADD CONSTRAINT fk_jenis_durian FOREIGN KEY (jenis_durian) REFERENCES jenis_durian(id);
ALTER TABLE tb_buah_raw ADD CONSTRAINT fk_blok_panen FOREIGN KEY (blok_panen) REFERENCES blok(id);
ALTER TABLE tb_buah_raw ADD CONSTRAINT fk_pohon_panen FOREIGN KEY (pohon_panen) REFERENCES pohon(id);