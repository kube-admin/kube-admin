# Kubernetes 管理系统

一个基于 Vue 3 + Go + Kubernetes 的全栈 K8s 集群管理系统。

## 项目结构

```
kube-admin/
├── backend/              # Go 后端服务
│   ├── cmd/             # 应用入口
│   ├── config/          # 配置文件
│   ├── internal/        # 内部代码
│   │   ├── api/        # API handlers
│   │   ├── middleware/ # 中间件 (JWT认证、CORS等)
│   │   ├── model/      # 数据模型
│   │   ├── router/     # 路由配置
│   │   └── service/    # 业务逻辑层
│   └── pkg/            # 可复用包
│       └── k8s/        # K8s客户端封装
└── frontend/           # Vue 3 前端
    ├── src/
    │   ├── apis/       # API 调用
    │   ├── views/      # 页面组件
    │   │   └── k8s/   # K8s 资源管理页面
    │   └── router/     # 路由配置
    └── vite.config.ts  # Vite 配置
```

## 技术栈

### 后端
- **Go 1.24+**: 编程语言
- **Gin**: Web 框架
- **client-go**: Kubernetes Go 客户端
- **JWT**: 身份认证
- **RESTful API**: 接口设计

### 前端
- **Vue 3**: 渐进式 JavaScript 框架
- **TypeScript**: 类型安全
- **Element Plus**: UI 组件库
- **Vue Router**: 路由管理
- **Pinia**: 状态管理
- **Axios**: HTTP 客户端
- **Vite**: 构建工具

## 功能特性

### 已实现功能

#### 后端功能
- ✅ **认证授权**: JWT Token 认证、用户登录
- ✅ **Dashboard API**: 集群资源统计、Pod 状态分布
- ✅ **Namespace**: 列表、创建、删除、资源统计
- ✅ **Node**: 列表、详情、资源状态、容量信息
- ✅ **Pod**: 列表、详情、删除、日志查看
- ✅ **Deployment**: 列表、详情、删除、扩缩容、重启
- ✅ **Service**: 列表、详情、删除
- ✅ **ConfigMap**: 列表、详情、创建、更新、删除
- ✅ **Secret**: 列表、详情、创建、更新、删除 (Base64 编解码)
- ✅ **日志系统**: HTTP 请求日志记录、错误日志
- ✅ **多命名空间支持**: 命名空间切换
- ✅ **实时状态**: 资源状态实时展示

#### 前端功能
- ✅ **Dashboard 页面**: 集群资源统计可视化、Pod 状态分布、快速操作
- ✅ **Namespace 管理**: 创建/删除、资源统计、详情查看
- ✅ **Node 管理**: 节点列表、详情查看、资源使用可视化、节点条件
- ✅ **Pod 管理**: 列表展示、删除、日志查看 (容器选择)
- ✅ **Deployment 管理**: 列表展示、扩缩容对话框、重启、删除
- ✅ **Service 管理**: 列表展示、端口配置、Selector 查看、删除
- ✅ **ConfigMap 管理**: 创建/编辑/删除、Key-Value 编辑器、详情查看
- ✅ **Secret 管理**: 创建/编辑/删除、明文/密文切换、Base64 自动处理
- ✅ **命名空间切换**: 所有资源页面支持命名空间快速切换
- ✅ **自动刷新**: 可配置的自动刷新功能 (5/10/30/60秒)
- ✅ **复用组件**: AutoRefresh、ResourceActions、StatusTag
- ✅ **工具函数库**: 40+ 个通用工具函数
- ✅ **响应式设计**: 适配各种屏幕尺寸
- ✅ **操作反馈**: Loading 状态、成功/失败提示、确认弹窗

### 可扩展功能
- Service 管理 (LoadBalancer、Ingress)
- StatefulSet / DaemonSet 管理
- PersistentVolume / PersistentVolumeClaim 管理
- 资源使用监控 (CPU/Memory Metrics)
- Helm Chart 部署
- YAML 文件编辑器
- 事件 (Events) 查看
- 日志聚合与搜索
- WebSocket 实时日志流

## 快速开始

### 前置要求
- Go 1.21+
- Node.js 16+
- Kubernetes 集群 (或 minikube/kind)
- kubectl 配置完成

### 后端启动

1. 进入后端目录
```bash
cd backend
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境变量 (可选)
```bash
export PORT=8080
export KUBECONFIG=/path/to/your/kubeconfig
export JWT_SECRET=your-secret-key
```

4. 启动服务
```bash
go run cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端启动

1. 进入前端目录
```bash
cd frontend
```

2. 安装依赖
```bash
npm install
# 或
pnpm install
```

3. 启动开发服务器
```bash
npm run dev
# 或
pnpm dev
```

前端服务将在 `http://localhost:3000` 启动

## API 接口文档

### 认证接口

#### 登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

Response:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1...",
    "expires_at": 1234567890,
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

#### 获取用户信息
```
GET /api/v1/auth/user
Authorization: Bearer <token>
```

### Dashboard 接口

```
GET /api/v1/dashboard/stats              # 获取集群统计信息
```

### Namespace 接口

