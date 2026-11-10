package router

import (
	"github.com/HoanNghi16/rentix-services/propety-service/internal/handler"
	"github.com/HoanNghi16/rentix-services/propety-service/internal/middleware"
	"github.com/HoanNghi16/rentix-services/propety-service/internal/repository"
	"github.com/HoanNghi16/rentix-services/propety-service/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Layer
	repo := repository.NewRepo(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)


	// Middlewares
	router.Use(middleware.InternalTokenHandler)
	router.Use(middleware.ErrorHandler)

	// endpoints
	// complex
	Complex := router.Group("/complex")
	{
		Complex.POST("/", handler.CreateComplexHandler)
		Complex.GET("/", handler.GetComplexesHandler)
		Complex.GET("/:complex_id") // Lấy complex theo complex_id
	}
	Room := Complex.Group("/:complex_id/room")

	// Room
	// --- Chưa làm
	Room.GET("/")         //Danh sách phòng
	Room.GET("/:room_id") //Chi tiết phòng
	Room.POST("/", handler.CreateRoomHandler) 	      //Tạo phòng
	

	return router
}