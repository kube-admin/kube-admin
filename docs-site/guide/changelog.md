# 更新日志

遵循 [Keep a Changelog](https://keepachangelog.com/) 风格。

## [0.1.0] - 2026-07

### Added

- **多集群管理**：kubeconfig 内容 / 文件路径 / ServerURL+Token 三种接入方式，凭据 AES-256-GCM 加密存储
- **核心资源 CRUD**：Namespace / Node / Pod / Deployment / Service / ConfigMap / Secret
- **通用资源浏览器**：基于 dynamic client 管理任意 GVR（StatefulSet / DaemonSet / Ingress / NetworkPolicy / PV / PVC / StorageClass / HPA / ServiceAccount / Role / RoleBinding）
- **实时监控**：接入 metrics-server，节点与 Pod CPU/内存实时使用率；Dashboard echarts 可视化
- **事件管理**：按命名空间 / 资源 / 类型查询
- **Web 终端**：Pod 交互式终端（WebSocket + xterm）
- **实时日志**：WebSocket 流式（follow / previous / 容器切换 / 搜索 / 暂停 / 下载）
- **YAML 工作台**：Monaco 编辑器在线编辑与 apply
- **安全**：JWT + bcrypt + RBAC 角色鉴权 + 操作审计日志
- **运维**：`/healthz` `/readyz` 健康探针、HTTP 优雅关闭
- **工程化**：后端单元测试 + golangci-lint、前端 vitest、GitHub Actions CI、Makefile、Docker 部署
- **前端**：i18n 框架、`v-permission` 权限指令

### Security

- 集群凭据统一加密存储，API 响应脱敏
- TLS 校验可配置（默认开启）
- 用户/集群/审计管理接口限制 admin 角色
- 默认 SQLite 持久化（替代内存模式），补全 `.gitignore` 防止数据库泄露
