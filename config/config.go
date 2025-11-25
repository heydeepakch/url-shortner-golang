package config

import (
    "log"
    "os"
    "strconv"
    
    "github.com/joho/godotenv"
)

type Config struct {
    DatabaseURL    string
    RedisAddr      string
    RedisPassword  string
    RedisDB        int
    Port           string
    BaseURL        string
    JWTSecret      string
    JWTExpiryHours int
}

var AppConfig *Config


func LoadConfig() {

    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
    
    redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
    jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
    
    AppConfig = &Config{
        DatabaseURL:    getEnv("DATABASE_URL", ""),
        RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPassword:  getEnv("REDIS_PASSWORD", ""),
        RedisDB:        redisDB,
        Port:           getEnv("PORT", "8080"),
        BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
        JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-me"),
        JWTExpiryHours: jwtExpiry,
    }
    
    log.Println("Configuration loaded successfully")
}


func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
