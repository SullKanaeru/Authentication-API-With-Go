

# Go Authentication API with Fiber & PostgreSQL

Proyek ini adalah sistem backend RESTful API yang tangguh untuk menangani autentikasi pengguna. Dibangun menggunakan **Go (Golang)** dengan *framework* **Fiber**, ORM **GORM**, dan database **PostgreSQL**. 

Sistem ini menerapkan standar keamanan modern menggunakan **JWT (JSON Web Tokens)** yang disimpan dengan aman di dalam *HTTP-Only Cookies*, serta dilengkapi sistem verifikasi ganda (OTP) via **Email (SMTP)** dan **WhatsApp (Fonnte)**.

## 🚀 Fitur Utama
* **Arsitektur Modular:** Menggunakan struktur folder yang rapi (pemisahan antara Handler, Service, Repository, dan Model) agar mudah dikelola dan diskalakan.
* **Registrasi & Login:** Mendukung penggunaan *Email* maupun *Nomor Telepon*.
* **Keamanan Lanjut:** Enkripsi *password* menggunakan `Bcrypt` dan penyimpanan token via *HTTP-Only Cookies* untuk mencegah serangan XSS.
* **Verifikasi OTP Otomatis:** Pengiriman kode OTP dinamis untuk aktivasi akun.
* **Auto Migration:** Tabel database otomatis dibuat berdasarkan *struct* model (GORM).
* **Protected Routes:** Akses ke data privat dilindungi oleh *middleware* JWT.

## 🛠️ Tech Stack
* **Language:** Go (Golang)
* **Web Framework:** [Fiber v2](https://gofiber.io/)
* **ORM:** [GORM](https://gorm.io/)
* **Database:** PostgreSQL
* **Security:** `golang-jwt/jwt/v5`, `golang.org/x/crypto/bcrypt`
* **Configuration:** `joho/godotenv`

## 📁 Struktur Folder
```text
my-backend/
├── .env                 # File rahasia (Tidak di-commit ke Git)
├── cmd/
│   └── api/
│       └── main.go      # Entry point utama aplikasi
├── internal/
│   ├── config/          # Pengaturan environment variables
│   ├── handler/         # Layer Controller (menerima HTTP request)
│   ├── middleware/      # Layer penengah (misal: pelindung JWT)
│   ├── model/           # Definisi struktur database dan DTO
│   ├── repository/      # Layer komunikasi langsung dengan PostgreSQL
│   └── service/         # Layer logika bisnis & notifikasi (OTP)
├── pkg/
│   └── database/        # Setup koneksi GORM dan Auto-Migrate
├── go.mod               # Daftar dependensi modul
└── go.sum               # Checksum keamanan modul

```

## ⚙️ Persyaratan Instalasi

Pastikan kamu telah menginstal:

1. [Go](https://go.dev/dl/) (versi 1.18 atau lebih baru)
2. [PostgreSQL](https://www.postgresql.org/download/) (berjalan di *local* atau *server*)

## 🚀 Cara Menjalankan Proyek

**1. Clone atau masuk ke direktori proyek**

```bash
cd my-backend

```

**2. Instal Dependensi**

```bash
go mod tidy

```

**3. Siapkan Database**
Buat database kosong di PostgreSQL kamu (misalnya dengan nama `mydb`).

**4. Konfigurasi Environment Variables**
Buat file `.env` di *root* direktori proyek dan isi dengan konfigurasi berikut:

```env
# Server
PORT=3000

# Database PostgreSQL
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=password_postgres_kamu
DB_NAME=mydb
DB_PORT=5432

# Keamanan
JWT_SECRET=super_secret_key_kamu_disini

# Notifikasi Email (SMTP Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=email_kamu@gmail.com
SMTP_PASS=password_aplikasi_gmail_kamu
SMTP_SENDER_NAME=Aplikasi Saya

# Notifikasi WhatsApp (Fonnte)
FONNTE_TOKEN=token_fonnte_kamu

```

**5. Jalankan Server**

```bash
go run cmd/api/main.go

```

*GORM akan otomatis melakukan migrasi dan membuat tabel `users` di database kamu. Server akan berjalan di `http://localhost:3000`.*

---

## 📖 Dokumentasi API Endpoint

Berikut adalah daftar *endpoint* yang tersedia untuk diuji coba melalui Postman atau *Front-end*.

### 1. Register User Baru

* **URL:** `/api/v1/auth/register`
* **Method:** `POST`
* **Body Request (JSON):**
```json
{
    "fullname": "John Doe",
    "username": "johndoe99",
    "email": "johndoe@example.com",
    "phone_number": "08123456789",
    "password": "password123",
    "confirm_password": "password123"
}

```


* **Response (201 Created):** Memicu pengiriman OTP ke Email atau WhatsApp.

### 2. Verifikasi OTP

* **URL:** `/api/v1/auth/verify`
* **Method:** `POST`
* **Body Request (JSON):**
```json
{
    "identifier": "08123456789", 
    "otp": "123456" 
}

```


* **Response (200 OK):** Mengubah status `is_verified` menjadi `true`.

### 3. Login

* **URL:** `/api/v1/auth/login`
* **Method:** `POST`
* **Body Request (JSON):**
```json
{
    "identifier": "johndoe@example.com", 
    "password": "password123"
}

```


* **Response (200 OK):** JWT Token akan otomatis tersimpan di *HTTP-Only Cookies* pada *browser* atau Postman.

### 4. Dapatkan Profil (*Protected Route*)

* **URL:** `/api/v1/users/me`
* **Method:** `GET`
* **Headers:** Membutuhkan `jwt_token` di dalam *Cookie*.
* **Response (200 OK):**
```json
{
    "message": "Berhasil mengambil data profil",
    "data": {
        "id": 1,
        "fullname": "John Doe",
        "username": "johndoe99",
        "email": "johndoe@example.com",
        "phone_number": "08123456789",
        "role": "user",
        "is_verified": true,
        "created_at": "2026-02-20T10:00:00Z"
    }
}

```



### 5. Logout

* **URL:** `/api/v1/auth/logout`
* **Method:** `POST`
* **Response (200 OK):** Menghapus *Cookie* `jwt_token` dari *client*.

```

***
