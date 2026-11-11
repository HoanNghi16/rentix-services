package handler

import (
	"net/http"
	"strings"

	"github.com/HoanNghi16/rentix-services/api-gateway/internal/authclient"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct{ auth *authclient.Client }

func RegisterRoutes(router gin.IRouter, auth *authclient.Client) {
	handler := &authHandler{auth: auth}
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/auth/register", handler.register)
	router.POST("/auth/login", handler.login)
}

type credentialsRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (handler *authHandler) register(c *gin.Context) {
	var input credentialsRequest
	if err := bindCredentials(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "request must contain a valid email and password"})
		return
	}
	_, err := handler.auth.Register(c.Request.Context(), &authclient.RegisterRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

func (handler *authHandler) login(c *gin.Context) {
	var input credentialsRequest
	if err := bindCredentials(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "request must contain a valid email and password"})
		return
	}
	response, err := handler.auth.Login(c.Request.Context(), &authclient.LoginRequest{Email: input.Email, Password: input.Password})
	if err != nil {
		writeGRPCError(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func bindCredentials(c *gin.Context, input *credentialsRequest) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
	if err := c.ShouldBindJSON(input); err != nil {
		return err
	}
	input.Email = strings.TrimSpace(input.Email)
	return nil
}

func writeGRPCError(c *gin.Context, err error) {
	switch status.Code(err) {
	case codes.InvalidArgument:
		c.JSON(http.StatusBadRequest, gin.H{"error_code": "INVALID_INPUT", "message": "email and password are required"})
	case codes.Unauthenticated:
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "INVALID_CREDENTIALS", "message": "invalid email or password"})
	case codes.AlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error_code": "EMAIL_ALREADY_USED", "message": "email is already registered"})
	case codes.Unavailable, codes.DeadlineExceeded:
		c.JSON(http.StatusServiceUnavailable, gin.H{"error_code": "AUTH_UNAVAILABLE", "message": "authentication service is unavailable"})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"error_code": "AUTH_REQUEST_FAILED", "message": "authentication request failed"})
	}
}
