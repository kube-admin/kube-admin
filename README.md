# Kube Admin

> 基于 Vue 3 + Go + Kubernetes 的开源多集群 Kubernetes 管理平台，提供资源管理、实时监控、Web 终端、YAML 编辑与 RBAC 鉴权等能力。

## ✨ 特性

| 能力 | 说明 |
|------|------|
| **多集群管理** | 通过 kubeconfig / 文件路径 / ServerURL+Token 接入多集群，凭据 AES-256-GCM 加密存储 |
| **资源管理** | Namespace / Node / Pod / Deployment / Service / ConfigMap / Secret，以及通过通用资源浏览器管理 StatefulSet / DaemonSet / Ingress / PVC / PV / StorageClass / HPA / SA / Role 等任意资源 |
| **实时监控** | 接入 metrics-server，节点与 Pod 的 CPU/内存实时使用率、Dashboard 可视化（echarts） |
| **Web 终端** | Pod 交互式终端（WebSocket + xterm），实时日志流（follow / previous / 搜索 / 下载） |
| **YAML 工作台** | Monaco 编辑器在线编辑、校验、应用任意资源（create or update） |
| **安全** | JWT 认证、bcrypt 密码、RBAC 角色鉴权（admin / operator / viewer）、操作审计日志 |
| **工程化** | 前端 i18n、单元测试、GitHub Actions CI、Docker 一键部署 |

## 🏗 技术栈

- **前端**：Vue 3 + TypeScript + Element Plus + Pinia + Vue Router + echarts + Monaco + xterm
- **后端**：Go 1.24 + Gin + client-go + GORM(SQLite) + gorilla/websocket
- **运维**：Docker / docker-compose / Kubernetes

## 🚀 快速开始

### 前置条件
- Go 1.24+
- Node.js 18+
- 可访问的 Kubernetes 集群（或 minikube / kind），本地 `~/.kube/config` 已配置

### 方式一：本地开发

```bash
# 后端
cd backend
go mod tidy
go run cmd/main.go          # 默认 :8080

# 前端（另开终端）
cd frontend
npm install
npm run dev                 # 默认 :3000
```

访问 http://localhost:3000 ，使用默认账号 `admin / admin123` 登录。

### 方式二：Make

```bash
make backend     # 启动后端
make frontend    # 启动前端
make test        # 运行前后端测试
```

### 方式三：Docker

```bash
./start-docker.sh   # 或 docker-compose up --build
```

## ⚙️ 配置

后端通过环境变量配置（见 `.env.example`）：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `8080` | 后端端口 |
| `KUBECONFIG` | `~/.kube/config` | 默认集群 kubeconfig 路径 |
| `JWT_SECRET` | （开发默认） | JWT 签名密钥，**生产必须修改** |
| `ENCRYPT_KEY` | （开发默认） | 集群凭据加密密钥，**生产必须修改** |
| `DB_PATH` | `data/kubeadm.db` | SQLite 数据库路径 |
| `TLS_SKIP_VERIFY` | `false` | 是否跳过集群 TLS 校验（仅开发） |
| `GIN_MODE` | `debug` | gin 运行模式 |

## 📡 API 概览

所有 K8s 操作需 `Authorization: Bearer <token>`，写操作需 `admin`/`operator` 角色。

```
POST   /api/v1/auth/login              登录
GET    /api/v1/auth/user               当前用户

# 集群与用户（仅 admin）
GET/POST/PUT/DELETE /api/v1/clusters
GET/POST/PUT/DELETE /api/v1/users
GET    /api/v1/audit/logs              审计日志

# K8s 资源（?cluster_id=&namespace=）
GET    /api/v1/dashboard/stats         集群统计 + 使用率
GET    /api/v1/nodes | /pods | /deployments | /services | ...
GET    /api/v1/events                  事件
GET    /api/v1/pods/:name/logs/stream  WebSocket 实时日志
GET    /api/v1/pods/:name/terminal     WebSocket 终端

# 通用资源（任意 GVR）
GET/DELETE /api/v1/resources[/:name]
POST   /api/v1/resources/apply         YAML apply
PATCH  /api/v1/resources/:name

GET    /healthz | /readyz              健康检查
```

## 🧪 测试

```bash
cd backend && go test ./...
cd frontend && npm run test
```

## 📦 Kubernetes 部署

后端镜像部署示例（已含健康探针）：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-admin-backend
spec:
  replicas: 2
  selector:
    matchLabels: { app: kube-admin-backend }
  template:
    metadata:
      labels: { app: kube-admin-backend }
    spec:
      serviceAccountName: kube-admin-sa
      containers:
      - name: backend
        image: your-registry/kube-admin-backend:latest
        ports: [{ containerPort: 8080 }]
        env:
        - { name: JWT_SECRET, valueFrom: { secretKeyRef: { name: kube-admin-secret, key: jwt-secret } } }
        - { name: ENCRYPT_KEY, valueFrom: { secretKeyRef: { name: kube-admin-secret, key: encrypt-key } } }
        livenessProbe:  { httpGet: { path: /healthz, port: 8080 }, initialDelaySeconds: 10, periodSeconds: 20 }
        readinessProbe: { httpGet: { path: /readyz,   port: 8080 }, initialDelaySeconds: 5,  periodSeconds: 10 }
```

> 集群内部署时，后端会自动使用 ServiceAccount 凭据访问当前集群；ServiceAccount 需按最小权限原则授权。

## 🗺 路线图

- [ ] Informer / Watch 实时资源变更推送
- [ ] Helm 应用市场
- [ ] Prometheus 历史监控集成
- [ ] 多租户与细粒度权限策略
- [ ] 前端全量国际化文案迁移

## 📁 项目结构

```
kube-admin/
├── backend/                 # Go 后端
│   ├── cmd/main.go          # 入口（优雅关闭）
│   ├── config/              # 配置
│   ├── database/            # GORM/SQLite
│   ├── internal/
│   │   ├── api/             # HTTP handler
│   │   ├── middleware/      # 认证 / RBAC / 审计 / CORS
│   │   ├── model/           # 数据模型
│   │   ├── router/          # 路由
│   │   └── service/         # 业务逻辑 + 通用资源 service
│   └── pkg/{k8s,crypto,logger}
├── frontend/                # Vue 3 前端
│   └── src/{apis,views,components,stores,locales,directives}
├── .github/workflows/ci.yml
└── Makefile
```

详见 [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md)。

## 🤝 贡献

欢迎 Issue 与 PR，请阅读 [CONTRIBUTING.md](./CONTRIBUTING.md)。

## 📄 License

[MIT](./LICENSE)
