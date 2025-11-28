package utils

import (
    "errors"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "github.com/heydeepakch/url-shortner-golang/config"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for user
func GenerateJWT(userID int, username, email string) (string, error) {
    expirationTime := time.Now().Add(
        time.Duration(config.AppConfig.JWTExpiryHours) * time.Hour,
    )
    
    claims := &Claims{
        UserID:   userID,
        Username: username,
        Email:    email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "url-shortener",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
    
    return tokenString, err
}

// ValidateJWT validates and parses JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return []byte(config.AppConfig.JWTSecret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return claims, nil
}
