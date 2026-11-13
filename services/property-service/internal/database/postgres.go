package database

import (
	"log"
	"os"

	"github.com/HoanNghi16/rentix-services/property-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	db_url := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(db_url))
	if err != nil {
		panic("connect failed" + err.Error())
	}

	log.Print("Connect successfully")

	if err := db.AutoMigrate(
		&model.Complex{},
		&model.Room{},
		&model.EmployeeAssignment{},
		&model.Outbox{},
	); err != nil {
		panic("database migration failed: " + err.Error())
	}

	return db
}
