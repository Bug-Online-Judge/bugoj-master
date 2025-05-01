package controller

import (
	"bugoj-master/config"
	"bugoj-master/model"
	"bugoj-master/utils"
	"time"

	"github.com/gin-gonic/gin"
)

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
	
	var existing model.User
	if err := config.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error; err == nil {
		utils.Fail(c, 4002, "username or email already exists")
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
		utils.Fail(c, 5000, "internal error")
		return
	}

	utils.Success(c, user)
}

func Me(c *gin.Context) {
	if user, ok := c.Get("user"); ok {
		utils.Success(c, user)
	} else {
		utils.Fail(c, 401, "unauthorized")
	}
}
