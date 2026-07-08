# 更新日志

本项目遵循 [Keep a Changelog](https://keepachangelog.com/) 风格。

## [Unreleased]

## [0.1.0] - 2026-07

### Added
- 多集群管理：kubeconfig 内容 / 文件路径 / ServerURL+Token 三种接入方式，凭据 AES-256-GCM 加密存储
- 资源管理：Namespace / Node / Pod / Deployment / Service / ConfigMap / Secret 全功能 CRUD
- 通用资源浏览器：基于 dynamic client 管理任意 GVR（StatefulSet / DaemonSet / ReplicaSet / Ingress / NetworkPolicy / PV / PVC / StorageClass / HPA / ServiceAccount / Role / RoleBinding 等）
- 实时监控：接入 metrics-server，节点与 Pod CPU/内存实时使用率；Dashboard echarts 可视化（仪表盘 + 饼图）
- Event 事件管理
- Web 终端：Pod 交互式终端（WebSocket + xterm）
- 实时日志：WebSocket 流式日志，支持 follow / previous / 容器切换 / 搜索 / 暂停 / 下载
- YAML 工作台：Monaco 编辑器在线编辑与 apply 任意资源
- 安全：JWT 认证、bcrypt 密码、RBAC 角色鉴权（admin/operator/viewer）、操作审计日志
- 健康检查端点 `/healthz` `/readyz`、HTTP 服务优雅关闭
- 数据库文件持久化（替代内存模式）、配置统一（JWT/加密/TLS）
- 工程化：后端单元测试、前端 vitest 测试、golangci-lint、GitHub Actions CI、Makefile、Docker 部署
- 前端 i18n 框架与 `v-permission` 权限指令

### Security
- 集群凭据不再明文存储，统一加密；TLS 校验可配置（默认开启）
- 用户/集群管理接口限制仅 admin 角色访问
