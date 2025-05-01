package controller

import (
	"bugoj-master/config"
	"bugoj-master/model"
	"bugoj-master/utils"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 4001, "invalid input")
		return
	}

	var user model.User
	if err := config.DB.
		Where("username = ? OR email = ?", req.Username, req.Username).
		First(&user).Error; err != nil {
		utils.Fail(c, 4002, "user not found")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Fail(c, 4003, "invalid password")
		return
	}

	accessToken, refreshToken, refreshKey := utils.GenerateTokens(user.ID, user.Username, user.Role)

	c.SetCookie("refresh_token", refreshToken, 3600*48, "/", "localhost", false, true)
	c.SetCookie("refresh_key", refreshKey, 3600*48, "/", "localhost", false, true)

	utils.Success(c, gin.H{
		"access_token": accessToken,
	})
}
