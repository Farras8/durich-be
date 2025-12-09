ALTER TABLE users ADD COLUMN current_location_id VARCHAR(27);
ALTER TABLE users ADD CONSTRAINT fk_users_current_location FOREIGN KEY (current_location_id) REFERENCES tb_tujuan_pengiriman(id);