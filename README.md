# Online Exam App 📚

Aplikasi ujian online berbasis web dengan fitur auto-save, timer sinkronisasi server, dan panel admin.

## 📋 Daftar Isi
1. [Fitur Utama](#fitur-utama)
2. [Teknologi yang Digunakan](#teknologi-yang-digunakan)
3. [Cara Install](#cara-install)
4. [Cara Menjalankan](#cara-menjalankan)
5. [Panduan Penggunaan](#panduan-penggunaan)
6. [Deployment ke Production](#deployment-ke-production)
7. [Troubleshooting](#troubleshooting)

## 🚀 Fitur Utama

### Peserta Ujian
- ✅ Login & register akun
- ✅ Lihat daftar ujian yang tersedia
- ✅ Mulai ujian dengan timer
- ✅ Auto-save jawaban setiap 30 detik
- ✅ Submit jawaban & lihat skor langsung

### Admin
- ✅ Panel admin terpisah
- ✅ Kelola soal ujian (tambah, edit, hapus)
- ✅ Kelola user (tambah, edit, hapus)
- ✅ Export hasil ujian (CSV)
- ✅ Lihat statistik peserta

### Sistem
- ✅ Auto-save dengan Redis
- ✅ Timer sinkronisasi server
- ✅ Keamanan dengan JWT
- ✅ Support HTTPS/SSL
- ✅ Optimasi database

## 🛠 Teknologi yang Digunakan

### Backend
- Go (Fiber Framework)
- PostgreSQL (Database)
- Redis (Cache & Session)
- JWT (Autentikasi)

### Frontend
- React.js
- Axios
- React Router
- Modern CSS

## 📥 Cara Install

### 1. Prasyarat
```bash
# Install Go (minimal versi 1.16)
https://golang.org/dl/

# Install Node.js (minimal versi 14)
https://nodejs.org/

# Install PostgreSQL
https://www.postgresql.org/download/

# Install Redis
https://redis.io/download
```

### 2. Clone Repository
```bash
git clone https://github.com/username/online-exam-app.git
cd online-exam-app
```

### 3. Setup Database
```sql
-- Buka PostgreSQL dan jalankan:
CREATE DATABASE examdb;
```

### 4. Install Dependencies

```bash
# Backend dependencies
cd backend
go mod download

# Frontend dependencies
cd frontend
npm install
```

### 5. Konfigurasi Environment

#### Backend (.env)
```env
# Development
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASS=admin
DB_NAME=examdb
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=
JWT_SECRET=your-secret-key
```

#### Frontend (.env)
```env
# Development
REACT_APP_API_URL=http://localhost:3000
```

## 🚀 Cara Menjalankan

### Mode Development

1. **Jalankan Redis**
```bash
# Windows
redis-server

# Linux/Mac
redis-server /etc/redis/redis.conf
```

2. **Jalankan Backend**
```bash
cd backend
go run main.go
```
Backend akan berjalan di http://localhost:3000

3. **Jalankan Frontend**
```bash
cd frontend
npm start
```
Frontend akan berjalan di http://localhost:3001

## 📖 Panduan Penggunaan

### 👨‍💻 Sebagai Peserta Ujian

1. **Register/Login**
   - Buka http://localhost:3001
   - Klik "Daftar Baru" jika belum punya akun
   - Isi email dan password
   - Login dengan akun yang sudah dibuat

2. **Mulai Ujian**
   - Setelah login, lihat daftar ujian di dashboard
   - Klik "Mulai Ujian" pada ujian yang ingin dikerjakan
   - Timer akan mulai berjalan
   - Jawaban tersimpan otomatis setiap 30 detik
   - Klik "Submit" jika sudah selesai

3. **Lihat Hasil**
   - Skor akan muncul setelah submit
   - Kembali ke dashboard untuk ujian lainnya

### 👨‍🏫 Sebagai Admin

1. **Akses Panel Admin**
   - Login dengan akun admin
   - Klik tab "Admin Panel"

2. **Kelola Soal**
   - Klik tab "Kelola Soal"
   - Tambah soal baru dengan form
   - Edit/hapus soal yang ada
   - Set bobot nilai tiap soal

3. **Kelola User**
   - Klik tab "Kelola User"
   - Lihat daftar user
   - Tambah/edit/hapus user
   - Set role user (admin/peserta)

4. **Export Hasil**
   - Klik "Export Hasil Ujian (CSV)"
   - File akan terdownload otomatis

## 🌐 Deployment ke Production

### 1. Persiapan Server
- Setup VPS/server dengan Ubuntu/Debian
- Install Nginx sebagai reverse proxy
- Setup domain dan SSL dengan Let's Encrypt

### 2. Environment Production
```env
# Backend
GO_ENV=production
DB_HOST=your-db-host
REDIS_HOST=your-redis-host
JWT_SECRET=strong-secret-key
SSL_CERT=/etc/letsencrypt/live/domain.com/fullchain.pem
SSL_KEY=/etc/letsencrypt/live/domain.com/privkey.pem

# Frontend
REACT_APP_API_URL=https://api.domain.com
```

### 3. Build & Deploy
```bash
# Build frontend
cd frontend
npm run build

# Copy ke server
scp -r build/* user@server:/var/www/html/

# Jalankan backend
cd backend
go build
./online-exam-app
```

### 4. Nginx Configuration
```nginx
# /etc/nginx/sites-available/online-exam
server {
    listen 80;
    server_name domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name domain.com;

    ssl_certificate /etc/letsencrypt/live/domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/domain.com/privkey.pem;

    location / {
        root /var/www/html;
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## ❗ Troubleshooting

### Backend Issues

1. **Error: Connection refused**
   - Pastikan PostgreSQL running
   - Cek credentials di .env
   - Cek port PostgreSQL

2. **Error: Redis connection failed**
   - Pastikan Redis server running
   - Cek Redis port dan password

### Frontend Issues

1. **API Error**
   - Pastikan backend running
   - Cek REACT_APP_API_URL di .env
   - Periksa network di DevTools

2. **Auto-save tidak berfungsi**
   - Cek koneksi Redis
   - Periksa token JWT
   - Lihat log di console

### Security Issues

1. **Invalid Token**
   - Token expired - login ulang
   - Cek JWT_SECRET di backend
   - Pastikan token dikirim di header

2. **CORS Error**
   - Cek origin di CORS config
   - Pastikan protocol (http/https) sama

## 📞 Support

Jika mengalami masalah atau butuh bantuan:
- Buat issue di GitHub
- Email: support@domain.com
- Dokumentasi API: /api/docs

## 📄 Lisensi
MIT License - Silakan gunakan dan modifikasi sesuai kebutuhan. 