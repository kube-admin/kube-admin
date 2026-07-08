package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
)

// RequireRole 要求当前用户具备指定角色之一，否则返回 403。
// 用于细粒度控制管理类接口（如用户管理、集群管理仅 admin）。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role := c.GetString("role")
		if _, ok := allowed[role]; ok {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, model.ErrorResponse(403, "权限不足，需要角色: "+joinRoles(roles)))
		c.Abort()
	}
}

// WriteAuth 写操作鉴权：读操作（GET/HEAD/OPTIONS）放行，
// 写操作要求 admin/operator/user 角色。viewer 等只读角色被拒绝。
func WriteAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}
		role := c.GetString("role")
		switch role {
		case "admin", "operator", "user":
			c.Next()
		default:
			c.JSON(http.StatusForbidden, model.ErrorResponse(403, "当前角色无写权限"))
			c.Abort()
		}
	}
}

// joinRoles 简单拼接角色名用于错误提示
func joinRoles(roles []string) string {
	out := ""
	for i, r := range roles {
		if i > 0 {
			out += "/"
		}
		out += r
	}
	return out
}
