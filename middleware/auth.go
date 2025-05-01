package middleware

import (
	"bugoj-master/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || len(tokenStr) <= 7 || tokenStr[:7] != "Bearer " {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "missing or malformed token"},
			)
			return
		}

		token := tokenStr[7:]

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid or expired access token"},
			)
			return
		}

		// access token 剩余时间少于 30 分钟，自动刷新
		if time.Until(claims.ExpiresAt.Time) < 30*time.Minute {
			refreshToken, err1 := c.Cookie("refresh_token")
			refreshKey, err2 := c.Cookie("refresh_key")

			if err1 == nil && err2 == nil {
				rtClaims, err := utils.ParseToken(refreshToken)
				if err == nil && utils.RefreshTokenExists(refreshKey) {
					utils.DeleteRefreshTokenFromRedis(refreshKey)

					// 生成新的 token 和 refreshKey
					newAT, newRT, newRK := utils.GenerateTokens(
						rtClaims.UserID,
						rtClaims.Username,
						rtClaims.Role,
					)

					// 写入新的 cookie
					c.SetCookie(
						"refresh_token",
						newRT, 3600*48,
						"/",
						"localhost",
						false,
						true,
					)
					c.SetCookie(
						"refresh_key",
						newRK,
						3600*48,
						"/",
						"localhost",
						false, true,
					)

					c.Header("X-New-Access-Token", newAT)

					claims = rtClaims
				}
			}
		}

		c.Set("user", claims)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
