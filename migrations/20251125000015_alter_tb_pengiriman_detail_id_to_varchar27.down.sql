ALTER TABLE tb_pengiriman_detail ALTER COLUMN id TYPE INTEGER USING id::integer;
ALTER TABLE tb_pengiriman_detail ALTER COLUMN id SET DEFAULT nextval('tb_pengiriman_detail_id_seq');