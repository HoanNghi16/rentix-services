package middleware

import "github.com/gin-gonic/gin"

// error Handler here

func ErrorHandler(cntx *gin.Context){
	cntx.Next()

	if cntx.Errors !=nil{
		cntx.JSON(500, cntx.Errors)
		return
	}
}