package main

import (
	"os"

	"github.com/HoanNghi16/rentix-services/propety-service/internal/database"
	"github.com/HoanNghi16/rentix-services/propety-service/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PROPETY_PORT")

	db := database.ConnectDB()

	r := router.SetupRoutes(db)
	r.Run(":"+port)
}
