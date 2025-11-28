package utils

import (
    "crypto/rand"
    "math/big"
    
    "github.com/heydeepakch/url-shortner-golang/storage"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const defaultLength = 7

// GenerateShortCode generates a unique random short code
func GenerateShortCode() (string, error) {
    return GenerateShortCodeWithLength(defaultLength)
}

// GenerateShortCodeWithLength generates short code of specific length
func GenerateShortCodeWithLength(length int) (string, error) {
    maxAttempts := 5
    
    for attempt := 0; attempt < maxAttempts; attempt++ {
        code, err := generateRandomCode(length)
        if err != nil {
            return "", err
        }
        
        // Check if code already exists
        exists, err := storage.ShortCodeExists(code)
        if err != nil {
            return "", err
        }
        
        if !exists {
            return code, nil
        }
    }
    
    // If collision after max attempts, increase length
    return GenerateShortCodeWithLength(length + 1)
}

// generateRandomCode creates a cryptographically secure random string
func generateRandomCode(length int) (string, error) {
    result := make([]byte, length)
    
    for i := range result {
        num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            return "", err
        }
        result[i] = charset[num.Int64()]
    }
    
    return string(result), nil
}

// ValidateCustomCode checks if custom code is valid
func ValidateCustomCode(code string) bool {
    if len(code) < 4 || len(code) > 20 {
        return false
    }
    
    // Check if contains only allowed characters
    for _, char := range code {
        valid := false
        for _, allowed := range charset {
            if char == allowed {
                valid = true
                break
            }
        }
        if !valid {
            return false
        }
    }
    
    return true
}