```
GET /api/v1/namespaces              # 获取所有命名空间
POST /api/v1/namespaces             # 创建命名空间
DELETE /api/v1/namespaces/:name     # 删除命名空间
```

### Node 接口

```
GET /api/v1/nodes                   # 获取Node列表
GET /api/v1/nodes/:name             # 获取Node详情
```

### Pod 接口

```
GET /api/v1/pods?namespace=default           # 获取Pod列表
GET /api/v1/pods/:name?namespace=default     # 获取Pod详情
DELETE /api/v1/pods/:name?namespace=default  # 删除Pod
GET /api/v1/pods/:name/logs?namespace=default&container=xxx&tail_lines=100  # 获取Pod日志
```

### Deployment 接口

```
GET /api/v1/deployments?namespace=default                      # 获取Deployment列表
GET /api/v1/deployments/:name?namespace=default                # 获取Deployment详情
DELETE /api/v1/deployments/:name?namespace=default             # 删除Deployment
PUT /api/v1/deployments/:name/scale?namespace=default&replicas=3    # 扩缩容
PUT /api/v1/deployments/:name/restart?namespace=default        # 重启Deployment
```

### ConfigMap 接口

```
GET /api/v1/configmaps?namespace=default          # 获取ConfigMap列表
GET /api/v1/configmaps/:name?namespace=default    # 获取ConfigMap详情
POST /api/v1/configmaps                          # 创建ConfigMap
PUT /api/v1/configmaps/:name?namespace=default   # 更新ConfigMap
DELETE /api/v1/configmaps/:name?namespace=default # 删除ConfigMap
```

### Secret 接口

```
GET /api/v1/secrets?namespace=default              # 获取Secret列表
GET /api/v1/secrets/:name?namespace=default&decode=true  # 获取Secret详情
POST /api/v1/secrets                               # 创建Secret
PUT /api/v1/secrets/:name?namespace=default        # 更新Secret
DELETE /api/v1/secrets/:name?namespace=default     # 删除Secret
```

## 配置说明

### 后端配置

环境变量:
- `PORT`: 服务端口 (默认: 8080)
- `KUBECONFIG`: kubeconfig 文件路径 (默认: ~/.kube/config)
- `JWT_SECRET`: JWT 密钥 (必须修改)

### Kubeconfig 配置

系统支持两种方式连接 Kubernetes 集群:

1. **本地 kubeconfig**: 自动读取 `~/.kube/config` 或通过 `KUBECONFIG` 环境变量指定
2. **集群内配置**: 在 K8s Pod 中运行时自动使用 ServiceAccount

### 默认用户

- 用户名: `admin`
- 密码: `admin123`

⚠️ **生产环境请务必修改默认密码和 JWT 密钥!**

## 部署

### Docker 部署 (后端)

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN go build -o server cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Docker 部署 (前端)

```dockerfile
# Dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY frontend/ .
RUN npm install
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Kubernetes 部署

```yaml
# backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-admin-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: kube-admin-backend
  template:
    metadata:
      labels:
        app: kube-admin-backend
    spec:
      serviceAccountName: kube-admin-sa
      containers:
      - name: backend
        image: your-registry/kube-admin-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: kube-admin-secret
              key: jwt-secret
---
apiVersion: v1
kind: Service
metadata:
  name: kube-admin-backend
spec:
  selector:
    app: kube-admin-backend
  ports:
  - port: 8080
    targetPort: 8080
```

## 开发指南

### 添加新的 K8s 资源管理

1. **后端 Service 层**: 在 `backend/internal/service/` 创建新的 service 文件
2. **后端 API 层**: 在 `backend/internal/api/` 创建对应的 API handler
3. **路由注册**: 在 `backend/internal/router/router.go` 注册新路由
4. **前端 API**: 在 `frontend/src/apis/k8s.ts` 添加 API 调用
5. **前端页面**: 在 `frontend/src/views/k8s/` 创建新页面
6. **路由配置**: 在 `frontend/src/router/menus.ts` 添加路由

### 代码规范

- Go 代码遵循 `gofmt` 格式化
- Vue 代码使用 ESLint + Prettier
- Commit 信息遵循 Conventional Commits

## 安全建议

1. ✅ 使用 HTTPS
2. ✅ 修改默认 JWT 密钥
3. ✅ 实施 RBAC 权限控制
4. ✅ 定期更新依赖
5. ✅ 最小权限原则配置 ServiceAccount
6. ✅ 启用 API Rate Limiting

## 常见问题

### 1. 连接 K8s 集群失败
- 检查 kubeconfig 文件路径是否正确
- 验证 kubectl 是否能正常访问集群
- 检查 ServiceAccount 权限 (集群内部署)

### 2. 前端无法连接后端
- 确认后端服务已启动
- 检查 Vite proxy 配置
- 查看浏览器控制台错误

### 3. JWT Token 无效
- 检查 JWT_SECRET 是否一致
- Token 可能已过期,重新登录

## 贡献

欢迎提交 Issue 和 Pull Request!

## 许可证

MIT License

## 联系方式

如有问题,请提交 Issue。
