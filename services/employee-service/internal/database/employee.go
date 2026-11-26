package database

import (
	"os"

	"github.com/HoanNghi16/rentix-services/employee-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	db_url := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(db_url))

	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Employee{}, &model.EmployeeAssignment{}, &model.EmployeeRole{}); err != nil {
		return nil, err
	}

	return db, nil
}
