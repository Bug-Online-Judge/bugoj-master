package controller

import (
	"bugoj-master/config"
	"bugoj-master/model"
	"bugoj-master/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 4001, "invalid input")
		return
	}

	hashedPwd, _ := utils.HashPassword(req.Password)

	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPwd,
		Role:      "user",
		CreatedAt: time.Now(),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		utils.Fail(c, 4002, "username or email exists")
		return
	}
	utils.Success(c, user)
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 4001, "invalid input")
		return
	}

	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, 4002, "user not found")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Fail(c, 4003, "invalid password")
		return
	}

	access, refresh := utils.GenerateTokens(user.ID, user.Username, user.Role)
	utils.Success(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		utils.Fail(c, 4001, "refresh token required")
		return
	}

	// 校验 refresh token 合法性和有效性
	claims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		utils.Fail(c, 401, "invalid or expired refresh token")
		return
	}

	// 生成新 access + refresh token（可选：refresh 可复用）
	access, refresh := utils.GenerateTokens(claims.UserID, claims.Username, claims.Role)

	utils.Success(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh, // 若想刷新 refresh，则返回新 refresh，否则可省略
	})
}

func Me(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}
