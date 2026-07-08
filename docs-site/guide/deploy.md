# 部署指南

Kube Admin 支持本地、Docker、Kubernetes 三种部署方式。

## 前置：配置密钥

生产环境必须设置：

```bash
export JWT_SECRET="<强随机字符串>"
export ENCRYPT_KEY="<强随机字符串>"   # 集群凭据加密密钥
# export TLS_SKIP_VERIFY=false        # 默认即校验 TLS
```

::: danger 密钥不可丢失
`ENCRYPT_KEY` 是解密集群凭据的唯一依据。密钥丢失将导致已存的集群凭据无法解密、集群连接失败。请妥善备份。
:::

## Docker Compose

```bash
# 编辑 .env 设置 JWT_SECRET / ENCRYPT_KEY
docker-compose up -d --build
```

默认暴露前端 `:80`（或 :3000）、后端 `:8080`。前端容器（nginx）反向代理 `/api` 到后端。

## Kubernetes 部署

### 1. 后端 Deployment

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
      serviceAccountName: kube-admin-sa   # 集群内访问当前集群
      containers:
      - name: backend
        image: your-registry/kube-admin-backend:latest
        ports: [{ containerPort: 8080 }]
        env:
        - { name: JWT_SECRET,  valueFrom: { secretKeyRef: { name: kube-admin-secret, key: jwt-secret } } }
        - { name: ENCRYPT_KEY, valueFrom: { secretKeyRef: { name: kube-admin-secret, key: encrypt-key } } }
        - { name: GIN_MODE, value: "release" }
        livenessProbe:  { httpGet: { path: /healthz, port: 8080 }, initialDelaySeconds: 10, periodSeconds: 20 }
        readinessProbe: { httpGet: { path: /readyz,   port: 8080 }, initialDelaySeconds: 5,  periodSeconds: 10 }
```

### 2. 前端

前端构建产物为静态文件，用 nginx 托管并将 `/api`、WebSocket 反向代理到后端 Service。

### 3. ServiceAccount 权限

集群内部署时后端自动使用 ServiceAccount 访问**当前集群**。请按最小权限原则为其绑定 Role/ClusterRole（如只读 + 指定资源的写权限）。若还需管理**外部集群**，在界面通过 kubeconfig/Token 接入即可（凭据加密存储）。

### 4. 数据库

默认使用 SQLite（`DB_PATH`）。多副本部署时，请替换为外部 MySQL/PostgreSQL（更换 GORM 驱动即可），避免多实例写同一文件。

## 健康检查

- `GET /healthz`：存活探针
- `GET /readyz`：就绪探针

后端支持优雅关闭：收到 `SIGINT` / `SIGTERM` 后等待最多 10 秒处理完在途请求再退出。

## 下一步

- [API 速查](./api.md)
- [常见问题](./faq.md)
