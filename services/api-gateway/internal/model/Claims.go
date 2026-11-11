package model

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	UserID       uint64 `json:"UserId"`
	Role         string `json:"Role"`
	TokenVersion uint64 `json:"TokenVersion"`
	jwt.RegisteredClaims
}
