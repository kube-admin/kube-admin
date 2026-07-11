//go:build !embed

package web

import "embed"

// Assets 为空：普通构建（双镜像部署 / 本地 dev）不内嵌前端，后端仅作 API 服务。
// all-in-one 镜像用 -tags embed 编译 embed.go，以真实内嵌覆盖此定义。
var Assets embed.FS
