package main

import (
	"bugoj-master/config"
	"bugoj-master/middleware"
	"bugoj-master/router"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.SetupRoutes(r)

	r.Run(":8080") // 启动服务
}
