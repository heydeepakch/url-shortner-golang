package storage

import (
    "database/sql"
    "errors"
    "time"
    
    "github.com/heydeepakch/url-shortner-golang/database"
    "github.com/heydeepakch/url-shortner-golang/models"
)


func CreateUser(user *models.User) error {
    query := `
        INSERT INTO users (username, email, password_hash, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    err := database.DB.QueryRow(
        query,
        user.Username,
        user.Email,
        user.PasswordHash,
        time.Now(),
    ).Scan(&user.ID)
    
    return err
}


func GetUserByEmail(email string) (*models.User, error) {
    query := `
        SELECT id, username, email, password_hash, created_at
        FROM users
        WHERE email = $1
    `
    
    user := &models.User{}
    err := database.DB.QueryRow(query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    
    return user, err
}


func GetUserByID(id int) (*models.User, error) {
    query := `
        SELECT id, username, email, password_hash, created_at
        FROM users
        WHERE id = $1
    `
    
    user := &models.User{}
    err := database.DB.QueryRow(query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    
    return user, err
}


func CreateURL(url *models.URL) error {
    query := `
        INSERT INTO urls (short_code, original_url, user_id, clicks, created_at, expires_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    
    err := database.DB.QueryRow(
        query,
        url.ShortCode,
        url.OriginalURL,
        url.UserID,
        0,
        time.Now(),
        url.ExpiresAt,
    ).Scan(&url.ID)
    
    return err
}

func GetURLByShortCode(shortCode string) (*models.URL, error) {
    query := `
        SELECT id, short_code, original_url, user_id, clicks, created_at, expires_at
        FROM urls
        WHERE short_code = $1 
        AND (expires_at IS NULL OR expires_at > NOW())
    `
    
    url := &models.URL{}
    err := database.DB.QueryRow(query, shortCode).Scan(
        &url.ID,
        &url.ShortCode,
        &url.OriginalURL,
        &url.UserID,
        &url.Clicks,
        &url.CreatedAt,
        &url.ExpiresAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, errors.New("URL not found or expired")
    }
    
    return url, err
}

func IncrementClicks(shortCode string) error {
    query := `UPDATE urls SET clicks = clicks + 1 WHERE short_code = $1`
    _, err := database.DB.Exec(query, shortCode)
    return err
}

func GetUserURLs(userID int) ([]models.URL, error) {
    query := `
        SELECT id, short_code, original_url, user_id, clicks, created_at, expires_at
        FROM urls
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
    
    rows, err := database.DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var urls []models.URL
    for rows.Next() {
        var url models.URL
        err := rows.Scan(
            &url.ID,
            &url.ShortCode,
            &url.OriginalURL,
            &url.UserID,
            &url.Clicks,
            &url.CreatedAt,
            &url.ExpiresAt,
        )
        if err != nil {
            return nil, err
        }
        urls = append(urls, url)
    }
    
    return urls, nil
}

func ShortCodeExists(shortCode string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)`
    var exists bool
    err := database.DB.QueryRow(query, shortCode).Scan(&exists)
    return exists, err
}
