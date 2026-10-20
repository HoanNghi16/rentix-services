package main

import (
	"os"

	"github.com/HoanNghi16/rentix-services/auth-service/internal/database"
	"github.com/HoanNghi16/rentix-services/auth-service/router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := database.ConnectDB()
	port := os.Getenv("AUTH_PORT")

	if err != nil{
		panic("Cannot connect database: " + err.Error())
	}

	ls := router.SetupRoute(db)
	ls.Run(":"+port)

}