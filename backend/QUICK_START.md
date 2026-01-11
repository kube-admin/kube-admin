# 后端服务快速启动指南

## ✅ 后端服务已成功修复并启动!

### 🔧 修复的问题

1. **语法错误修复**
   - ✅ 修复了 `logger.go` 中重复的 `package logger` 声明
   - ✅ 修复了 `pod.go` 文件末尾的多余 `package service` 声明
   
2. **缺失文件补充**
   - ✅ 创建了 `internal/api/service.go` (Service API 处理器)
   
3. **路由配置完善**
   - ✅ 添加了所有缺失的 Service 和 API 初始化
   - ✅ 注册了所有 30+ 个 API 路由

---

## 🚀 启动方式

### 方式一：使用启动脚本（推荐）

```bash
cd backend
./start.sh
```

### 方式二：直接运行

```bash
cd backend
go run cmd/main.go
```

### 方式三：编译后运行

```bash
cd backend
go build -o kube-admin-server cmd/main.go
./kube-admin-server
```

---

## 📋 环境要求

### 必需
- ✅ Go 1.21 或更高版本
- ✅ Kubernetes 集群（本地 minikube 或远程集群）
- ✅ kubeconfig 配置文件（通常在 `~/.kube/config`）

### 可选
- 环境变量配置：
  ```bash
  export PORT=8080                    # 服务端口（默认: 8080）
  export KUBECONFIG=~/.kube/config   # K8s 配置文件路径
  export JWT_SECRET=your-secret-key  # JWT 密钥
  ```

---

## ✅ 服务状态检查

### 1. 检查服务是否启动

看到以下输出说明启动成功：

```
✅ Successfully connected to Kubernetes cluster
✅ Server starting on :8080
✅ Listening and serving HTTP on :8080
```

### 2. 测试 API 接口

```bash
# 测试登录接口
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取 Pod 列表（需要先登录获取 token）
curl http://localhost:8080/api/v1/pods?namespace=default \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 📊 已注册的 API 路由

### 认证接口
- ✅ `POST /api/v1/auth/login` - 用户登录
- ✅ `GET /api/v1/auth/user` - 获取用户信息

### Dashboard
- ✅ `GET /api/v1/dashboard/stats` - 集群统计信息

### Namespace
- ✅ `GET /api/v1/namespaces` - 列表
- ✅ `POST /api/v1/namespaces` - 创建
- ✅ `DELETE /api/v1/namespaces/:name` - 删除

### Node
- ✅ `GET /api/v1/nodes` - 列表
- ✅ `GET /api/v1/nodes/:name` - 详情

### Pod
- ✅ `GET /api/v1/pods` - 列表
- ✅ `GET /api/v1/pods/:name` - 详情
- ✅ `DELETE /api/v1/pods/:name` - 删除
- ✅ `GET /api/v1/pods/:name/logs` - 日志

### Deployment
- ✅ `GET /api/v1/deployments` - 列表
- ✅ `GET /api/v1/deployments/:name` - 详情
- ✅ `DELETE /api/v1/deployments/:name` - 删除
- ✅ `PUT /api/v1/deployments/:name/scale` - 扩缩容
- ✅ `PUT /api/v1/deployments/:name/restart` - 重启

### Service
- ✅ `GET /api/v1/services` - 列表
- ✅ `GET /api/v1/services/:name` - 详情
- ✅ `DELETE /api/v1/services/:name` - 删除

### ConfigMap
- ✅ `GET /api/v1/configmaps` - 列表
- ✅ `GET /api/v1/configmaps/:name` - 详情
- ✅ `POST /api/v1/configmaps` - 创建
- ✅ `PUT /api/v1/configmaps/:name` - 更新
- ✅ `DELETE /api/v1/configmaps/:name` - 删除

### Secret
- ✅ `GET /api/v1/secrets` - 列表
- ✅ `GET /api/v1/secrets/:name` - 详情
- ✅ `POST /api/v1/secrets` - 创建
- ✅ `PUT /api/v1/secrets/:name` - 更新
- ✅ `DELETE /api/v1/secrets/:name` - 删除

---

## 🔍 故障排查

### 问题1：无法连接到 K8s 集群

**错误信息**:
```
Failed to create k8s client: ...
```

**解决方案**:
```bash
# 检查 kubeconfig 文件是否存在
ls ~/.kube/config

# 测试 kubectl 命令
kubectl cluster-info

# 如果使用 minikube
minikube start
```

### 问题2：端口被占用

**错误信息**:
```
bind: address already in use
```

**解决方案**:
```bash
# 查看占用端口的进程
lsof -i :8080

# 或者更改端口
export PORT=8081
./start.sh
```

### 问题3：依赖下载失败

**解决方案**:
```bash
# 配置 Go 代理
export GOPROXY=https://goproxy.cn,direct

# 重新下载依赖
go mod tidy
```

---

## 📝 默认配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| 端口 | 8080 | HTTP 服务端口 |
| Kubeconfig | `~/.kube/config` | K8s 配置文件路径 |
| JWT Secret | `your-secret-key-change-in-production` | JWT 签名密钥（生产环境需修改） |
| 默认用户名 | `admin` | 登录用户名 |
| 默认密码 | `admin123` | 登录密码 |

---

## 🎯 下一步

1. ✅ 后端服务已启动
2. 📱 启动前端服务：
   ```bash
   cd ../frontend
   npm install
   npm run dev
   ```
3. 🌐 访问系统：`http://localhost:3000`

---

## 📞 技术支持

如遇问题，请检查：
1. Go 版本是否 >= 1.21
2. Kubernetes 集群是否可访问
3. 端口是否被占用
4. 依赖是否完整安装

**日志位置**: 控制台输出

---

## 🎉 成功标志

看到以下输出表示服务正常运行：

```
✅ Successfully connected to Kubernetes cluster
✅ Server starting on :8080
✅ [GIN-debug] Listening and serving HTTP on :8080
✅ 30+ API routes registered
```

**现在可以使用前端访问后端 API 了! 🚀**
