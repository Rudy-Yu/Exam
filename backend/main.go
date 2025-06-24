package main

import (
    "fmt"
    "log"
    "time"
    "encoding/csv"
    "strconv"
    "bytes"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/go-redis/redis/v8"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v4"
)

// User model
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"unique"`
    Password string
    Role     string `gorm:"default:user"`
}

// Question model
type Question struct {
    ID            uint   `gorm:"primaryKey" json:"id"`
    ExamID        uint   `json:"exam_id"`
    QuestionText  string `json:"question_text"`
    CorrectAnswer string `json:"correct_answer"`
}

// Answer model
type Answer struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    ParticipantID uint     `json:"participant_id"`
    QuestionID   uint      `json:"question_id"`
    AnswerText   string    `json:"answer_text"`
    SubmittedAt  time.Time `json:"submitted_at"`
    IsDraft      bool      `json:"is_draft"`
}

// Database connection
func connectDB() *gorm.DB {
    dsn := "host=localhost user=postgres password=admin dbname=examdb port=5433 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    return db
}

// Redis connection
func connectRedis() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,   // use default DB
    })
    return rdb
}

func adminMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Token tidak ditemukan"})
    }
    tokenString := authHeader[len("Bearer "):] // "Bearer <token>"
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })
    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Token tidak valid"})
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["role"] != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Akses admin diperlukan"})
    }
    return c.Next()
}

func main() {
    // Connect to PostgreSQL
    db := connectDB()
    db.AutoMigrate(&User{}, &Question{}, &Answer{}) // Migrasi tabel users, questions, answers

    // Connect to Redis
    rdb := connectRedis()
    defer rdb.Close()

    // Initialize Fiber app
    app := fiber.New()

    // Add CORS middleware
    app.Use(cors.New())

    // Example route: Get all exams
    app.Get("/api/exams", func(c *fiber.Ctx) error {
        var exams []map[string]interface{}
        result := db.Raw("SELECT * FROM exams").Scan(&exams)
        if result.Error != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": "Failed to fetch exams"})
        }
        return c.JSON(exams)
    })

    // Real Register Endpoint
    app.Post("/api/register", func(c *fiber.Ctx) error {
        var data map[string]string
        if err := c.BodyParser(&data); err != nil {
            return err
        }

        // Check if email already exists
        var existingUser User
        db.Where("email = ?", data["email"]).First(&existingUser)
        if existingUser.ID != 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Email sudah terdaftar"})
        }

        password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
        user := User{
            Email:    data["email"],
            Password: string(password),
            Role:     "user", // default role
        }

        result := db.Create(&user)
        if result.Error != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mendaftarkan user"})
        }

        // Verify user was created
        var createdUser User
        db.First(&createdUser, user.ID)
        if createdUser.ID == 0 {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal memverifikasi pendaftaran"})
        }

        return c.JSON(fiber.Map{"success": true, "message": "Pendaftaran berhasil"})
    })

    // Real Login Endpoint
    app.Post("/api/login", func(c *fiber.Ctx) error {
        var data map[string]string
        if err := c.BodyParser(&data); err != nil {
            return err
        }

        var user User
        db.Where("email = ?", data["email"]).First(&user)

        if user.ID == 0 {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "User tidak ditemukan"})
        }

        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Password salah"})
        }

        // Generate JWT token
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID,
            "email": user.Email,
            "role": user.Role,
            "exp": time.Now().Add(time.Hour * 24).Unix(),
        })
        tokenString, err := token.SignedString([]byte("secret"))
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal generate token"})
        }

        return c.JSON(fiber.Map{
            "success": true,
            "message": "Login berhasil",
            "token": tokenString,
            "role": user.Role,
        })
    })

    // Endpoint to get questions for an exam
    app.Get("/api/exam/:id/questions", func(c *fiber.Ctx) error {
        // Dummy questions
        questions := []map[string]interface{}{
            {"id": 1, "text": "Ibu kota Indonesia adalah?", "options": []string{"Jakarta", "Bandung", "Surabaya", "Medan"}, "correct_answer": "Jakarta"},
            {"id": 2, "text": "2 + 2 * 2 = ?", "options": []string{"8", "6", "4", "10"}, "correct_answer": "6"},
            {"id": 3, "text": "Siapakah presiden pertama Indonesia?", "options": []string{"Soeharto", "B.J. Habibie", "Soekarno", "Joko Widodo"}, "correct_answer": "Soekarno"},
            {"id": 4, "text": "Planet terbesar di tata surya adalah?", "options": []string{"Bumi", "Mars", "Jupiter", "Saturnus"}, "correct_answer": "Jupiter"},
            {"id": 5, "text": "HTML adalah singkatan dari?", "options": []string{"HyperText Markup Language", "High-level Text Language", "Hyper Transfer Language", "HyperText Main Language"}, "correct_answer": "HyperText Markup Language"},
        }
        return c.JSON(questions)
    })

    // Contoh endpoint admin (hanya admin)
    app.Get("/api/admin/only", adminMiddleware, func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"success": true, "message": "Halo Admin!"})
    })

    // CRUD Question (admin only)
    app.Post("/api/admin/questions", adminMiddleware, func(c *fiber.Ctx) error {
        var q Question
        if err := c.BodyParser(&q); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
        }
        if err := db.Create(&q).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal menambah soal"})
        }
        return c.JSON(fiber.Map{"success": true, "question": q})
    })
    app.Get("/api/admin/questions", adminMiddleware, func(c *fiber.Ctx) error {
        var questions []Question
        if err := db.Find(&questions).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data soal"})
        }
        return c.JSON(questions)
    })
    app.Put("/api/admin/questions/:id", adminMiddleware, func(c *fiber.Ctx) error {
        id := c.Params("id")
        var q Question
        if err := db.First(&q, id).Error; err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Soal tidak ditemukan"})
        }
        var data Question
        if err := c.BodyParser(&data); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
        }
        q.ExamID = data.ExamID
        q.QuestionText = data.QuestionText
        q.CorrectAnswer = data.CorrectAnswer
        if err := db.Save(&q).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal update soal"})
        }
        return c.JSON(fiber.Map{"success": true, "question": q})
    })
    app.Delete("/api/admin/questions/:id", adminMiddleware, func(c *fiber.Ctx) error {
        id := c.Params("id")
        if err := db.Delete(&Question{}, id).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal hapus soal"})
        }
        return c.JSON(fiber.Map{"success": true, "message": "Soal dihapus"})
    })

    // Auto-save jawaban peserta (draft)
    app.Post("/api/answers/draft", func(c *fiber.Ctx) error {
        var data struct {
            ParticipantID uint   `json:"participant_id"`
            QuestionID    uint   `json:"question_id"`
            AnswerText    string `json:"answer_text"`
        }
        if err := c.BodyParser(&data); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
        }
        var answer Answer
        if err := db.Where("participant_id = ? AND question_id = ?", data.ParticipantID, data.QuestionID).First(&answer).Error; err == nil {
            answer.AnswerText = data.AnswerText
            answer.IsDraft = true
            answer.SubmittedAt = time.Now()
            db.Save(&answer)
        } else {
            answer = Answer{
                ParticipantID: data.ParticipantID,
                QuestionID: data.QuestionID,
                AnswerText: data.AnswerText,
                IsDraft: true,
                SubmittedAt: time.Now(),
            }
            db.Create(&answer)
        }
        return c.JSON(fiber.Map{"success": true, "message": "Draft jawaban disimpan"})
    })
    // Submit jawaban final
    app.Post("/api/answers/submit", func(c *fiber.Ctx) error {
        var data struct {
            ParticipantID uint   `json:"participant_id"`
            QuestionID    uint   `json:"question_id"`
            AnswerText    string `json:"answer_text"`
        }
        if err := c.BodyParser(&data); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Data tidak valid"})
        }
        var answer Answer
        if err := db.Where("participant_id = ? AND question_id = ?", data.ParticipantID, data.QuestionID).First(&answer).Error; err == nil {
            answer.AnswerText = data.AnswerText
            answer.IsDraft = false
            answer.SubmittedAt = time.Now()
            db.Save(&answer)
        } else {
            answer = Answer{
                ParticipantID: data.ParticipantID,
                QuestionID: data.QuestionID,
                AnswerText: data.AnswerText,
                IsDraft: false,
                SubmittedAt: time.Now(),
            }
            db.Create(&answer)
        }
        return c.JSON(fiber.Map{"success": true, "message": "Jawaban final disimpan"})
    })

    // Export hasil ujian (admin only)
    app.Get("/api/admin/export", adminMiddleware, func(c *fiber.Ctx) error {
        var answers []Answer
        if err := db.Find(&answers).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data"})
        }
        var buf bytes.Buffer
        writer := csv.NewWriter(&buf)
        writer.Write([]string{"ID", "ParticipantID", "QuestionID", "AnswerText", "SubmittedAt", "IsDraft"})
        for _, a := range answers {
            writer.Write([]string{
                strconv.Itoa(int(a.ID)),
                strconv.Itoa(int(a.ParticipantID)),
                strconv.Itoa(int(a.QuestionID)),
                a.AnswerText,
                a.SubmittedAt.Format("2006-01-02 15:04:05"),
                strconv.FormatBool(a.IsDraft),
            })
        }
        writer.Flush()
        c.Set("Content-Type", "text/csv")
        c.Set("Content-Disposition", "attachment;filename=hasil_ujian.csv")
        return c.Send(buf.Bytes())
    })

    // Auto-save timer
    go func() {
        for {
            time.Sleep(30 * time.Second)
            fmt.Println("Auto-save triggered")
            // Implement auto-save logic here
        }
    }()

    // Start server
    log.Fatal(app.Listen(":3000"))
}