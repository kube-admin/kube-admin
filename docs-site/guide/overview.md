# 项目介绍

Kube Admin 是一个开源的多集群 Kubernetes 管理平台，提供从资源管理、实时监控到 Web 终端、YAML 工作台的完整能力，适合作为团队内部的 K8s 运维入口。

## 能力一览

| 模块 | 能力 |
|------|------|
| 多集群 | 接入多个集群并在界面切换；集群凭据加密存储 |
| 仪表盘 | 资源计数、集群 CPU/内存使用率、Pod 状态分布可视化 |
| 核心资源 | Namespace / Node / Pod / Deployment / Service / ConfigMap / Secret 全功能管理 |
| 资源浏览器 | 基于 dynamic client 管理任意 K8s 资源（StatefulSet / Ingress / PVC / HPA / SA / Role …） |
| 监控 | 节点与 Pod 实时 CPU/内存使用率（metrics-server） |
| 终端 | Pod 交互式终端（WebSocket + xterm） |
| 日志 | WebSocket 实时日志流，支持 follow / 上个容器 / 搜索 / 暂停 / 下载 |
| YAML 工作台 | Monaco 编辑器在线编辑与 apply 任意资源 |
| 事件 | 按命名空间/资源查询 Warning 与 Normal 事件 |
| 安全 | JWT 认证、bcrypt、RBAC 角色、操作审计日志 |
| 运维 | 健康探针、优雅关闭、Docker / k8s 部署 |

## 技术栈

- **前端**：Vue 3 + TypeScript + Element Plus + Pinia + Vue Router + echarts + Monaco + xterm
- **后端**：Go 1.24 + Gin + client-go + GORM(SQLite) + gorilla/websocket
- **运维**：Docker / docker-compose / Kubernetes

## 架构

```
浏览器 SPA ──HTTP/WebSocket──▶ Go Backend (Gin) ──client-go──▶ Cluster A / B / N
```

后端分层：`router → middleware(认证/RBAC/审计) → api → service → pkg/k8s`。多集群通过 `ClusterMiddleware` 按请求注入对应客户端；通用资源走 `dynamic client` + RESTMapper。

详见项目根 [docs/ARCHITECTURE.md](https://github.com/)（仓库内）。

## 角色与权限

| 角色 | 能力 |
|------|------|
| `admin` | 全部操作，含用户/集群/审计管理 |
| `operator` / `user` | 资源读写（含终端/日志/YAML apply），不能管理用户与集群 |
| `viewer`（预留） | 只读 |

## 下一步

- [快速上手](./getting-started.md)：本地 5 分钟跑起来
- [多集群管理](./clusters.md)：接入你的集群
