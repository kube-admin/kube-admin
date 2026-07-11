# Kube Admin Helm Chart

基于 Vue3 + Go 的多集群 Kubernetes 管理平台。

## 快速安装

```bash
# 1. 安装（自动创建 namespace，密钥用强随机值）
helm install kube-admin deploy/helm/kube-admin \
  --create-namespace -n kube-admin \
  --set secrets.jwtSecret="$(openssl rand -base64 32)" \
  --set secrets.encryptKey="$(openssl rand -base64 32)"

# 2. 访问（未启用 Ingress 时 port-forward）
kubectl -n kube-admin port-forward svc/kube-admin-frontend 8080:80
# 浏览器打开 http://localhost:8080 ，登录 admin / admin123
```

## 常用配置

| 场景 | 参数 |
|------|------|
| 启用 Ingress | `--set ingress.enabled=true --set ingress.host=kube-admin.example.com` |
| 切外部 MySQL（可多副本） | `--set database.driver=mysql --set secrets.dbDSN='user:pass@tcp(host:3306)/kubeadm?...' --set backend.strategy=RollingUpdate --set backend.replicas=2` |
| 收窄权限（禁用通用资源 apply） | `--set serviceAccount.clusterRole.wide=false` |
| 自定义镜像 | `--set image.backend.repository=... --set image.backend.tag=...` |

## 验证

```bash
helm lint deploy/helm/kube-admin
helm template kube-admin deploy/helm/kube-admin -n kube-admin
```

## 生产密钥建议

values 中的密钥仅用于快速开始。生产环境建议用 **ExternalSecrets Operator** 或 **sealed-secrets** 外部化管理，避免密钥进入 Git 或 Helm values 文件。

## 卸载

```bash
helm uninstall kube-admin -n kube-admin
kubectl delete namespace kube-admin   # 可选，连同数据一并清理
```
