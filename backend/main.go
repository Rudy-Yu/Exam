package main

import (
    "fmt"
    "log"
    "time"
    "encoding/csv"
    "strconv"
    "bytes"
    "context"
    "encoding/json"
    "sync"
    "crypto/tls"
    "./utils"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/go-redis/redis/v8"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/golang-jwt/jwt/v4"
    "github.com/jackc/pgx/v4/pgxpool"
)

// Global variables
var (
    rdb *redis.Client
    store *session.Store
    examTimers sync.Map
    ctx = context.Background()
    config *utils.Config
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
    Weight        int    `json:"weight" gorm:"default:1"` // Bobot soal
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
    // Setup connection pool
    pgConfig, err := pgxpool.ParseConfig(config.GetDBConnString())
    if err != nil {
        log.Fatal("Error parsing config:", err)
    }

    // Set max connections
    pgConfig.MaxConns = 20

    // Create pool
    pool, err := pgxpool.ConnectConfig(context.Background(), pgConfig)
    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }

    // Convert pool to *sql.DB
    db, err := gorm.Open(postgres.New(postgres.Config{
        Conn: pool,
    }), &gorm.Config{})

    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    return db
}

// Redis connection
func connectRedis() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     config.GetRedisAddr(),
        Password: config.RedisPass,
        DB:       0,
        PoolSize: 100,
    })
    return rdb
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

    // Connect to Redis
    rdb = connectRedis()
    defer rdb.Close()

    // Initialize session store with Redis
    store = session.New(session.Config{
        Storage: rdb,
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
        err := rdb.Set(ctx, sessionKey, sessionData, time.Hour).Err()
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
        err := rdb.Set(ctx, key, answer.AnswerText, time.Hour).Err()
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
                rdb.Del(ctx, key)
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
        sessionData, err := rdb.Get(ctx, sessionKey).Result()
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Sesi ujian tidak ditemukan",
            })
        }
        
        var session ExamSession
        json.Unmarshal([]byte(sessionData), &session)
        
        remainingTime := session.Duration - int(time.Since(session.StartTime).Seconds())
        if remainingTime < 0 {
            remainingTime = 0
        }
        
        return c.JSON(fiber.Map{
            "success": true,
            "remaining_time": remainingTime,
        })
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