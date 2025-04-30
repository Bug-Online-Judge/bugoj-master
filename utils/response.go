package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`           // 业务状态码
	Message   string      `json:"message"`        // 描述信息
	Data      interface{} `json:"data,omitempty"` // 返回数据
	Timestamp int64       `json:"timestamp"`      // 时间戳
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "ok",
		Data:      data,
		Timestamp: time.Now().Unix(),
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:      code,
		Message:   msg,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	})
}
