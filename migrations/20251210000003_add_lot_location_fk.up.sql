ALTER TABLE tb_stok_lot ADD CONSTRAINT fk_stok_lot_location FOREIGN KEY (current_location_id) REFERENCES tb_tujuan_pengiriman(id);
