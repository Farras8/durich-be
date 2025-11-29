CREATE TABLE estate (
    id VARCHAR(27) PRIMARY KEY,
    kode VARCHAR(5) UNIQUE NOT NULL,
    nama TEXT NOT NULL,
    company_id VARCHAR(27),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_estate_company FOREIGN KEY (company_id) REFERENCES company(id)
);

INSERT INTO estate (id, kode, nama, company_id) VALUES
('4SRlQ8zX9vJ2mN5P6Q7R8S9T001', 'ES01', 'Kebun Durian', '3SRlQ8zX9vJ2mN5P6Q7R8S9T001');
