package main

import (
    "fmt"
    "log"
    "time"
    "context"
    "encoding/json"
    "sync"
    "os"
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/gofiber/storage/redis"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v4"
    "online-exam-app-backend/utils"
)

// Global variables
var (
    store *session.Store
    examTimers sync.Map
    ctx = context.Background()
    config *utils.Config
)

// User model
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `json:"name"`
    Email    string `gorm:"unique"`
    Password string
    Role     string `gorm:"default:user"`
}

// Question model
type Question struct {
    ID            uint     `gorm:"primaryKey" json:"id"`
    ExamID        uint     `json:"exam_id"`
    QuestionText  string   `json:"question_text"`
    CorrectAnswer string   `json:"correct_answer"`
    Weight        int      `json:"weight" gorm:"default:1"`
    Type          string   `json:"type" gorm:"default:'pilihan_ganda'"`
    Options       []string `gorm:"type:json" json:"options"`
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

// ExamSession model untuk Redis
type ExamSession struct {
    UserID    uint      `json:"user_id"`
    ExamID    uint      `json:"exam_id"`
    StartTime time.Time `json:"start_time"`
    Duration  int       `json:"duration"` // dalam detik
}

// Database connection dengan connection pooling
func connectDB() *gorm.DB {
    dsn := config.GetDBConnString()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    return db
}

// Middleware untuk autentikasi JWT
func authMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "message": "Token tidak ditemukan",
        })
    }

    tokenString := authHeader[len("Bearer "):]
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })

    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "success": false,
            "message": "Token tidak valid",
        })
    }

    claims := token.Claims.(jwt.MapClaims)
    c.Locals("user_id", claims["user_id"])
    c.Locals("role", claims["role"])
    return c.Next()
}

// Middleware untuk admin
func adminMiddleware(c *fiber.Ctx) error {
    if c.Locals("role") != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "success": false,
            "message": "Akses admin diperlukan",
        })
    }
    return c.Next()
}

