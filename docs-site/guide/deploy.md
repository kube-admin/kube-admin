# 部署指南

Kube Admin 提供三种部署方式，部署文件统一在 `deploy/` 目录：

| 方式 | 适用 | 入口 |
|------|------|------|
| **Docker Compose** | 单机 / Demo / 个人 | `make docker-up`（→ `deploy/deploy.sh`） |
| **Kubernetes 原生清单** | 集群内直接 apply | `kubectl apply -f deploy/k8s/` |
| **Helm Chart** | 多环境 / 参数化 | `helm install ... deploy/helm/kube-admin` |

::: tip 镜像来源
三种方式共用 CI 推送到 ghcr.io 的镜像（`ghcr.io/kube-admin/kube-admin-backend` / `-frontend`）。push 到 main 或打 tag 后由 `Release Images` workflow 自动构建。
:::

## 前置：密钥

生产环境必须设置强随机密钥：

```bash
openssl rand -base64 32   # 用作 JWT_SECRET
openssl rand -base64 32   # 用作 ENCRYPT_KEY
```

::: danger ENCRYPT_KEY 不可丢失
`ENCRYPT_KEY` 是解密集群凭据的唯一依据。丢失将导致已存的集群凭据无法解密。请妥善备份。
:::

## 一、Docker Compose（单机）

```bash
cd deploy
cp .env.example .env
# 编辑 .env：填写 JWT_SECRET / ENCRYPT_KEY（必填）；按需设 DB_DRIVER/DB_DSN
bash deploy.sh up          # 或退回根目录 make docker-up
```

`deploy.sh` 会校验密钥、构建并启动、轮询 `/healthz` 健康检查。

- 前端：http://localhost
- 后端：http://localhost:8080
- 登录：`admin / admin123`（首次登录后立即改密）

其他命令：`bash deploy.sh down | restart | logs [backend|frontend] | ps`。

## 二、Kubernetes（集群内部署）

manifests 已在仓库内，直接 apply：

```bash
# 1. 必改：编辑 deploy/k8s/secret.yaml，替换 jwt-secret / encrypt-key 占位
kubectl apply -f deploy/k8s/
```

镜像默认引用 ghcr.io 公开镜像。默认集群通过 **ServiceAccount** 自动接入（后端 kubeconfig 为空时走 InClusterConfig），无需额外配置即可管理当前集群。

未启用 Ingress 时访问：
```bash
kubectl -n kube-admin port-forward svc/kube-admin-frontend 8080:80
```

### 集群内 vs 管理外部集群

| 场景 | 方式 |
|------|------|
| 管当前集群 | 自动用 Pod 的 ServiceAccount，零配置 |
| 管外部集群 | 登录后界面「集群管理」→ 填 kubeconfig/Token（凭据 AES-256-GCM 加密入库） |

## 三、Helm

```bash
helm install kube-admin deploy/helm/kube-admin \
  --create-namespace -n kube-admin \
  --set secrets.jwtSecret="$(openssl rand -base64 32)" \
  --set secrets.encryptKey="$(openssl rand -base64 32)"
```

常用参数：`ingress.enabled`、`database.driver`、`backend.replicas`、`serviceAccount.clusterRole.wide`。详见 `deploy/helm/kube-admin/README.md`。

::: warning 生产密钥
Helm values 中的密钥仅用于快速开始。生产环境建议用 **ExternalSecrets Operator** 或 **sealed-secrets** 外部化管理，避免密钥进 Git。
:::

## 数据库

默认 SQLite（单机/单副本，数据卷持久化）。接入现有 MySQL/PostgreSQL：

```bash
# 以 compose 为例（.env）：
DB_DRIVER=mysql
DB_DSN=user:pass@tcp(host:3306)/kubeadm?charset=utf8mb4&parseTime=True&loc=Local
# 或 postgres：
DB_DRIVER=postgres
DB_DSN=host=pg user=kubeadmin password=secret dbname=kubeadm port=5432 sslmode=disable TimeZone=Asia/Shanghai
```

::: tip 多副本
SQLite 为单写者，多副本部署请切换到外部 MySQL/PostgreSQL，并将 backend 副本数调大、策略改 RollingUpdate（K8s/Helm 均已预留参数）。
:::

`AutoMigrate` 会在首次启动自动建表，无需手动 DDL。

## 运行时调优

- `K8S_REQUEST_TIMEOUT`：k8s API 单次请求超时（秒，默认 `10`）。某集群不可达时让请求快速失败，而非 client-go 默认挂起 30s 拖垮前端体验。网络较差或多跳代理环境可适当调大。

## 健康检查与优雅关闭

- `GET /healthz`：存活探针
- `GET /readyz`：就绪探针
- 收到 SIGINT/SIGTERM 后等待最多 10s 处理完在途请求再退出

## 下一步

- [API 速查](./api.md)
- [常见问题](./faq.md)
