package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/Baalamurgan/coin-selling-backend/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB = nil

func GetDB() *gorm.DB {
	if db != nil {
		return db
	}
	db = Connect()
	return db
}

func Connect() *gorm.DB {
	log.Println("Connecting to the database", config.DB_URI)
	sqlDB, err := sql.Open("postgres", config.DB_URI)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	log.Println("Connected to the database!")

	return db
}
