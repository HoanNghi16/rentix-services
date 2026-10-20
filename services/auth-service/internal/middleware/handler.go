package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)


type ErrorCode struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var ERROR_CODE = map[string]ErrorCode{
	"EMAIL_ALREADY_USED":{
		Status: 409,
		Message: "Email đã tồn tại",
	},
	"WRONG_PASSWORD": {
		Status:  400,
		Message: "Sai mật khẩu",
	},
	"USER_NOT_FOUND":{
		Status: 404,
		Message: "Không tìm thấy người dùng",
	},
}


func ErrorHandler(cntx *gin.Context){
	cntx.Next()
	log.Print("Đã tới đây")
	err_code, ok := cntx.Get("ERROR_CODE")
	if (!ok){
		return
	}
	err := ERROR_CODE[err_code.(string)]
	cntx.JSON(err.Status, gin.H{
		"error_code": err_code,
		"message": err.Message,
	})

}