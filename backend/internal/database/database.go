package database

import (
	"log"
	"time"

	"backend/internal/config"
	"backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	logMode := logger.Info
	if cfg.GinMode == "release" {
		logMode = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pooling
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Run auto migrations
	log.Println("Running database auto-migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.URL{},
		&models.URLClick{},
	)
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("Database connection & migration completed successfully.")
	return db, nil
}
