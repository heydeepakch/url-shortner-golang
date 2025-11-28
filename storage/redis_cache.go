package storage

import (
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/heydeepakch/url-shortner-golang/database"
    "github.com/heydeepakch/url-shortner-golang/models"
)	

const urlCacheTTL = 24 * time.Hour

// CacheURL stores URL in Redis cache
func CacheURL(url *models.URL) error {
    key := fmt.Sprintf("url:%s", url.ShortCode)
    
    data, err := json.Marshal(url)
    if err != nil {
        return err
    }
    
    return database.RedisClient.Set(database.Ctx, key, data, urlCacheTTL).Err()
}

// GetCachedURL retrieves URL from Redis cache
func GetCachedURL(shortCode string) (*models.URL, error) {
    key := fmt.Sprintf("url:%s", shortCode)
    
    data, err := database.RedisClient.Get(database.Ctx, key).Bytes()
    if err != nil {
        return nil, err // Returns redis.Nil if not found
    }
    
    var url models.URL
    if err := json.Unmarshal(data, &url); err != nil {
        return nil, err
    }
    
    return &url, nil
}

// DeleteCachedURL removes URL from cache
func DeleteCachedURL(shortCode string) error {
    key := fmt.Sprintf("url:%s", shortCode)
    return database.RedisClient.Del(database.Ctx, key).Err()
}

// IncrementClicksCache increments clicks in cache
func IncrementClicksCache(shortCode string) error {
    // key := fmt.Sprintf("url:%s", shortCode)
    
    
    url, err := GetCachedURL(shortCode)
    if err != nil {
        return err
    }
    
    url.Clicks++
    return CacheURL(url)
}