func main() {
    // Load configuration
    config = utils.LoadConfig()

    // Connect to PostgreSQL with connection pooling
    db := connectDB()
    db.AutoMigrate(&User{}, &Question{}, &Answer{})

    // Initialize session store with Redis (Fiber Storage)
    store = session.New(session.Config{
        Storage: redis.New(redis.Config{
            Host:     config.RedisHost,
            Port:     getRedisPort(config.RedisPort), // konversi string ke int
            Password: config.RedisPass,
            Database: 0,
        }),
        Expiration: 24 * time.Hour,
    })

    // Initialize Fiber app with custom config
    app := fiber.New(fiber.Config{
        Prefork: true, // Enable untuk multi-core processing
    })

    // Add CORS middleware
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // ========== AUTH & USER ENDPOINTS ==========
    // Register endpoint
    app.Post("/api/register", func(c *fiber.Ctx) error {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Format data tidak valid",
            })
        }
        if req.Email == "" || req.Password == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Email dan password wajib diisi",
            })
        }
        // Cek apakah user sudah ada
        var user User
        if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "success": false,
                "message": "Email sudah terdaftar",
            })
        }
        // Simpan user baru
        user = User{Email: req.Email, Password: req.Password}
        if err := db.Create(&user).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal mendaftar user",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Pendaftaran berhasil",
        })
    })

    // Login endpoint
    app.Post("/api/login", func(c *fiber.Ctx) error {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Format data tidak valid",
            })
        }
        var user User
        if err := db.Where("email = ? AND password = ?", req.Email, req.Password).First(&user).Error; err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "success": false,
                "message": "Email atau password salah",
            })
        }
        // Generate JWT
        claims := jwt.MapClaims{
            "user_id": user.ID,
            "role": user.Role,
            "exp": time.Now().Add(time.Hour * 24).Unix(),
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString([]byte("secret"))
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal membuat token",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "token": tokenString,
            "user": fiber.Map{
                "id": user.ID,
                "email": user.Email,
                "role": user.Role,
            },
        })
    })

    // ========== EXAM ENDPOINTS ==========
    // List all exams
    app.Get("/api/exams", authMiddleware, func(c *fiber.Ctx) error {
        var exams []struct {
            ID       uint   `json:"id"`
            Title    string `json:"title"`
            Duration int    `json:"duration"`
        }
        // Ambil dari tabel Exam jika ada
        type ExamDB struct {
            ID       uint
            Title    string
            Duration int
        }
        var dbExams []ExamDB
        if err := db.Table("exams").Select("id, title, duration").Scan(&dbExams).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal mengambil data ujian",
            })
        }
        for _, e := range dbExams {
            exams = append(exams, struct {
                ID       uint   `json:"id"`
                Title    string `json:"title"`
                Duration int    `json:"duration"`
            }{e.ID, e.Title, e.Duration})
        }
        return c.JSON(exams)
    })

    // Get questions for an exam
    app.Get("/api/exam/:id/questions", authMiddleware, func(c *fiber.Ctx) error {
        examID := c.Params("id")
        var questions []struct {
            ID            uint     `json:"id"`
            QuestionText  string   `json:"question_text"`
            Options       []string `json:"options"`
            CorrectAnswer string   `json:"correct_answer"`
            Weight        int      `json:"weight"`
        }
        // Ambil dari tabel questions
        type QuestionDB struct {
            ID            uint
            ExamID        uint
            QuestionText  string
            CorrectAnswer string
            Weight        int
            // Asumsi ada kolom options bertipe text/json di DB, jika tidak, dummy
        }
        var dbQuestions []QuestionDB
        if err := db.Table("questions").Where("exam_id = ?", examID).Find(&dbQuestions).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal mengambil soal",
            })
        }
        for _, q := range dbQuestions {
            // Dummy options jika tidak ada di DB
            options := []string{"A", "B", "C", "D"}
            questions = append(questions, struct {
                ID            uint     `json:"id"`
                QuestionText  string   `json:"question_text"`
                Options       []string `json:"options"`
                CorrectAnswer string   `json:"correct_answer"`
                Weight        int      `json:"weight"`
            }{q.ID, q.QuestionText, options, q.CorrectAnswer, q.Weight})
        }
        return c.JSON(questions)
    })

    // Session handling endpoint
    app.Post("/api/exam/:id/start", authMiddleware, func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(float64)
        examID, _ := strconv.Atoi(c.Params("id"))
        
        // Create exam session
        session := ExamSession{
            UserID:    uint(userID),
            ExamID:    uint(examID),
            StartTime: time.Now(),
            Duration:  3600, // 1 hour default
        }
        
        // Store in Redis
        sessionKey := fmt.Sprintf("exam_session:%d:%d", uint(userID), examID)
        sessionData, _ := json.Marshal(session)
        err := store.Storage.Set(sessionKey, sessionData, time.Hour)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal memulai ujian",
            })
        }
        
        return c.JSON(fiber.Map{
            "success": true,
            "start_time": session.StartTime,
            "duration": session.Duration,
        })
    })

    // Auto-save answer endpoint
    app.Post("/api/answers/draft", authMiddleware, func(c *fiber.Ctx) error {
        var answer struct {
            QuestionID uint   `json:"question_id"`
            AnswerText string `json:"answer_text"`
        }
        
        if err := c.BodyParser(&answer); err != nil {
            return err
        }
        
        userID := c.Locals("user_id").(float64)
        
        // Store draft answer in Redis
        key := fmt.Sprintf("draft_answer:%d:%d", uint(userID), answer.QuestionID)
        err := store.Storage.Set(key, []byte(answer.AnswerText), time.Hour)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menyimpan jawaban sementara",
            })
        }
        
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Jawaban tersimpan sementara",
        })
    })

    // Submit final answer endpoint with goroutine
    app.Post("/api/answers/submit", authMiddleware, func(c *fiber.Ctx) error {
        var answers []Answer
        if err := c.BodyParser(&answers); err != nil {
            return err
        }
        
        userID := c.Locals("user_id").(float64)
        
        // Process answers in goroutine
        go func() {
            for _, answer := range answers {
                answer.ParticipantID = uint(userID)
                answer.SubmittedAt = time.Now()
                answer.IsDraft = false
                db.Create(&answer)
                
                // Clean up Redis draft
                key := fmt.Sprintf("draft_answer:%d:%d", uint(userID), answer.QuestionID)
                store.Storage.Delete(key)
            }
        }()
        
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Jawaban berhasil disubmit",
        })
    })

    // Get exam timer endpoint
    app.Get("/api/exam/:id/timer", authMiddleware, func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(float64)
        examID := c.Params("id")
        
        // Get session from Redis
        sessionKey := fmt.Sprintf("exam_session:%d:%s", uint(userID), examID)
        sessionData, err := store.Storage.Get(sessionKey)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Sesi ujian tidak ditemukan",
            })
        }
        
        var session ExamSession
        if err := json.Unmarshal(sessionData, &session); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal membaca data sesi",
            })
        }
        
        remainingTime := session.Duration - int(time.Since(session.StartTime).Seconds())
        if remainingTime < 0 {
            remainingTime = 0
        }
        
        return c.JSON(fiber.Map{
            "success": true,
            "remaining_time": remainingTime,
        })
    })

    // ========== ADMIN ENDPOINTS ==========
    admin := app.Group("/api/admin", authMiddleware, adminMiddleware)

    // List all users
    admin.Get("/users", func(c *fiber.Ctx) error {
        var users []User
        if err := db.Find(&users).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal mengambil user",
            })
        }
        return c.JSON(users)
    })

    // Add user
    admin.Post("/users", func(c *fiber.Ctx) error {
        var req struct {
            Name     string `json:"name"`
            Email    string `json:"email"`
            Password string `json:"password"`
            Role     string `json:"role"`
        }
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Format data tidak valid",
            })
        }
        if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Nama, email, password, dan role wajib diisi",
            })
        }
        // Cek apakah user sudah ada
        var user User
        if err := db.Where("email = ?", req.Email).First(&user).Error; err == nil {
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "success": false,
                "message": "Email sudah terdaftar",
            })
        }
        user = User{Name: req.Name, Email: req.Email, Password: req.Password, Role: req.Role}
        if err := db.Create(&user).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menambah user",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "User berhasil ditambah",
        })
    })

    // Edit user
    admin.Put("/users/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
        var user User
        if err := db.First(&user, id).Error; err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "User tidak ditemukan",
            })
        }
        var req struct {
            Name     string `json:"name"`
            Email    string `json:"email"`
            Password string `json:"password"`
            Role     string `json:"role"`
        }
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Format data tidak valid",
            })
        }
        user.Name = req.Name
        user.Email = req.Email
        user.Password = req.Password
        user.Role = req.Role
        if err := db.Save(&user).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal update user",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "User berhasil diupdate",
        })
    })

    // Delete user
    admin.Delete("/users/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
        if err := db.Delete(&User{}, id).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menghapus user",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "User berhasil dihapus",
        })
    })

    // List all questions
    admin.Get("/questions", func(c *fiber.Ctx) error {
        var questions []Question
        if err := db.Find(&questions).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal mengambil soal",
            })
        }
        return c.JSON(questions)
    })

    // Add question
    admin.Post("/questions", func(c *fiber.Ctx) error {
        var q Question
        if err := c.BodyParser(&q); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Format data tidak valid",
            })
        }
        if q.Type == "" {
            q.Type = "pilihan_ganda"
        }
        if q.Type == "pilihan_ganda" && len(q.Options) < 2 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Opsi pilihan ganda minimal 2!",
            })
        }
        if err := db.Create(&q).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menambah soal",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Soal berhasil ditambah",
        })
    })

    // Edit question
    admin.Put("/questions/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
        var q Question
        if err := db.First(&q, id).Error; err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Soal tidak ditemukan",
            })
        }
        var update Question
        if err := c.BodyParser(&update); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Format data tidak valid",
            })
        }
        q.ExamID = update.ExamID
        q.QuestionText = update.QuestionText
        q.CorrectAnswer = update.CorrectAnswer
        q.Weight = update.Weight
        q.Type = update.Type
        q.Options = update.Options
        if q.Type == "" {
            q.Type = "pilihan_ganda"
        }
        if q.Type == "pilihan_ganda" && len(q.Options) < 2 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Opsi pilihan ganda minimal 2!",
            })
        }
        if err := db.Save(&q).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal update soal",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Soal berhasil diupdate",
        })
    })

    // Delete question
    admin.Delete("/questions/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
        if err := db.Delete(&Question{}, id).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "success": false,
                "message": "Gagal menghapus soal",
            })
        }
        return c.JSON(fiber.Map{
            "success": true,
            "message": "Soal berhasil dihapus",
        })
    })

    // Export hasil ujian (CSV)
    admin.Get("/export", func(c *fiber.Ctx) error {
        var answers []Answer
        if err := db.Find(&answers).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Gagal mengambil data jawaban")
        }
        csv := "ID,ParticipantID,QuestionID,AnswerText,SubmittedAt,IsDraft\n"
        for _, a := range answers {
            csv += fmt.Sprintf("%d,%d,%d,%s,%s,%t\n", a.ID, a.ParticipantID, a.QuestionID, a.AnswerText, a.SubmittedAt.Format(time.RFC3339), a.IsDraft)
        }
        c.Set("Content-Type", "text/csv")
        c.Set("Content-Disposition", "attachment; filename=hasil_ujian.csv")
        return c.SendString(csv)
    })

    // Development/Production mode switch
    if os.Getenv("GO_ENV") == "production" {
        // Production mode with SSL
        // Start HTTP server (for redirect)
        go func() {
            redirectApp := fiber.New()
            redirectApp.Use(func(c *fiber.Ctx) error {
                return c.Redirect("https://" + c.Hostname() + c.OriginalURL())
            })
            log.Fatal(redirectApp.Listen(":80"))
        }()

        // Start HTTPS server
        log.Printf("Server is running in production mode (HTTPS) on port 443")
        log.Fatal(app.ListenTLS(":443", config.SSLCert, config.SSLKey))
    } else {
        // Development mode without SSL
        log.Printf("Server is running in development mode (HTTP) on port 3000")
        log.Fatal(app.Listen(":3000"))
    }
}

// Tambahkan fungsi bantu di bawah main()
func getRedisPort(portStr string) int {
    port := 6379
    if p, err := strconv.Atoi(portStr); err == nil {
        port = p
    }
    return port
}