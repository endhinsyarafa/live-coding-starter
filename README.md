# live-coding-starter

Starter/base project fullstack untuk live coding interview dengan React + Vite, Go + Gin, GORM, dan PostgreSQL.

Project ini sengaja dibuat sederhana: sudah runnable, memiliki contoh CRUD `products`, dan mudah dijadikan blueprint ketika interviewer meminta entity baru.

## Stack

- Frontend: React, Vite, JavaScript, React Router, Axios, CSS biasa
- Backend: Go, Gin, GORM, PostgreSQL, godotenv
- Database: PostgreSQL 14 lokal
- Target development: macOS Apple Silicon / M1

## Structure

```text
live-coding-starter/
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── components/
│   │   ├── pages/products/
│   │   └── services/
│   ├── .env.example
│   ├── package.json
│   └── vite.config.js
├── backend/
│   ├── cmd/server/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── dto/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── router/
│   │   └── service/
│   ├── pkg/response/
│   ├── .env.example
│   └── go.mod
├── .gitignore
├── Makefile
└── README.md
```

## Database Setup

Buat database lokal:

```bash
createdb -h localhost -p 5434 -U macends live_coding_db
```

Test koneksi database:

```bash
psql -h localhost -p 5434 -U macends -d live_coding_db
```

Default konfigurasi database ada di `backend/.env.example`:

```env
DB_HOST=localhost
DB_PORT=5434
DB_USER=macends
DB_PASSWORD=
DB_NAME=live_coding_db
DB_SSLMODE=disable
```

## Backend Setup

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/server
```

Atau dari root project:

```bash
make backend
```

Backend berjalan di:

```text
http://localhost:8080
```

Health check:

```text
http://localhost:8080/health
```

## Frontend Setup

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Atau dari root project:

```bash
make frontend
```

Frontend berjalan di:

```text
http://localhost:5173
```

Untuk menjalankan backend dan frontend sekaligus dari root project:

```bash
make run
```

## Environment

Frontend:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Backend:

```env
APP_ENV=development
HTTP_PORT=8080

ALLOWED_ORIGIN=http://localhost:5173

JWT_SECRET=change-me-development-secret
JWT_EXPIRATION_HOURS=24

DB_HOST=localhost
DB_PORT=5434
DB_USER=macends
DB_PASSWORD=
DB_NAME=live_coding_db
DB_SSLMODE=disable

SEED_DATA=false
SEED_ADMIN=true
ADMIN_NAME=Administrator
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=Password123!
```

`SEED_DATA=true` akan membuat 3 sample products jika table masih kosong. Jika table sudah berisi data, seed tidak dijalankan lagi.

`SEED_ADMIN=true` akan membuat admin development jika email belum ada:

```text
Email: admin@example.com
Password: Password123!
```

## API Endpoints

```text
GET    /health
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/auth/me
GET    /api/v1/admin/dashboard
GET    /api/v1/admin/users
GET    /api/v1/admin/users/:id
POST   /api/v1/admin/users
PATCH  /api/v1/admin/users/:id
PATCH  /api/v1/admin/users/:id/status
GET    /api/v1/products
GET    /api/v1/products/:id
POST   /api/v1/products
PUT    /api/v1/products/:id
DELETE /api/v1/products/:id
```

## Admin Dashboard

Login admin:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Password123!"
  }'
```

Frontend pages:

```text
/login
/admin
/admin/users
```

Admin API menggunakan Bearer token:

```bash
curl http://localhost:8080/api/v1/admin/dashboard \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Quick Test

Health:

```bash
curl http://localhost:8080/health
```

Create:

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Keyboard Mechanical",
    "description": "Mechanical keyboard",
    "price": 750000,
    "stock": 10
  }'
```

List:

```bash
curl http://localhost:8080/api/v1/products
```

Detail:

```bash
curl http://localhost:8080/api/v1/products/1
```

Update:

```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Keyboard Mechanical Updated",
    "description": "Updated description",
    "price": 800000,
    "stock": 15
  }'
```

Delete:

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1
```

## Architecture

Frontend:

```text
Page -> Service -> Axios -> REST API
```

- Page berisi UI, state loading/error, dan event handler.
- Service berisi semua request API untuk entity tertentu.
- Axios instance membaca base URL dari environment variable.

Backend:

```text
Router -> Handler -> Service -> Repository -> PostgreSQL
```

- Router mendefinisikan endpoint dan middleware.
- Handler membaca request, melakukan binding/validation, dan menentukan HTTP status.
- Service berisi business logic sederhana dan mapping DTO ke model.
- Repository berisi query database menggunakan GORM.
- PostgreSQL menyimpan data, dengan AutoMigrate untuk mempercepat setup live coding.

## How to Add a New Entity

Misalnya interviewer meminta entity `employees`.

Backend checklist:

1. Tambah model di `backend/internal/model`.
2. Tambah DTO di `backend/internal/dto`.
3. Tambah repository di `backend/internal/repository`.
4. Tambah service di `backend/internal/service`.
5. Tambah handler di `backend/internal/handler`.
6. Register route di `backend/internal/router`.
7. Tambah model ke AutoMigrate di `backend/internal/database/postgres.go`.

Frontend checklist:

1. Tambah service di `frontend/src/services`.
2. Tambah page list.
3. Tambah form.
4. Tambah detail page jika diperlukan.
5. Tambah route di `frontend/src/App.jsx`.

## Verification

Command yang sebaiknya dijalankan sebelum interview:

```bash
cd backend
go mod tidy
go test ./...
go vet ./...
```

```bash
cd frontend
npm install
npm run build
```

Prinsip project:

```text
simple > clever
readable > abstract
ready for live coding > production complexity
```
