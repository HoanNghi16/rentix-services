package handler

import (
	"github.com/HoanNghi16/rentix-services/auth-service/internal/model"
	"github.com/HoanNghi16/rentix-services/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct{
	service *service.AuthService
}

func NewHandler(service *service.AuthService)*AuthHandler{
	return &AuthHandler{
		service,
	}
}

func (handler *AuthHandler) LoginHandler(cntx *gin.Context){
	var input model.LoginRequest
	if bd_err := cntx.ShouldBindBodyWithJSON(&input); bd_err != nil{
		cntx.JSON(400, gin.H{
			"error_code": "INVALID_INPUT",
		})
		return
	}

	result, lg_err := handler.service.Login(input.Email, input.Password)
	if lg_err != nil{
		cntx.Set("ERROR_CODE", lg_err.Error())
		return
	}
	cntx.JSON(200, result)
}


func (handler *AuthHandler) RegisterHandler(cntx *gin.Context){
	var input model.RegisterRequest
	cntx.ShouldBindBodyWithJSON(&input)
	
	if err := handler.service.Register(input.Email, input.Password); err != nil{
		cntx.Set("ERROR_CODE", err.Error())
		return
	}
	cntx.JSON(201, gin.H{
		"message": "Đăng ký thành công",
	})
}