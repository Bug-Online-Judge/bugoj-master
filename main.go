package main

import (
	"bugoj-master/config"
	"bugoj-master/middleware"
	"bugoj-master/router"
	"github.com/gin-gonic/gin"

	_ "bugoj-master/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title BugOJ API
// @version 1.0
// @description API documentation for BugOJ backend
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	config.InitDB()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.SetupRoutes(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":8080")
	if err != nil {
		return
	} // Start up the server
}
