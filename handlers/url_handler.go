package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    
    "github.com/heydeepakch/url-shortner-golang/config"
    "github.com/heydeepakch/url-shortner-golang/models"
    "github.com/heydeepakch/url-shortner-golang/storage"
    "github.com/heydeepakch/url-shortner-golang/utils"
)

func ShortenURL(c *gin.Context) {
    var req models.ShortenRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data",
            "details": err.Error(),
        })
        return
    }
    
    var shortCode string
    var err error
    
    if req.CustomCode != "" {
        if !utils.ValidateCustomCode(req.CustomCode) {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid custom code (4-20 alphanumeric characters only)",
            })
            return
        }
        
        exists, _ := storage.ShortCodeExists(req.CustomCode)
        if exists {
            c.JSON(http.StatusConflict, gin.H{
                "error": "Custom code already in use",
            })
            return
        }
        
        shortCode = req.CustomCode
    } else {
        shortCode, err = utils.GenerateShortCode()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to generate short code",
            })
            return
        }
    }
    
    var userID *int
    if id, exists := c.Get("user_id"); exists {
        uid := id.(int)
        userID = &uid
    }
    
    var expiresAt *time.Time
    if req.ExpiresInHrs > 0 {
        expiry := time.Now().Add(time.Duration(req.ExpiresInHrs) * time.Hour)
        expiresAt = &expiry
    }
    
    url := &models.URL{
        ShortCode:   shortCode,
        OriginalURL: req.URL,
        UserID:      userID,
        ExpiresAt:   expiresAt,
    }
    
    if err := storage.CreateURL(url); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create short URL",
        })
        return
    }
    
    storage.CacheURL(url)
    
    response := models.ShortenResponse{
        ShortCode:   shortCode,
        ShortURL:    config.AppConfig.BaseURL + "/" + shortCode,
        OriginalURL: req.URL,
        ExpiresAt:   expiresAt,
    }
    
    c.JSON(http.StatusCreated, response)
}

func RedirectURL(c *gin.Context) {
    shortCode := c.Param("code")
    
    url, err := storage.GetCachedURL(shortCode)
    
    if err == redis.Nil || url == nil {
        url, err = storage.GetURLByShortCode(shortCode)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "URL not found or expired",
            })
            return
        }
        
        // Update cache
        storage.CacheURL(url)
    }
    
    // Increment clicks asynchronously
    go storage.IncrementClicks(shortCode)
    go storage.IncrementClicksCache(shortCode)
    
    c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func GetURLStats(c *gin.Context) {
    shortCode := c.Param("code")
    
    url, err := storage.GetURLByShortCode(shortCode)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "URL not found",
        })
        return
    }
    
    if url.UserID != nil {
        if userID, exists := c.Get("user_id"); exists {
            if *url.UserID != userID.(int) {
                c.JSON(http.StatusForbidden, gin.H{
                    "error": "You don't have permission to view these stats",
                })
                return
            }
        }
    }
    
    response := models.URLStatsResponse{
        URL:      *url,
        ShortURL: config.AppConfig.BaseURL + "/" + shortCode,
    }
    
    c.JSON(http.StatusOK, response)
}

func GetMyURLs(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized",
        })
        return
    }
    
    urls, err := storage.GetUserURLs(userID.(int))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch URLs",
        })
        return
    }
    
    var response []models.URLStatsResponse
    for _, url := range urls {
        response = append(response, models.URLStatsResponse{
            URL:      url,
            ShortURL: config.AppConfig.BaseURL + "/" + url.ShortCode,
        })
    }
    
    c.JSON(http.StatusOK, gin.H{
        "count": len(response),
        "urls":  response,
    })
}
