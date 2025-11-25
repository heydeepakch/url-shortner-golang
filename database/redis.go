package database

import (
    "context"
    "log"
    "time"
    
    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()


func InitRedis(addr, password string, db int) error {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:         addr,
        Password:     password,
        DB:           db,
        PoolSize:     10,
        MinIdleConns: 5,
        MaxRetries:   3,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })
    
 
    _, err := RedisClient.Ping(Ctx).Result()
    if err != nil {
        return err
    }
    
    log.Println("Redis connected successfully")
    return nil
}


func CloseRedis() {
    if RedisClient != nil {
        RedisClient.Close()
        log.Println("Redis connection closed")
    }
}
