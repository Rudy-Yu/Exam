# API Testing Guide

Berikut adalah daftar link langsung ke semua endpoint utama backend yang bisa kamu klik atau copy-paste ke browser/API client (Postman, Insomnia, dsb).

> Untuk endpoint POST/PUT/DELETE, gunakan Postman atau tool serupa karena browser hanya bisa GET.

---

## AUTH & USER

- **Register User:**  
  `POST` [`http://localhost:3000/api/register`](http://localhost:3000/api/register)
- **Login User:**  
  `POST` [`http://localhost:3000/api/login`](http://localhost:3000/api/login)

---

## DASHBOARD PESERTA

- **List Ujian:**  
  `GET` [`http://localhost:3000/api/exams`](http://localhost:3000/api/exams)
- **Ambil Soal Ujian:**  
  `GET` [`http://localhost:3000/api/exam/1/questions`](http://localhost:3000/api/exam/1/questions)  
  (ganti 1 dengan id ujian lain)

---

## JAWABAN PESERTA

- **Mulai Ujian:**  
  `POST` [`http://localhost:3000/api/exam/1/start`](http://localhost:3000/api/exam/1/start)
- **Auto-save Jawaban Draft:**  
  `POST` [`http://localhost:3000/api/answers/draft`](http://localhost:3000/api/answers/draft)
- **Submit Jawaban:**  
  `POST` [`http://localhost:3000/api/answers/submit`](http://localhost:3000/api/answers/submit)
- **Ambil Timer Ujian:**  
  `GET` [`http://localhost:3000/api/exam/1/timer`](http://localhost:3000/api/exam/1/timer)

---

## ADMIN - SOAL

- **List Semua Soal:**  
  `GET` [`http://localhost:3000/api/admin/questions`](http://localhost:3000/api/admin/questions)
- **Tambah Soal:**  
  `POST` [`http://localhost:3000/api/admin/questions`](http://localhost:3000/api/admin/questions)
- **Edit Soal:**  
  `PUT` [`http://localhost:3000/api/admin/questions/1`](http://localhost:3000/api/admin/questions/1)  
  (ganti 1 dengan id soal)
- **Hapus Soal:**  
  `DELETE` [`http://localhost:3000/api/admin/questions/1`](http://localhost:3000/api/admin/questions/1)  
  (ganti 1 dengan id soal)

---

## ADMIN - USER

- **List Semua User:**  
  `GET` [`http://localhost:3000/api/admin/users`](http://localhost:3000/api/admin/users)
- **Tambah User:**  
  `POST` [`http://localhost:3000/api/admin/users`](http://localhost:3000/api/admin/users)
- **Edit User:**  
  `PUT` [`http://localhost:3000/api/admin/users/1`](http://localhost:3000/api/admin/users/1)  
  (ganti 1 dengan id user)
- **Hapus User:**  
  `DELETE` [`http://localhost:3000/api/admin/users/1`](http://localhost:3000/api/admin/users/1)  
  (ganti 1 dengan id user)

---

## ADMIN - EXPORT

- **Export Hasil Ujian (CSV):**  
  `GET` [`http://localhost:3000/api/admin/export`](http://localhost:3000/api/admin/export)

---

> **Catatan:**
> - Untuk endpoint yang butuh login/admin, pastikan kirim header `Authorization: Bearer <token>`.
> - Untuk POST/PUT, gunakan Postman/Insomnia dan kirim body JSON sesuai kebutuhan.
> - Ganti angka ID pada URL sesuai data di database.

---

# Cara Testing API dengan Postman & Insomnia

## 1. Testing dengan Postman

### A. Install Postman
- Download & install dari [postman.com/downloads](https://www.postman.com/downloads/)

### B. Langkah Umum
1. Buka Postman
2. Klik **New** → **HTTP Request**
3. Pilih **method** (GET, POST, PUT, DELETE)
4. Masukkan **URL endpoint** (misal: `http://localhost:3000/api/login`)
5. Jika POST/PUT, klik tab **Body** → pilih **raw** → pilih **JSON** → isi data JSON
6. Jika butuh token, klik tab **Headers** → tambahkan:
   ```
   Key: Authorization
   Value: Bearer <token>
   ```
7. Klik **Send**
8. Lihat response di bawah

### C. Contoh Testing

#### Register User
- Method: POST  
- URL: `http://localhost:3000/api/register`
- Body (raw, JSON):
  ```json
  {
    "name": "Nama Peserta",
    "email": "user1@example.com",
    "password": "password"
  }
  ```

#### Login User
- Method: POST  
- URL: `http://localhost:3000/api/login`
- Body (raw, JSON):
  ```json
  {
    "email": "user1@example.com",
    "password": "password"
  }
  ```
- **Copy token** dari response untuk request berikutnya.

#### GET Endpoint yang Butuh Token
- Method: GET  
- URL: `http://localhost:3000/api/exams`
- Header:
  ```
  Authorization: Bearer <token>
  ```

#### Tambah Soal (Admin)
- Method: POST  
- URL: `http://localhost:3000/api/admin/questions`
- Header:
  ```
  Authorization: Bearer <token_admin>
  ```
- Body (raw, JSON):
  ```json
  {
    "exam_id": 1,
    "question_text": "Contoh soal?",
    "type": "pilihan_ganda",
    "options": ["A", "B", "C", "D"],
    "correct_answer": "A"
  }
  ```

---

## 2. Testing dengan Insomnia

### A. Install Insomnia
- Download & install dari [insomnia.rest/download](https://insomnia.rest/download)

### B. Langkah Umum
1. Buka Insomnia
2. Klik **New Request**
3. Pilih **method** (GET, POST, PUT, DELETE)
4. Masukkan **URL endpoint**
5. Untuk POST/PUT, klik tab **Body** → pilih **JSON** → isi data JSON
6. Untuk token, klik tab **Header** → tambahkan:
   ```
   Name: Authorization
   Value: Bearer <token>
   ```
7. Klik **Send**
8. Lihat response di bawah

---

## 3. Tips
- Untuk endpoint yang butuh token, **login dulu** dan copy token dari response.
- Untuk POST/PUT, **pilih raw/JSON** dan isi body sesuai kebutuhan.
- Untuk GET/DELETE, cukup URL dan header jika perlu.

---

## 4. Troubleshooting
- **401 Unauthorized:** Token salah/expired/tidak dikirim.
- **405 Method Not Allowed:** Method tidak sesuai (misal, POST ke endpoint yang hanya GET).
- **400 Bad Request:** Body JSON tidak sesuai.
- **500 Internal Server Error:** Cek log backend.

---

Jika ingin contoh file koleksi Postman/Insomnia (import .json), kabari saja!
Jika ada error, lampirkan screenshot atau pesan error detail agar bisa dibantu lebih lanjut. 