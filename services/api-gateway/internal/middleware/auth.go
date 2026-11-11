package middleware

import (
	"net/http"
	"strings"

	"github.com/HoanNghi16/rentix-services/api-gateway/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthenticatedUserIDKey = "authenticated_user_id"
	AuthenticatedRoleKey   = "authenticated_role"
)

func AuthMiddleware(accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}
		if strings.TrimSpace(accessSecret) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error_code": "AUTH_CONFIG_ERROR", "message": "authentication is not configured"})
			return
		}

		claims := new(model.TokenClaims)
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(accessSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || token == nil || !token.Valid || claims.ExpiresAt == nil || claims.UserID == 0 || strings.TrimSpace(claims.Role) == "" {
			unauthorized(c)
			return
		}

		c.Set(AuthenticatedUserIDKey, claims.UserID)
		c.Set(AuthenticatedRoleKey, claims.Role)
		c.Next()
	}
}

func AuthenticatedUserID(c *gin.Context) (uint64, bool) {
	userID, ok := c.Get(AuthenticatedUserIDKey)
	if !ok {
		return 0, false
	}
	value, ok := userID.(uint64)
	return value, ok && value != 0
}

func AuthenticatedRole(c *gin.Context) (string, bool) {
	role, ok := c.Get(AuthenticatedRoleKey)
	if !ok {
		return "", false
	}
	value, ok := role.(string)
	return value, ok && strings.TrimSpace(value) != ""
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHENTICATED", "message": "valid bearer token required"})
}
