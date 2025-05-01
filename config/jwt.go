package config

import (
	"fmt"
	"os"

	"bugoj-master/global"
	"github.com/joho/godotenv"
)

func InitJWT() {
	// 加载 .env（如果 main 已经加载一次可以省略）
	_ = godotenv.Load(".env")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set in .env")
	}

	global.JWTKey = []byte(secret)
	fmt.Println("JWT 密钥加载成功")
}
