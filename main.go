package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/gin-gonic/gin"
    
    "github.com/heydeepakch/url-shortner-golang/config"
    "github.com/heydeepakch/url-shortner-golang/database"
    "github.com/heydeepakch/url-shortner-golang/handlers"
    "github.com/heydeepakch/url-shortner-golang/middleware"
)

func main() {

    config.LoadConfig()
    
    if err := database.InitPostgres(config.AppConfig.DatabaseURL); err != nil {
        log.Fatal("Failed to connect to PostgreSQL:", err)
    }
    defer database.CloseDB()
    
    if err := database.InitRedis(
        config.AppConfig.RedisAddr,
        config.AppConfig.RedisPassword,
        config.AppConfig.RedisDB,
    ); err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    defer database.CloseRedis()
    
    router := gin.Default()
    
    router.Use(corsMiddleware())
    
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    
    router.POST("/api/register", handlers.Register)
    router.POST("/api/login", handlers.Login)

    router.POST("/api/shorten", middleware.OptionalAuthMiddleware(), handlers.ShortenURL)
    
    router.GET("/:code", handlers.RedirectURL)
    
    router.GET("/api/url/:code/stats", middleware.OptionalAuthMiddleware(), handlers.GetURLStats)
    
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", handlers.GetProfile)
        protected.GET("/my-urls", handlers.GetMyURLs)
    }
    
    go func() {
        if err := router.Run(":" + config.AppConfig.Port); err != nil {
            log.Fatal("Failed to start server:", err)
        }
    }()
    
    log.Printf("Server started on port %s", config.AppConfig.Port)
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
