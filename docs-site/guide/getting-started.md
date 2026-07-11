# 快速上手

本地 5 分钟把 Kube Admin 跑起来。

## 环境要求

| 工具 | 版本 | 说明 |
|---|---|---|
| Go | `>= 1.24` | 后端 |
| Node.js | `>= 18` | 前端 |
| Kubernetes | 任意版本 | 或 minikube / kind，本地 `~/.kube/config` 已配置 |

## 安装与启动

### 方式一：分别启动（推荐开发）

```bash
# 后端
cd backend
go mod tidy
go run cmd/main.go            # 默认 :8080

# 前端（另开终端）
cd frontend
npm install
npm run dev                   # 默认 :3000
```

浏览器打开 http://localhost:3000 。

### 方式二：Make

```bash
make backend      # 启动后端
make frontend     # 启动前端
make test         # 运行前后端测试
```

### 方式三：Docker 一键

```bash
make docker-up   # 经 deploy/deploy.sh 构建并启动（含密钥校验 + 健康等待）
```

## 默认账号

| 用户名 | 密码 | 角色 |
|---|---|---|
| `admin` | `admin123` | admin |

::: warning 生产环境务必修改
首次登录后请立即在「系统设置 → 用户管理」修改密码，并通过环境变量设置强 `JWT_SECRET` 与 `ENCRYPT_KEY`。
:::

## 配置项

后端通过环境变量配置（参考 `.env.example`）：

| 变量 | 默认值 | 说明 |
|---|---|---|
| `PORT` | `8080` | 后端端口 |
| `KUBECONFIG` | `~/.kube/config` | 默认集群 kubeconfig 路径 |
| `JWT_SECRET` | 开发默认 | JWT 签名密钥，**生产必须修改** |
| `ENCRYPT_KEY` | 开发默认 | 集群凭据加密密钥，**生产必须修改** |
| `DB_PATH` | `data/kubeadm.db` | SQLite 数据库路径 |
| `TLS_SKIP_VERIFY` | `false` | 是否跳过集群 TLS 校验（仅开发） |
| `GIN_MODE` | `debug` | gin 运行模式 |

## 常用命令

| 命令 | 作用 |
|---|---|
| `make backend` / `make frontend` | 启动后端 / 前端 |
| `make build` | 构建前后端 |
| `make test` | 全部测试 |
| `cd backend && go test ./...` | 后端单元测试 |
| `cd frontend && npm run test` | 前端 vitest |
| `cd frontend && npm run build` | 前端生产构建 |
| `cd docs-site && npm run docs:dev` | 文档站本地预览 |

## 下一步

- [多集群管理](./clusters.md)：接入业务集群
- [核心资源](./workloads.md)：管理 Pod / Deployment 等
- [部署指南](./deploy.md)：上线生产
