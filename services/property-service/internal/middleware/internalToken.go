package middleware

import (
	"strconv"

	"github.com/HoanNghi16/rentix-services/propety-service/internal/model"
	"github.com/gin-gonic/gin"
)

func InternalTokenHandler(cntx *gin.Context) {

	// InternalToken here before create UserClaims


	userID, prs_err :=  strconv.ParseUint(
		cntx.GetHeader("X-User-ID"), 10, 10,
	)

	if (prs_err!=nil){
		cntx.JSON(400, gin.H{
			"message": "UserID không hợp lệ",
		})
	}

	userRole := cntx.GetHeader("X-User-Role")
	user := model.UserClaims{
		UserID: uint(userID),
		Role: userRole,
	}

	cntx.Set("User", user)
	cntx.Next()
}