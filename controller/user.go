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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
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
		c.JSON(http.StatusConflict, gin.H{"error": "username or email exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registered successfully"})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}
	token := utils.GenerateToken(user.ID, user.Username, user.Role)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Me(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, user)
}
