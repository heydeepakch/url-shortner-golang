package database

import (
    "database/sql"
    "log"
    
    _ "github.com/lib/pq"
)

var DB *sql.DB


func InitPostgres(connectionString string) error {
    var err error
    
    
    DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        return err
    }
    
    
    if err = DB.Ping(); err != nil {
        return err
    }
    
    
    DB.SetMaxOpenConns(25)
    DB.SetMaxIdleConns(5)
    
    log.Println("PostgreSQL connected successfully")
    
  
    if err := createTables(); err != nil {
        return err
    }
    
    return nil
}

func createTables() error {
 
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`
    
    // URLs table
    urlTable := `
    CREATE TABLE IF NOT EXISTS urls (
        id SERIAL PRIMARY KEY,
        short_code VARCHAR(10) UNIQUE NOT NULL,
        original_url TEXT NOT NULL,
        user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
        clicks BIGINT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        expires_at TIMESTAMP
    );
    
    CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
    CREATE INDEX IF NOT EXISTS idx_user_id ON urls(user_id);`
    
 
    if _, err := DB.Exec(userTable); err != nil {
        return err
    }
    
    if _, err := DB.Exec(urlTable); err != nil {
        return err
    }
    
    log.Println("Database tables created/verified")
    return nil
}


func CloseDB() {
    if DB != nil {
        DB.Close()
        log.Println("Database connection closed")
    }
}
