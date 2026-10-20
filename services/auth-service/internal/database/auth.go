package database

import (
	"os"

	"github.com/HoanNghi16/rentix-services/auth-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil{
		return nil, err
	}

	db.AutoMigrate(model.User{})

	return db, nil
}