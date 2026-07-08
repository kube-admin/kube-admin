# 常见问题

## 启动与连接

### 后端启动报 `bind: address already in use`（8080 被占）

8080 端口被其他程序占用（可用 `lsof -i :8080` 排查）。两种处理：

1. 停掉占用端口的进程；
2. 让 Kube Admin 用其他端口：`PORT=8090 make backend`，并把 `frontend/vite.config.ts` 中 proxy 的 `target` 改为 `http://localhost:8090`。

### 后端日志显示 `Failed to create default k8s client`

本机 `~/.kube/config` 不可用或无集群。这不影响启动——可在「集群管理」页面手动接入集群。

### 前端登录提示 401

- 确认后端已启动且端口正确；
- 确认 `JWT_SECRET` 前后端一致（前端无需配置，由后端签发）；
- 默认账号 `admin / admin123`，注意是否已被修改。

## 监控

### 节点/Pod 使用率为空

实时使用率依赖 [metrics-server](https://github.com/kubernetes-sigs/metrics-server)。未安装时使用率显示为空，Dashboard 会出现黄色提示。安装方法见[仪表盘与监控](./dashboard.md)。

### Dashboard 数据不刷新

页面每 30 秒自动刷新一次；也可手动点击「刷新」。若持续无数据，检查集群连接与 ServiceAccount 权限。

## 终端与日志

### 终端连不上 / 无响应

- 确认 Pod 处于 `Running`；
- 极简镜像（如 `scratch` / distroless）可能不含 `/bin/sh`，无法 exec；
- 检查反向代理是否放行 WebSocket 升级。

### 日志流中断

网络抖动会断开 WebSocket，重新打开日志即可。`follow=false` 的一次性日志不受影响。

## 安全

### 忘记 admin 密码

若无其他 admin 用户，需直接操作数据库重置密码（bcrypt），或删除 `data/kubeadm.db` 后重启（会丢失用户与集群配置，集群需重新接入）。

### 更换 `ENCRYPT_KEY` 后集群连不上

`ENCRYPT_KEY` 变更后，旧密钥加密的集群凭据无法解密。需在「集群管理」重新录入凭据。生产环境请在初始化时确定密钥并备份。

## 资源管理

### 资源浏览器里找不到某 CRD

资源类型下拉集中了常用资源。要管理其他 CRD，在前端 `Resources.vue` 的 `RESOURCE_TYPES` 中追加一项（指定 `group/version/resource`）即可，无需后端改动。

### apply YAML 报 `无法识别资源类型`

后端通过 RESTMapper 解析 GVK。若集群未安装对应 CRD，则无法识别。先在集群中安装 CRD 再 apply。

## 其他

### 数据存储在哪里

默认 SQLite，路径由 `DB_PATH`（默认 `data/kubeadm.db`）决定。多副本生产部署请改用外部数据库。

### 如何参与贡献

参见 [CONTRIBUTING.md](https://github.com/)（仓库根）。
