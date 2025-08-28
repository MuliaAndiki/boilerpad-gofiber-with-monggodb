# Boilerpad-GoFiber 🚀

[![Go](https://img.shields.io/badge/Go-1.24-blue?logo=go&logoColor=white)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-2.50.0-brightgreen)](https://gofiber.io/)
[![MongoDB](https://img.shields.io/badge/MongoDB-5.0-green?logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-blue)](LICENSE)

**Boilerplate backend menggunakan GoFiber dan MongoDB** untuk membangun REST API cepat, scalable, dan mudah dikembangkan.  
Cocok untuk aplikasi dengan autentikasi, manajemen user, produk, event, dan modul lainnya.

---

## Fitur ✨

- **Authentication**: Register, Login, Logout
- **JWT Middleware**: Autentikasi token & role-based access
- **MongoDB**: Koneksi database siap pakai
- **Fiber Framework**: Ringan, cepat, dan mudah dikembangkan
- **Env Config**: Semua konfigurasi rahasia ada di `.env`
- **Extensible**: Tambah module baru (product, event, dsb) dengan mudah

---

## Prasyarat ⚡

- Go >= 1.24
- MongoDB (Atlas atau lokal)
- Git
- Postman / Insomnia (opsional untuk testing API)

---

## Instalasi & Setup 🛠️

### 1️⃣ Clone Repository

````bash
git clone https://github.com/MuliaAndiki/boilerpad-gofiber-with-monggodb.git
cd boilerpad-gofiber-with-monggodb


## Instalasi 🛠️

1. Clone repository:
```bash
git clone https://github.com/MuliaAndiki/boilerpad-gofiber-with-monggodb.git
cd boilerpad-gofiber-with-monggodb


2️⃣ Install Dependencies
```bash
go mod tidy

3️⃣ Buat File .env

touch .env
MONGO_URI=<your_mongodb_uri>
MONGO_DB=<your_db_name>
PORT=5000
JWT_SECRET=<your_jwt_secret>

4️⃣ Jalankan Server
```bash
go run main.go

server runing on port
http://localhost:5000


Struktur Project 📁
boilerpad-go/
│
├─ config/        # Config database & env
├─ controllers/   # Logic endpoint (Auth, User, dsb)
├─ middleware/    # JWT & role middleware
├─ models/        # MongoDB models
├─ routes/        # Routing API
├─ main.go        # Entry point aplikasi
├─ go.mod
├─ go.sum
└─ .env           # Environment variable (jangan di-push ke git)


Contoh Endpoint
POST /api/auth/register
Body (JSON):
{
  "fullname": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "role": "user"
}

POST /api/auth/login
Body (JSON):
{
  "email": "john@example.com",
  "password": "password123"
}
Response:
{
  "message": "Login berhasil",
  "token": "<JWT_TOKEN>",
  "user": {
    "id": "64f4a2b...xyz",
    "fullname": "John Doe",
    "email": "john@example.com",
    "role": "user"
  }
}


POST /api/auth/logout
Header:
Authorization: Bearer <JWT_TOKEN>
Response:
{
  "message": "Logout berhasil, token di-blacklist"
}
````
