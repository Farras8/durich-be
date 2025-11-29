ALTER TABLE tb_stok_lot DROP CONSTRAINT IF EXISTS fk_stok_lot_jenis_durian;
ALTER TABLE tb_stok_lot ALTER COLUMN jenis_durian_id TYPE TEXT;
ALTER TABLE tb_stok_lot RENAME COLUMN jenis_durian_id TO jenis_durian;