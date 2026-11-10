package handler

import (
	"errors"
	"strconv"

	"github.com/HoanNghi16/rentix-services/propety-service/internal/model"
	"github.com/HoanNghi16/rentix-services/propety-service/internal/service"
	"github.com/gin-gonic/gin"
)


func GetUserClaims (cntx *gin.Context)(*model.UserClaims, error){
	claims, clm_exst := cntx.Get("User")
	
	if (!clm_exst){
		return nil, errors.New("INVALID_USER")
	}
	
	userClaims,ok := claims.(model.UserClaims)
	
	if (!ok){
		return nil, errors.New("PARSE_CLAIMS_FAILED")
	}
	return &userClaims, nil
}


type PropHandler struct {
	service *service.PropService
}

func NewHandler(service *service.PropService) *PropHandler{
	return &PropHandler{
		service,
	}
}

func (handler *PropHandler) GetComplexesHandler(cntx *gin.Context){
	userClaims, clm_err := GetUserClaims(cntx)

	if clm_err != nil{
		cntx.Error(clm_err)
		return
	}

	if userClaims.Role != "owner"{
		// Check permission
	}

	complexes, err := handler.service.GetComplexes(userClaims.UserID)

	if err != nil{
		cntx.Error(err)
		return
	}

	cntx.JSON(200, complexes)
}

func (handler *PropHandler) CreateComplexHandler(cntx *gin.Context){

	userClaims, clm_err := GetUserClaims(cntx)
	if clm_err != nil{
		cntx.Error(clm_err)
		cntx.Next()
		return
	}
	

	if userClaims.Role != "owner"{
		cntx.AbortWithStatusJSON(403, gin.H{
			"message": "Bạn không có quyền này!",
		})
	}


	var input *model.Complex

	if prs_err := cntx.ShouldBindBodyWithJSON(&input); prs_err != nil{
		cntx.Status(400)
		return
	}
	

	cmplx, crt_err := handler.service.NewComplex(userClaims.UserID, input)

	if (crt_err != nil){
		
		cntx.Error(crt_err)
		cntx.Next()
		return
	}

	cntx.JSON(201, cmplx)

}


func (handler *PropHandler) CreateRoomHandler(cntx *gin.Context){
	value, ok := cntx.Params.Get("complex_id")
	
	if !ok{
		cntx.Status(404)
		return
	}

	complexID, prs_err := strconv.ParseUint(value, 10, 10)

	if prs_err != nil{
		cntx.Error(prs_err)
		return
	}

	var input *model.Room

	if bnd_err := cntx.ShouldBindBodyWithJSON(&input); bnd_err != nil{
		cntx.Error(bnd_err)
		return
	}

	newRoom, crt_err := handler.service.NewRoom(uint(complexID), input)

	if (crt_err != nil){
		cntx.Error(crt_err)
		return
	}

	resRoom := model.RoomResponse{
		ID:	newRoom.ID,
		ComplexID: newRoom.ComplexID,
		RoomNumber: newRoom.RoomNumber,
		Status: newRoom.Status,
		Price: newRoom.Price,
	}

	cntx.JSON(201, resRoom)

}