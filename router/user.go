package router

import (
	"bugoj-master/controller"
	"bugoj-master/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/api/register", controller.Register)
	r.POST("/api/login", controller.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", controller.Me)
	}
}
