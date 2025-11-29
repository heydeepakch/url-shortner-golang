package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    
    "github.com/heydeepakch/url-shortner-golang/models"
    "github.com/heydeepakch/url-shortner-golang/storage"
    "github.com/heydeepakch/url-shortner-golang/utils"
)

func Register(c *gin.Context) {
    var req models.RegisterRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data",
            "details": err.Error(),
        })
        return
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(req.Password),
        bcrypt.DefaultCost,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to process password",
        })
        return
    }
    
    user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
    }
    
    if err := storage.CreateUser(user); err != nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": "Username or email already exists",
        })
        return
    }
    
    token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate token",
        })
        return
    }
    
    c.JSON(http.StatusCreated, models.LoginResponse{
        Token: token,
        User:  *user,
    })
}

func Login(c *gin.Context) {
    var req models.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data",
        })
        return
    }
    
    user, err := storage.GetUserByEmail(req.Email)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid email or password",
        })
        return
    }
    
    err = bcrypt.CompareHashAndPassword(
        []byte(user.PasswordHash),
        []byte(req.Password),
    )
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid email or password",
        })
        return
    }
    
    token, err := utils.GenerateJWT(user.ID, user.Username, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to generate token",
        })
        return
    }
    
    c.JSON(http.StatusOK, models.LoginResponse{
        Token: token,
        User:  *user,
    })
}
			
func GetProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Unauthorized",
        })
        return
    }
    
    user, err := storage.GetUserByID(userID.(int))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "User not found",
        })
        return
    }
    
    c.JSON(http.StatusOK, user)
}
