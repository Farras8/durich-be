# List Endpoints API Durich

## Authentication & User Management

### Admin Registration
- `POST /v1/admin/register-admin` - Public (no auth)
- `POST /v1/admin/register-warehouse` - Public (no auth)
- `POST /v1/admin/register-sales` - Public (no auth)
- `POST /v1/admin/users/reset-password` - Admin

### Authentication
- `POST /v1/authentications/login` - Public (no auth)
- `POST /v1/authentications/refresh-token` - Public (no auth)
- `POST /v1/authentications/logout` - Admin, Warehouse, Sales (authenticated)

### Profile
- `PUT /v1/profile/password` - Admin, Warehouse, Sales (authenticated)

## Buah Raw (Raw Fruit)
- `POST /v1/buah-raw` - Admin, Warehouse
- `POST /v1/buah-raw/bulk` - Admin, Warehouse
- `GET /v1/buah-raw` - Admin, Warehouse
- `GET /v1/buah-raw/:id` - Admin, Warehouse
- `PUT /v1/buah-raw/:id` - Admin, Warehouse
- `DELETE /v1/buah-raw/:id` - Admin, Warehouse
- `GET /v1/buah-raw/unsorted` - Admin, Warehouse

## Lots
- `POST /v1/lots` - Admin, Warehouse
- `GET /v1/lots` - Admin, Warehouse
- `GET /v1/lots/:id` - Admin, Warehouse
- `POST /v1/lots/:id/items` - Admin, Warehouse
- `DELETE /v1/lots/:id/items` - Admin, Warehouse
- `POST /v1/lots/:id/finalize` - Admin, Warehouse

## Shipments
- `POST /v1/shipments` - Admin, Warehouse
- `GET /v1/shipments` - Admin, Warehouse
- `GET /v1/shipments/:id` - Admin, Warehouse
- `POST /v1/shipments/:id/items` - Admin, Warehouse
- `DELETE /v1/shipments/:id/items` - Admin, Warehouse
- `POST /v1/shipments/:id/finalize` - Admin, Warehouse
- `PATCH /v1/shipments/:id/status` - Admin, Sales

## Sales
- `POST /v1/sales` - Admin, Sales
- `GET /v1/sales` - Admin, Sales
- `GET /v1/sales/:id` - Admin, Sales
- `PUT /v1/sales/:id` - Admin, Sales
- `DELETE /v1/sales/:id` - Admin, Sales

## Dashboard
- `GET /v1/dashboard/stok` - Admin, Warehouse
- `GET /v1/dashboard/sales` - Admin, Sales

## Traceability
- `GET /v1/trace/lot/:id` - Admin, Warehouse, Sales
- `GET /v1/trace/fruit/:buah_raw_id` - Admin, Warehouse, Sales
- `GET /v1/trace/shipment/:id` - Admin, Warehouse, Sales

## Master Data

### Companies
- `POST /v1/companies/` - Admin
- `GET /v1/companies/` - Admin, Warehouse
- `GET /v1/companies/:id` - Admin, Warehouse
- `PUT /v1/companies/:id` - Admin
- `DELETE /v1/companies/:id` - Admin

### Estates
- `POST /v1/estates/` - Admin
- `GET /v1/estates/` - Admin, Warehouse
- `GET /v1/estates/:id` - Admin, Warehouse
- `PUT /v1/estates/:id` - Admin
- `DELETE /v1/estates/:id` - Admin

### Divisi
- `POST /v1/divisi/` - Admin
- `GET /v1/divisi/` - Admin, Warehouse
- `GET /v1/divisi/:id` - Admin, Warehouse
- `PUT /v1/divisi/:id` - Admin
- `DELETE /v1/divisi/:id` - Admin

### Bloks
- `POST /v1/bloks/` - Admin
- `GET /v1/bloks/` - Admin, Warehouse
- `GET /v1/bloks/:id` - Admin, Warehouse
- `PUT /v1/bloks/:id` - Admin
- `DELETE /v1/bloks/:id` - Admin

### Jenis Durian
- `POST /v1/jenis-durian/` - Admin
- `GET /v1/jenis-durian/` - Admin, Warehouse
- `GET /v1/jenis-durian/:id` - Admin, Warehouse
- `PUT /v1/jenis-durian/:id` - Admin
- `DELETE /v1/jenis-durian/:id` - Admin

### Pohon
- `POST /v1/pohon/` - Admin
- `GET /v1/pohon/` - Admin, Warehouse
- `GET /v1/pohon/:id` - Admin, Warehouse
- `PUT /v1/pohon/:id` - Admin
- `DELETE /v1/pohon/:id` - Admin

TOTAL ENDPOINTS: 66