package router

import (
	"github.com/HoanNghi16/rentix-services/auth-service/internal/handler"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/middleware"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/repository"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoute(db *gorm.DB) *gin.Engine {
	repo := repository.NewRepo(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	
	server := gin.Default()
	server.Use(middleware.ErrorHandler)
	server.POST("/login", handler.LoginHandler)
	server.POST("/register", handler.RegisterHandler)

	return server
}