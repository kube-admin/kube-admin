package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/pkg/logger"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录日志
		duration := time.Since(startTime)
		logger.Info("%s %s %d %v %s",
			c.Request.Method,
			c.Request.RequestURI,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}
