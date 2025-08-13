package internal

import (
	_ "database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB return instance of db
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("[DB] Database not initialized. Call InitDB() first.")
	}
	return db
}

// InitDB init the conn to postgresql for singleton
func InitDB(host, user, password, dbname string, port int) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			host, user, password, dbname, port,
		)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
		if err != nil {
			log.Fatalf("[DB] Failed to connect to database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("[DB] Failed to get SQL DB from GORM: %v", err)
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)

		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("[DB] Ping failed: %v", err)
		}

		log.Println("[DB] Connected to database ✔")
	})
}
