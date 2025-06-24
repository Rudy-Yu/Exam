package utils

import (
    "fmt"
    "os"
)

type Config struct {
    // Database
    DBHost     string
    DBPort     string
    DBUser     string
    DBPass     string
    DBName     string
    
    // Redis
    RedisHost string
    RedisPort string
    RedisPass string
    
    // JWT
    JWTSecret string
    
    // SSL/TLS
    SSLCert string
    SSLKey  string
}

func LoadConfig() *Config {
    config := &Config{
        // Database
        DBHost: getEnv("DB_HOST", "localhost"),
        DBPort: getEnv("DB_PORT", "5433"),
        DBUser: getEnv("DB_USER", "postgres"),
        DBPass: getEnv("DB_PASS", "admin"),
        DBName: getEnv("DB_NAME", "examdb"),
        
        // Redis
        RedisHost: getEnv("REDIS_HOST", "localhost"),
        RedisPort: getEnv("REDIS_PORT", "6379"),
        RedisPass: getEnv("REDIS_PASS", ""),
        
        // JWT
        JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
        
        // SSL/TLS
        SSLCert: getEnv("SSL_CERT", "./cert.pem"),
        SSLKey:  getEnv("SSL_KEY", "./key.pem"),
    }
    
    return config
}

func (c *Config) GetDBConnString() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName,
    )
}

func (c *Config) GetRedisAddr() string {
    return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
} 