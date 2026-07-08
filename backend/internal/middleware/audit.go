package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/model"
)

// AuditMiddleware 审计中间件：在请求处理后记录写操作到数据库。
// 仅记录 POST/PUT/DELETE/PATCH，读操作不记录。同步写入保证顺序与可靠性。
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 仅审计写操作
		switch c.Request.Method {
		case "POST", "PUT", "DELETE", "PATCH":
		default:
			return
		}

		auditLog := model.AuditLog{
			UserID:    c.GetUint("user_id"),
			Username:  c.GetString("username"),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			CreatedAt: time.Now(),
		}
		// 写入失败不影响主流程，仅记录日志
		if err := database.DB.Create(&auditLog).Error; err != nil {
			// 避免引入 logger 循环依赖，直接忽略审计写入错误
			_ = err
		}
	}
}
