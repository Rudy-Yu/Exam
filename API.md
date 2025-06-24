# 📚 API Documentation

## 🔑 Authentication

### Register User
```http
POST /api/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}

Response:
{
    "success": true,
    "message": "Pendaftaran berhasil"
}
```

### Login
```http
POST /api/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}

Response:
{
    "success": true,
    "message": "Login berhasil",
    "token": "jwt_token_here",
    "role": "user"
}
```

## 📝 Exam Endpoints

### Get Available Exams
```http
GET /api/exams
Authorization: Bearer <token>

Response:
[
    {
        "id": 1,
        "title": "Matematika Dasar",
        "duration": 3600
    }
]
```

### Start Exam
```http
POST /api/exam/:id/start
Authorization: Bearer <token>

Response:
{
    "success": true,
    "start_time": "2024-01-20T10:00:00Z",
    "duration": 3600
}
```

### Get Exam Questions
```http
GET /api/exam/:id/questions
Authorization: Bearer <token>

Response:
[
    {
        "id": 1,
        "text": "2 + 2 = ?",
        "options": ["3", "4", "5", "6"],
        "weight": 1
    }
]
```

### Get Exam Timer
```http
GET /api/exam/:id/timer
Authorization: Bearer <token>

Response:
{
    "success": true,
    "remaining_time": 1800
}
```

### Auto-save Answer
```http
POST /api/answers/draft
Authorization: Bearer <token>
Content-Type: application/json

{
    "question_id": 1,
    "answer_text": "4"
}

Response:
{
    "success": true,
    "message": "Jawaban tersimpan sementara"
}
```

### Submit Final Answers
```http
POST /api/answers/submit
Authorization: Bearer <token>
Content-Type: application/json

[
    {
        "question_id": 1,
        "answer_text": "4"
    }
]

Response:
{
    "success": true,
    "message": "Jawaban berhasil disubmit"
}
```

## 👨‍🏫 Admin Endpoints

### Get All Users
```http
GET /api/admin/users
Authorization: Bearer <token>

Response:
[
    {
        "id": 1,
        "email": "user@example.com",
        "role": "user"
    }
]
```

### Create User
```http
POST /api/admin/users
Authorization: Bearer <token>
Content-Type: application/json

{
    "email": "newuser@example.com",
    "password": "password123",
    "role": "user"
}

Response:
{
    "success": true,
    "message": "User berhasil dibuat"
}
```

### Update User
```http
PUT /api/admin/users/:id
Authorization: Bearer <token>
Content-Type: application/json

{
    "email": "updated@example.com",
    "role": "admin"
}

Response:
{
    "success": true,
    "message": "User berhasil diupdate"
}
```

### Delete User
```http
DELETE /api/admin/users/:id
Authorization: Bearer <token>

Response:
{
    "success": true,
    "message": "User berhasil dihapus"
}
```

### Get All Questions
```http
GET /api/admin/questions
Authorization: Bearer <token>

Response:
[
    {
        "id": 1,
        "exam_id": 1,
        "question_text": "2 + 2 = ?",
        "correct_answer": "4",
        "weight": 1
    }
]
```

### Create Question
```http
POST /api/admin/questions
Authorization: Bearer <token>
Content-Type: application/json

{
    "exam_id": 1,
    "question_text": "2 + 2 = ?",
    "correct_answer": "4",
    "weight": 1
}

Response:
{
    "success": true,
    "message": "Soal berhasil ditambah"
}
```

### Update Question
```http
PUT /api/admin/questions/:id
Authorization: Bearer <token>
Content-Type: application/json

{
    "exam_id": 1,
    "question_text": "Updated question",
    "correct_answer": "Updated answer",
    "weight": 2
}

Response:
{
    "success": true,
    "message": "Soal berhasil diupdate"
}
```

### Delete Question
```http
DELETE /api/admin/questions/:id
Authorization: Bearer <token>

Response:
{
    "success": true,
    "message": "Soal berhasil dihapus"
}
```

### Export Results
```http
GET /api/admin/export
Authorization: Bearer <token>

Response: CSV File Download
```

## 🔒 Error Responses

### Unauthorized
```http
Status: 401
{
    "success": false,
    "message": "Token tidak valid"
}
```

### Forbidden
```http
Status: 403
{
    "success": false,
    "message": "Akses admin diperlukan"
}
```

### Not Found
```http
Status: 404
{
    "success": false,
    "message": "Data tidak ditemukan"
}
```

### Server Error
```http
Status: 500
{
    "success": false,
    "message": "Terjadi kesalahan server"
}
```

## 📝 Notes

- Semua request yang memerlukan autentikasi harus menyertakan header `Authorization: Bearer <token>`
- Token JWT memiliki masa berlaku 24 jam
- Response selalu dalam format JSON kecuali untuk endpoint export
- Semua timestamp menggunakan format ISO 8601
- Error response selalu menyertakan field `success` dan `message` 