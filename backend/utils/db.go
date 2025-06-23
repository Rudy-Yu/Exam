package utils

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "log"
)

func ConnectDB() *gorm.DB {
    dsn := "host=localhost user=postgres password=password dbname=examdb port=5432 sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    return db
}