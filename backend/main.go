package main

import (
    "fmt"
    "log"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/go-redis/redis/v8"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// User model
type User struct {
    ID       uint   `gorm:"primaryKey"`
    Email    string `gorm:"unique"`
    Password string
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

func main() {
    // Connect to PostgreSQL
    db := connectDB()
    db.AutoMigrate(&User{}) // Creates the users table automatically

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

        return c.JSON(fiber.Map{"success": true, "message": "Login berhasil"})
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