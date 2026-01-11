package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否是WebSocket升级请求
		if c.GetHeader("Upgrade") == "websocket" {
			// 对于WebSocket请求，尝试从查询参数获取token
			tokenString := c.Query("token")
			if tokenString == "" {
				// 如果查询参数中没有token，尝试从Cookie获取
				tokenString, _ = c.Cookie("token")
			}

			if tokenString == "" {
				// 如果还是没有token，检查Authorization头
				authHeader := c.GetHeader("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenString = authHeader[7:] // "Bearer "之后的部分
				}
			}

			if tokenString != "" {
				// 验证token
				claims, err := model.ParseToken(tokenString)
				if err == nil {
					// Token有效，将用户信息存入上下文
					c.Set("user_id", claims.UserID)
					c.Set("username", claims.Username)
					c.Set("role", claims.Role)
					c.Next()
					return
				}
			}

			// WebSocket连接如果没有有效的token，直接返回401
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(401, "未经授权的访问"))
			c.Abort()
			return
		}

		// 对于普通HTTP请求，保持原有的认证逻辑
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			// 尝试从Cookie获取token
			tokenString, _ = c.Cookie("token")
		}

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(401, "缺少访问令牌"))
			c.Abort()
			return
		}

		tokenString = tokenString[7:] // "Bearer "之后的部分

		claims, err := model.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse(401, "无效的访问令牌"))
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
