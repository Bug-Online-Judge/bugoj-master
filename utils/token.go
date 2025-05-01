package utils

import (
	"bugoj-master/global"
	"bugoj-master/model"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ctx = context.Background()

type Claims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Role     []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateTokens 返回 accessToken、refreshToken、refreshKey
func GenerateTokens(id uint, username string, roles interface{}) (accessToken, refreshToken, refreshKey string) {
	// Access token: 8小时
	atExp := time.Now().Add(8 * time.Hour)
	// Refresh token: 2天
	rtExp := time.Now().Add(48 * time.Hour)

	// Convert []Role to []string
	var roleNames []string
	switch v := roles.(type) {
	case []model.Role:
		roleNames = make([]string, len(v))
		for i, r := range v {
			roleNames[i] = r.Name
		}
	case []string:
		roleNames = v
	default:
		roleNames = []string{}
	}

	accessClaims := &Claims{
		UserID:   id,
		Username: username,
		Role:     roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(atExp),
		},
	}
	refreshClaims := &Claims{
		UserID:   id,
		Username: username,
		Role:     roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rtExp),
		},
	}

	accessToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(global.JWTKey)
	refreshToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(global.JWTKey)

	// 构造 refresh key 并写入 Redis
	refreshKey = fmt.Sprintf("refresh:%d:%s", id, uuid.New().String())
	SaveRefreshTokenToRedis(refreshKey, 48*time.Hour)

	return accessToken, refreshToken, refreshKey
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return global.JWTKey, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func SaveRefreshTokenToRedis(key string, ttl time.Duration) {
	global.Redis.Set(ctx, key, 1, ttl)
}

func DeleteRefreshTokenFromRedis(key string) {
	global.Redis.Del(ctx, key)
}

func RefreshTokenExists(key string) bool {
	return global.Redis.Exists(ctx, key).Val() == 1
}
