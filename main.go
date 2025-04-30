package main

import (
	"bugoj-master/config"
	"bugoj-master/middleware"
	"bugoj-master/router"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	config.InitRedis()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	router.SetupRoutes(r)

	err := r.Run(":8080")
	if err != nil {
		return
	} // Start up the server
}
