# kube-admin Backend

## 运行

```bash
# 安装依赖
go mod tidy

# 运行
go run cmd/main.go

# 编译
go build -o kube-admin-server cmd/main.go
```

## 环境变量

- `PORT`: 服务端口,默认 8080
- `KUBECONFIG`: kubeconfig 文件路径,默认 ~/.kube/config
- `JWT_SECRET`: JWT 密钥,生产环境必须修改

## 目录结构

```
backend/
├── cmd/              # 应用入口
│   └── main.go
├── config/           # 配置管理
│   └── config.go
├── internal/         # 内部代码
│   ├── api/         # API handlers
│   │   ├── auth.go
│   │   ├── deployment.go
│   │   ├── namespace.go
│   │   └── pod.go
│   ├── middleware/  # 中间件
│   │   ├── auth.go
│   │   └── cors.go
│   ├── model/       # 数据模型
│   │   ├── k8s.go
│   │   ├── response.go
│   │   └── user.go
│   ├── router/      # 路由
│   │   └── router.go
│   └── service/     # 业务逻辑
│       ├── deployment.go
│       ├── namespace.go
│       ├── pod.go
│       └── service.go
└── pkg/             # 可复用包
    └── k8s/
        └── client.go
```
