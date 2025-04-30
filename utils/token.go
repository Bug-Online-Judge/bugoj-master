package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokens(id uint, username, role string) (accessToken, refreshToken string) {
	// Access token: 8小时
	atExp := time.Now().Add(8 * time.Hour)
	// Refresh token: 2天
	rtExp := time.Now().Add(2 * 24 * time.Hour)

	accessClaims := &Claims{
		UserID:   id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(atExp),
		},
	}
	refreshClaims := &Claims{
		UserID:   id,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rtExp),
		},
	}

	accessToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtKey)
	refreshToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtKey)
	return
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
