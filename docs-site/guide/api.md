# API 速查

所有接口前缀 `/api/v1`。除 `POST /auth/login` 与 `GET /healthz|/readyz` 外，均需 `Authorization: Bearer <token>`。

## 认证与用户

| 方法 | 路径 | 说明 |
|---|---|---|
| POST | `/auth/login` | 登录，返回 JWT |
| GET | `/auth/user` | 当前用户信息 |
| GET/POST/PUT/DELETE | `/users[/:id]` | 用户管理（admin） |
| GET | `/audit/logs` | 审计日志（admin） |

## 集群（admin）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET/POST/PUT/DELETE | `/clusters[/:id]` | 集群 CRUD |
| POST | `/clusters/test-connection` | 凭明文凭据测试（创建前预检） |
| POST | `/clusters/:id/test-connection` | 按已存集群 ID 测试 |

## K8s 资源

所有 K8s 资源接口支持查询参数 `cluster_id`、`namespace`。

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/dashboard/stats` | 集群统计 + 实时使用率 |
| GET | `/namespaces` / `/nodes` / `/pods` / `/deployments` / `/services` / `/configmaps` / `/secrets` | 列表 |
| GET | `/pods/:name`、`/deployments/:name` 等 | 详情 |
| DELETE | `/pods/:name`、`/deployments/:name` 等 | 删除 |
| PUT | `/deployments/:name/scale` | 扩缩容（`?replicas=N`） |
| PUT | `/deployments/:name/restart` | 滚动重启 |
| POST | `/pods/yaml`、`/deployments/yaml`、`/services/yaml` | 由 YAML 创建 |
| GET | `/pods/:name/logs` | 一次性日志（HTTP） |
| GET | `/pods/:name/logs/stream` | WebSocket 实时日志 |
| GET | `/pods/:name/terminal` | WebSocket 终端 |
| GET | `/events` | 事件（`?kind=&name=` 过滤） |

## 通用资源（任意 GVR）

查询参数：`group`、`version`、`resource`、`namespace`。

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/resources` | 列表 |
| GET | `/resources/:name` | 详情 |
| DELETE | `/resources/:name` | 删除 |
| POST | `/resources/apply` | YAML 创建或更新（`{ "yaml": "..." }`） |
| PATCH | `/resources/:name` | 补丁（`{ "patch_type": "strategic", "data": "..." }`） |

## 健康检查

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/healthz` | 存活 |
| GET | `/readyz` | 就绪 |

## WebSocket 鉴权

WebSocket 连接无法使用 `Authorization` 头时，可通过 `?token=<jwt>` 传递，后端在升级握手时校验。

## 响应格式

```json
{ "code": 0, "message": "success", "data": { ... } }
```

错误时 `code` 非 0，`data` 缺省，`message` 为错误描述。
