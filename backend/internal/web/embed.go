//go:build embed

package web

import "embed"

// Assets 内嵌前端静态资源。
// dist 由 all-in-one Dockerfile 在构建时从前端构建产物拷入 backend/internal/web/dist。
// 仅 -tags embed 构建生效（all-in-one 单镜像）；普通构建编译 noembed.go，Assets 为空。
//
//go:embed all:dist
var Assets embed.FS
