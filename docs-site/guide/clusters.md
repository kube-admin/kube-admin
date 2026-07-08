# 多集群管理

Kube Admin 支持接入多个 Kubernetes 集群并在界面内一键切换，所有集群凭据在数据库中加密存储。

## 接入集群

进入「Kubernetes → Clusters」，点击「创建 Clusters」。支持三种连接方式（任选其一）：

| 方式 | 字段 | 适用场景 |
|---|---|---|
| kubeconfig 内容（推荐） | `Config文件内容` | 粘贴整个 kubeconfig 文本，最通用 |
| kubeconfig 文件路径 | `Config文件路径` | 后端可访问的服务器文件路径 |
| 服务地址 + Token | `服务器地址` + `Token` | ServiceAccount Bearer Token 接入 |

::: tip 凭据安全
- 集群 `Token` 与 `kubeconfig 内容` 写入数据库前经 **AES-256-GCM** 加密，界面返回时脱敏（仅显示「已配置」标记）。
- 编辑集群时凭据字段留空表示**保留原值**，无需每次重新粘贴。
:::

## 测试连接

列表中点击「测试连接」可验证集群连通性并返回集群版本（基于已保存的凭据，后端解密后探测）。

::: warning TLS 校验
默认对集群 API Server 的 TLS 证书做校验。仅当集群使用自签证书且需放行时，设置环境变量 `TLS_SKIP_VERIFY=true`（仅限开发/受控环境）。
:::

## 切换集群

点击列表中的「切换」，当前集群会保存到本地并全局生效——所有资源页面（Pod / Deployment / 资源浏览器等）自动指向该集群，无需逐页切换。

## 多集群原理

- 每个请求携带 `cluster_id` 参数（前端自动注入），后端 `ClusterMiddleware` 据此从数据库取出集群配置、解密凭据，通过客户端缓存获取对应的 `client-go` 实例并注入请求上下文。
- 未指定 `cluster_id` 时使用启动时从本地 `KUBECONFIG` 创建的默认集群客户端。

## 下一步

- [仪表盘与监控](./dashboard.md)：查看集群资源使用率
- [核心资源](./workloads.md)：开始管理资源
