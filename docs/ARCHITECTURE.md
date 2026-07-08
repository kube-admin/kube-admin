# 架构说明

## 总览

Kube Admin 采用前后端分离架构，后端通过 `client-go` 与多个 Kubernetes 集群通信，前端为 SPA。

```
┌──────────────┐      HTTPS / WebSocket       ┌──────────────────────────┐
│   浏览器 SPA  │  ──────────────────────────▶ │        Go Backend        │
│  (Vue 3)     │                              │          (Gin)           │
└──────────────┘                              └───────────┬──────────────┘
                                                          │ client-go / dynamic
                                       ┌──────────────────┼──────────────────┐
                                       ▼                  ▼                  ▼
                                 ┌──────────┐       ┌──────────┐       ┌──────────┐
                                 │ Cluster A│       │ Cluster B│       │ Cluster N│
                                 └──────────┘       └──────────┘       └──────────┘
```

## 后端分层

| 层 | 职责 | 关键文件 |
|----|------|----------|
| `router` | 路由注册、中间件编排 | `internal/router/router.go` |
| `middleware` | 认证(JWT) / RBAC / 审计 / CORS / 集群选择 | `internal/middleware/*` |
| `api` | HTTP / WebSocket handler | `internal/api/*` |
| `service` | 业务逻辑，K8s 资源操作 | `internal/service/*` |
| `pkg/k8s` | K8s 客户端封装、多集群管理器 | `pkg/k8s/{client,manager}.go` |
| `pkg/crypto` | 凭据加解密 | `pkg/crypto/crypto.go` |

## 多集群机制

1. 启动时创建默认客户端（本地 kubeconfig）与多集群 `Manager`。
2. 每个请求通过 `ClusterMiddleware` 解析 `cluster_id`，从 DB 取集群配置（凭据解密），通过 `Manager` 获取/缓存对应的 dynamic client，注入到请求上下文。
3. 未指定 `cluster_id` 时使用默认集群客户端。

## 通用资源管理

`ResourceService` 基于 `dynamic.Interface` + `restmapper`，将 GVK 映射为 GVR，实现任意资源的 List / Get / Delete / Apply / Patch，避免为每种资源重复实现 typed service。前端 `Resources.vue` 通过 GVR 切换覆盖所有资源类型。

## 安全模型

- 用户密码 bcrypt 存储；JWT（HS256）鉴权，密钥由 `JWT_SECRET` 注入。
- 集群 `Token` / `ConfigContent` 写入数据库前 AES-256-GCM 加密，读取时解密。
- `AuthMiddleware` 校验 token；`RequireRole` / `WriteAuth` 实现接口级 RBAC；`AuditMiddleware` 记录所有写操作。
- 前端 `v-permission` 指令按角色控制元素显隐。

## 实时能力

- **终端 / 日志**：通过 `gorilla/websocket` 升级连接，后端桥接 K8s `remotecommand`（exec）与日志流到前端。WebSocket token 通过 query 参数传递（受 `AuthMiddleware` 校验）。
- **监控**：调用 metrics-server 的 `MetricsV1beta1` 接口聚合节点/容器实时使用率，优雅降级（metrics-server 缺失时使用率为空）。

## 数据持久化

使用 SQLite（GORM）存储用户、集群、审计日志，文件路径由 `DB_PATH` 配置。如需多副本生产部署，建议替换为 MySQL/PostgreSQL（仅需更换 GORM 驱动）。
