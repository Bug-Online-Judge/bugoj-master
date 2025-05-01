package controller

import (
	"bugoj-master/config"
	"bugoj-master/model"
	"bugoj-master/utils"
	"net/http"
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

	var userRole model.Role
	if err := config.DB.FirstOrCreate(&userRole, model.Role{Name: "user"}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ensure default role"})
		return
	}

	user := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPwd,
		Roles:     []model.Role{userRole},
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

func UserProfile(c *gin.Context) {
	username := c.Query("username")
	email := c.Query("email")

	var user model.User
	var err error

	switch {
	case username != "":
		err = config.DB.Where("username = ?", username).First(&user).Error
	case email != "":
		err = config.DB.Where("email = ?", email).First(&user).Error
	default:
		utils.Fail(c, 1001, "username or email required")
		return
	}

	if err != nil {
		utils.Fail(c, 1002, "user not found")
		return
	}

	utils.Success(c, model.ToUserDTO(user))
}
