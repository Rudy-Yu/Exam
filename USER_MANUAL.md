# 👥 User Manual

## 📖 Panduan Penggunaan Aplikasi Online Exam

Panduan lengkap untuk menggunakan aplikasi ujian online sebagai peserta dan admin.

## 🎯 Daftar Isi
1. [Peserta Ujian](#peserta-ujian)
2. [Administrator](#administrator)
3. [FAQ](#faq)
4. [Troubleshooting](#troubleshooting)

---

## 👨‍🎓 Peserta Ujian

### 1. Registrasi & Login

#### Registrasi Akun Baru
1. Buka aplikasi di browser: `https://yourdomain.com`
2. Klik tombol **"Daftar Baru"**
3. Isi form registrasi:
   - **Email**: Masukkan email yang valid
   - **Password**: Minimal 8 karakter
4. Klik **"Daftar"**
5. Pesan sukses akan muncul

#### Login ke Sistem
1. Masukkan **email** dan **password**
2. Klik **"Login"**
3. Jika berhasil, akan diarahkan ke dashboard

### 2. Dashboard Peserta

#### Melihat Daftar Ujian
- Dashboard menampilkan semua ujian yang tersedia
- Informasi yang ditampilkan:
  - Nama ujian
  - Durasi ujian
  - Status (tersedia/selesai)

#### Memulai Ujian
1. Pilih ujian yang ingin dikerjakan
2. Klik tombol **"Mulai Ujian"**
3. Timer akan mulai berjalan
4. Soal akan dimuat otomatis

### 3. Mengerjakan Ujian

#### Interface Ujian
```
┌─────────────────────────────────────┐
│ Timer: 45:30 | Ujian: Matematika   │
├─────────────────────────────────────┤
│ 1. Berapakah hasil dari 2 + 2?     │
│    ○ 3                             │
│    ● 4                             │ ← Pilihan yang dipilih
│    ○ 5                             │
│    ○ 6                             │
├─────────────────────────────────────┤
│ [Submit Jawaban] [Kembali]         │
└─────────────────────────────────────┘
```

#### Fitur Auto-save
- ✅ Jawaban tersimpan otomatis setiap 30 detik
- ✅ Notifikasi "Jawaban tersimpan otomatis" akan muncul
- ✅ Tidak perlu khawatir kehilangan jawaban

#### Navigasi Soal
- Soal ditampilkan satu per satu
- Gunakan tombol navigasi untuk berpindah soal
- Semua soal harus dijawab sebelum submit

### 4. Submit & Hasil

#### Submit Jawaban
1. Pastikan semua soal sudah dijawab
2. Klik tombol **"Submit Jawaban"**
3. Konfirmasi submit
4. Tunggu proses penilaian

#### Melihat Hasil
- Skor akan ditampilkan segera setelah submit
- Format: "Skor Anda: X dari Y"
- Klik **"Kembali ke Dashboard"** untuk melanjutkan

### 5. Logout
- Klik tombol **"Logout"** di pojok kanan atas
- Sesi akan berakhir
- Kembali ke halaman login

---

## 👨‍🏫 Administrator

### 1. Akses Panel Admin

#### Login sebagai Admin
1. Login dengan akun admin
2. Klik tab **"Admin Panel"**
3. Panel admin akan terbuka dengan 2 tab:
   - **Kelola Soal**
   - **Kelola User**

### 2. Kelola Soal

#### Menambah Soal Baru
1. Klik tab **"Kelola Soal"**
2. Isi form tambah soal:
   - **Exam ID**: ID ujian (1, 2, 3, dst)
   - **Pertanyaan**: Teks soal
   - **Jawaban Benar**: Jawaban yang benar
3. Klik **"Tambah"**
4. Soal akan muncul di tabel

#### Edit Soal
1. Klik tombol **"Edit"** pada soal yang ingin diubah
2. Form akan terisi dengan data soal
3. Ubah data yang diperlukan
4. Klik **"Update"**

#### Hapus Soal
1. Klik tombol **"Hapus"** pada soal
2. Konfirmasi penghapusan
3. Soal akan dihapus dari database

#### Export Hasil Ujian
1. Klik tombol **"Export Hasil Ujian (CSV)"**
2. File CSV akan terdownload otomatis
3. File berisi data:
   - Email peserta
   - Skor per ujian
   - Waktu submit

### 3. Kelola User

#### Melihat Daftar User
- Tabel menampilkan semua user terdaftar
- Informasi: ID, Email, Role

#### Menambah User Baru
1. Klik tab **"Kelola User"**
2. Isi form:
   - **Email**: Email user baru
   - **Role**: Pilih "user" atau "admin"
3. Klik **"Tambah"**
4. Password akan digenerate otomatis

#### Edit User
1. Klik tombol **"Edit"** pada user
2. Ubah email atau role
3. Klik **"Update"**

#### Hapus User
1. Klik tombol **"Hapus"** pada user
2. Konfirmasi penghapusan
3. User akan dihapus dari sistem

### 4. Monitoring & Laporan

#### Dashboard Admin
- Jumlah user terdaftar
- Jumlah ujian aktif
- Statistik penggunaan

#### Log Aktivitas
- Login/logout user
- Aktivitas admin
- Error sistem

---

## ❓ FAQ (Frequently Asked Questions)

### Peserta Ujian

**Q: Apa yang terjadi jika internet terputus saat ujian?**
A: Jawaban tersimpan otomatis setiap 30 detik di server. Setelah internet kembali, jawaban akan tetap ada.

**Q: Bisakah saya mengubah jawaban setelah submit?**
A: Tidak, setelah submit jawaban tidak bisa diubah. Pastikan semua jawaban sudah benar sebelum submit.

**Q: Apa yang terjadi jika waktu ujian habis?**
A: Ujian akan otomatis submit dan jawaban yang sudah diisi akan dinilai.

**Q: Bagaimana jika saya tidak bisa login?**
A: Pastikan email dan password benar. Jika lupa password, hubungi admin.

**Q: Bisakah saya mengerjakan ujian yang sama berkali-kali?**
A: Tergantung kebijakan admin. Biasanya setiap ujian hanya bisa dikerjakan sekali.

### Administrator

**Q: Bagaimana cara reset password user?**
A: Edit user dan set password baru, atau hapus user dan buat ulang.

**Q: Bisakah saya melihat jawaban detail per user?**
A: Ya, melalui fitur export CSV yang berisi detail jawaban.

**Q: Bagaimana cara backup data?**
A: Database di-backup otomatis setiap hari. Hubungi sistem administrator untuk restore.

**Q: Bisakah saya menambah ujian baru?**
A: Ya, dengan menambah soal baru dan set Exam ID yang sesuai.

**Q: Bagaimana cara mengatur durasi ujian?**
A: Durasi ujian diatur di backend. Hubungi developer untuk perubahan.

---

## 🔧 Troubleshooting

### Masalah Umum

#### 1. Tidak Bisa Login
**Gejala**: Error "Email atau password salah"
**Solusi**:
- Periksa email dan password
- Pastikan caps lock tidak aktif
- Hubungi admin jika lupa password

#### 2. Ujian Tidak Muncul
**Gejala**: Dashboard kosong
**Solusi**:
- Refresh halaman
- Logout dan login ulang
- Hubungi admin jika masalah berlanjut

#### 3. Auto-save Tidak Berfungsi
**Gejala**: Tidak ada notifikasi auto-save
**Solusi**:
- Periksa koneksi internet
- Refresh halaman
- Coba browser lain

#### 4. Timer Tidak Akurat
**Gejala**: Timer tidak sinkron
**Solusi**:
- Refresh halaman
- Timer akan sinkron otomatis dengan server

#### 5. Submit Gagal
**Gejala**: Error saat submit
**Solusi**:
- Periksa koneksi internet
- Pastikan semua soal dijawab
- Coba submit ulang

### Masalah Admin

#### 1. Tidak Bisa Akses Panel Admin
**Gejala**: Tab admin tidak muncul
**Solusi**:
- Pastikan login dengan akun admin
- Logout dan login ulang
- Hubungi sistem administrator

#### 2. Export Gagal
**Gejala**: File CSV tidak terdownload
**Solusi**:
- Periksa popup blocker
- Coba browser lain
- Hubungi sistem administrator

#### 3. Soal Tidak Tersimpan
**Gejala**: Error saat tambah/edit soal
**Solusi**:
- Periksa semua field terisi
- Pastikan format data benar
- Hubungi sistem administrator

---

## 📞 Kontak Support

### Peserta Ujian
- **Email**: support@company.com
- **WhatsApp**: +62-xxx-xxx-xxxx
- **Jam Kerja**: Senin-Jumat 08:00-17:00 WIB

### Administrator
- **Email**: admin@company.com
- **Telepon**: +62-xxx-xxx-xxxx
- **Emergency**: +62-xxx-xxx-xxxx

### Dokumentasi Teknis
- **API Docs**: https://yourdomain.com/api/docs
- **GitHub**: https://github.com/username/online-exam-app
- **Wiki**: https://wiki.company.com

---

## 📋 Checklist Penggunaan

### Sebelum Ujian
- [ ] Pastikan koneksi internet stabil
- [ ] Siapkan browser yang kompatibel
- [ ] Tutup aplikasi lain yang tidak perlu
- [ ] Siapkan ruangan yang tenang

### Saat Ujian
- [ ] Baca instruksi dengan teliti
- [ ] Jawab semua soal
- [ ] Perhatikan timer
- [ ] Submit sebelum waktu habis

### Setelah Ujian
- [ ] Catat skor yang diperoleh
- [ ] Logout dari sistem
- [ ] Simpan bukti ujian (jika diperlukan)

---

## 🔄 Update & Versi

### Versi Terbaru
- **Versi**: 1.0.0
- **Tanggal**: Januari 2024
- **Fitur Baru**: Auto-save, Timer sinkronisasi, Panel admin

### Riwayat Update
- **v1.0.0**: Release awal dengan fitur dasar
- **v1.1.0**: Penambahan auto-save
- **v1.2.0**: Timer sinkronisasi server
- **v1.3.0**: Panel admin lengkap

---

**Catatan**: Panduan ini akan diperbarui sesuai dengan perkembangan aplikasi. Pastikan selalu menggunakan versi terbaru. 