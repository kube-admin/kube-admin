package web

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegisterSPA 注册 SPA 静态资源与 index.html fallback。
//
// 仅当内嵌了前端 dist（-tags embed，all-in-one 镜像）时生效：
// Gin 已注册的 /api、/ws 路由优先匹配，NoRoute 仅兜底未命中的前端路由与静态资源。
// 普通 API-only 构建（Assets 为空）下 no-op，不影响双镜像部署形态。
func RegisterSPA(r *gin.Engine) {
	sub, err := fs.Sub(Assets, "dist")
	if err != nil {
		return
	}
	// dist 无 index.html 视为未内嵌前端（如本地误带 -tags embed 但未拷前端产物）
	if f, err := sub.Open("index.html"); err != nil {
		return
	} else {
		_ = f.Close()
	}
	r.NoRoute(func(c *gin.Context) {
		serveSPA(c, sub)
	})
}

// serveSPA 处理未命中路由的请求：静态文件直接返回，其余 fallback 到 index.html。
func serveSPA(c *gin.Context, fsys fs.FS) {
	p := c.Request.URL.Path
	// API / WebSocket 未匹配路由：返回 JSON 404，不 fallback 到 index.html
	if strings.HasPrefix(p, "/api/") || strings.HasPrefix(p, "/ws/") {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}
	name := strings.TrimPrefix(p, "/")
	if name == "" {
		name = "index.html"
	}
	if data, err := fs.ReadFile(fsys, name); err == nil {
		c.Data(http.StatusOK, contentType(name), data)
		return
	}
	// SPA fallback：未命中的前端路由交由前端 router 处理
	if idx, err := fs.ReadFile(fsys, "index.html"); err == nil {
		c.Data(http.StatusOK, "text/html; charset=utf-8", idx)
		return
	}
	c.Status(http.StatusNotFound)
}

// contentType 按扩展名返回静态资源 Content-Type，避免依赖 mime 包的模糊匹配与系统差异。
func contentType(name string) string {
	switch path.Ext(name) {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".js", ".mjs":
		return "application/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".map":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
