package model

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	UserId       uint
	Role		 string
	TokenVersion uint
	jwt.RegisteredClaims
}