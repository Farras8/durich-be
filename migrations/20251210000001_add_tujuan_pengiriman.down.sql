ALTER TABLE tb_pengiriman DROP CONSTRAINT fk_pengiriman_tujuan;
ALTER TABLE tb_pengiriman DROP COLUMN tujuan_id;
DROP TABLE tb_tujuan_pengiriman;