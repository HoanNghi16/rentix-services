package handler

import (
	"net/http"

	"github.com/HoanNghi16/rentix-services/api-gateway/internal/employeeclient"
	"github.com/HoanNghi16/rentix-services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type employeeHandler struct{ employee *employeeclient.Client }

func RegisterEmployeeRoutes(router gin.IRouter, employee *employeeclient.Client, accessSecret string) {
	routes := router.Group("/employee")
	routes.Use(middleware.AuthMiddleware(accessSecret))
	routes.GET("", (&employeeHandler{employee: employee}).getEmployees)
}

func (handler *employeeHandler) getEmployees(c *gin.Context) {
	ownerID, ok := middleware.AuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	role, ok := middleware.AuthenticatedRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
		return
	}
	response, err := handler.employee.GetEmployees(c.Request.Context(), &employeeclient.GetEmployeesRequest{
		OwnerID: ownerID, Role: role,
	})
	if err != nil {
		writeEmployeeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func writeEmployeeGRPCError(c *gin.Context, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "request is invalid"})
	case codes.PermissionDenied:
		c.JSON(http.StatusForbidden, gin.H{"error_code": "FORBIDDEN", "message": "only owners can access employees"})
	case codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
	case codes.Unavailable, codes.DeadlineExceeded:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error_code": "EMPLOYEE_UNAVAILABLE", "message": "employee service is unavailable"})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error_code": "EMPLOYEE_REQUEST_FAILED", "message": "employee request failed"})
	}
}
