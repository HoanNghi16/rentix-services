package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/HoanNghi16/rentix-services/api-gateway/internal/middleware"
	"github.com/HoanNghi16/rentix-services/api-gateway/internal/propertyclient"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type propertyHandler struct{ property *propertyclient.Client }

type createComplexRequest struct {
	Name            string                   `json:"name" binding:"required"`
	Avatar          string                   `json:"avatar"`
	Address         string                   `json:"address"`
	EmployeeUserIDs []uint64                 `json:"employee_user_ids"`
	Services        *propertyclient.Services `json:"services"`
}

func RegisterPropertyRoutes(router gin.IRouter, property *propertyclient.Client, accessSecret string) {
	handler := &propertyHandler{property: property}
	propertyRoutes := router.Group("/property")
	propertyRoutes.Use(middleware.AuthMiddleware(accessSecret))
	propertyRoutes.GET("/complex", handler.getComplexes)
	propertyRoutes.GET("/complex/:complex_id/room", handler.getRooms)
	propertyRoutes.POST("/complex", handler.createComplex)
	// propertyRoutes.POST("/complex/:complex_id/room") 
}

func (handler *propertyHandler) getComplexes(c *gin.Context) {
	userID, ok := middleware.AuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	role, ok := middleware.AuthenticatedRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	response, err := handler.property.GetComplexes(c.Request.Context(), &propertyclient.GetComplexesRequest{UserID: userID, Role: role})
	if err != nil {
		writePropertyGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (handler *propertyHandler) getRooms(c *gin.Context) {
	complexID, err := strconv.ParseInt(c.Param("complex_id"), 10, 64)
	if err != nil || complexID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "complex_id must be a positive integer"})
		return
	}
	userID, ok := middleware.AuthenticatedUserID(c)
	if !ok || userID > uint64(^uint64(0)>>1) {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	role, ok := middleware.AuthenticatedRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	response, err := handler.property.GetRooms(c.Request.Context(), &propertyclient.GetRoomsRequest{
		ComplexID: complexID,
		UserID:    int64(userID),
		Role:      role,
	})

	jsondt, _ := json.Marshal(response)

	log.Print(string(jsondt))
	if err != nil {
		writePropertyGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (handler *propertyHandler) createComplex(c *gin.Context) {
	var input createComplexRequest
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "a complex name is required"})
		return
	}
	userID, ok := middleware.AuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	role, ok := middleware.AuthenticatedRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	response, err := handler.property.CreateComplex(c.Request.Context(), &propertyclient.CreateComplexRequest{
		UserID:          userID,
		Role:            role,
		Name:            strings.TrimSpace(input.Name),
		Avatar:          input.Avatar,
		Address:         input.Address,
		EmployeeUserIds: input.EmployeeUserIDs,
		Services:        input.Services,
	})
	if err != nil {
		writePropertyGRPCError(c, err)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func writePropertyGRPCError(c *gin.Context, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "request is invalid"})
	case codes.PermissionDenied:
		c.JSON(http.StatusForbidden, gin.H{"error_code": "FORBIDDEN", "message": "you are not allowed to perform this action"})
	case codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
	case codes.Unavailable, codes.DeadlineExceeded:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error_code": "PROPERTY_UNAVAILABLE", "message": "property service is unavailable"})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error_code": "PROPERTY_REQUEST_FAILED", "message": "property request failed"})
	}
}
