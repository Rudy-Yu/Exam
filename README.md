# Online Exam App

Aplikasi Ujian Online berbasis web dengan backend Go (Fiber, GORM, PostgreSQL) dan frontend ReactJS.

---

## Status Fitur MVP (Update)

### Backend
- [x] Autentikasi user & admin (JWT, role)
- [x] CRUD soal (admin only)
- [x] Simpan jawaban peserta (auto-save & submit final)
- [x] Export hasil ujian (CSV, admin only)
- [x] Penilaian otomatis sederhana (benar/salah)
- [x] Struktur modular, mudah dikembangkan

### Frontend
- [x] Login & register user
- [x] Dashboard peserta (daftar ujian)
- [x] Halaman ujian (soal, input jawaban, timer)
- [x] Auto-save jawaban peserta
- [x] Submit jawaban & tampilkan skor
- [x] Logout
- [x] Panel admin CRUD soal (tambah, edit, hapus, export)
- [x] Notifikasi success/error pada aksi CRUD & submit
- [x] Validasi input form
- [x] Proteksi route (hanya admin bisa akses /admin)

### Fitur Lanjutan (Belum Ada)
- [ ] Manajemen user (CRUD user oleh admin)
- [ ] Laporan detail hasil ujian per peserta (frontend)
- [ ] Penilaian otomatis lanjutan (bobot soal, analisis jawaban)
- [ ] Fitur keamanan lanjutan (reset password, verifikasi email, dsb)
- [ ] Notifikasi/email
- [ ] UI/UX panel admin lebih lengkap (filter, search, pagination)

---

## Catatan
- Semua fitur inti MVP SUDAH SELESAI dan siap diuji/digunakan.
- Fitur lanjutan dapat dikembangkan sesuai kebutuhan selanjutnya.

---

## Fitur Utama
- **Registrasi & Login User** (dengan validasi ke database, password di-hash)
- **Dashboard Peserta**: Menampilkan daftar ujian yang tersedia
- **Mulai Ujian**: Peserta dapat memulai ujian, mengerjakan 5 soal, dan melihat skor
- **Timer Ujian**: Penghitung waktu mundur otomatis saat ujian berlangsung
- **Logout**: Peserta dapat keluar dari sesi
- **Lupa Password & Registrasi**: Halaman dummy sudah tersedia
- **Integrasi Database**: Data user, ujian, dan soal diambil dari database/PostgreSQL
- **API Backend**: Menggunakan Fiber (Go) dengan endpoint RESTful
- **Frontend Modern**: Menggunakan ReactJS dan react-router-dom

---

## Struktur Proyek

```
online-exam-app/
├── backend/           # Backend Go (Fiber, GORM, PostgreSQL)
│   ├── main.go
│   ├── models/
│   └── utils/
├── frontend/          # Frontend ReactJS
│   ├── public/
│   └── src/
│       ├── app.js
│       └── components/
│           ├── login.js
│           ├── dashboard.js
│           ├── Register.js
│           ├── ExamPage.js
│           └── ForgotPassword.js
├── schema.sql         # Skema database
└── README.md
```

---

## Cara Menjalankan

### 1. **Backend (Go + PostgreSQL)**
- Pastikan Go, PostgreSQL, dan Redis sudah terinstall.
- Edit `backend/main.go` jika perlu menyesuaikan koneksi database.
- Jalankan migrasi database (atau gunakan `schema.sql` jika perlu).
- Jalankan backend:
  ```sh
  cd backend
  go run main.go
  ```
- Backend berjalan di: `http://localhost:3000`

### 2. **Frontend (ReactJS)**
- Pastikan Node.js & npm sudah terinstall.
- Install dependencies:
  ```sh
  cd frontend
  npm install
  ```
- Jalankan frontend:
  ```sh
  npm start
  ```
- Frontend berjalan di: `http://localhost:3001`

---

## Fitur yang Sudah Berfungsi
- [x] Registrasi user baru (disimpan di database, password di-hash)
- [x] Login user (validasi ke database, hanya user terdaftar yang bisa login)
- [x] Dashboard peserta (menampilkan daftar ujian dari database)
- [x] Mulai ujian (soal diambil dari backend, 5 soal dummy)
- [x] Penghitung waktu mundur saat ujian
- [x] Submit jawaban & skor otomatis
- [x] Logout
- [x] Halaman registrasi & lupa password (dummy, bisa dikembangkan)

---

## Kontribusi
Pull request dan issue sangat diterima!

---

## Lisensi
MIT 