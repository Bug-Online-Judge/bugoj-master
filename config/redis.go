package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"bugoj-master/global"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if addr == "" {
		addr = "localhost:6379"
	}
	if dbStr == "" {
		dbStr = "0"
	}
	db, _ := strconv.Atoi(dbStr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis 连接失败: %v", err))
	}

	global.Redis = rdb
	fmt.Println("Redis 初始化成功")
}
