# 核心资源管理

Kube Admin 为高频资源提供了专门的管理页面，每个页面都支持命名空间筛选与自动刷新。

## 命名空间

「Kubernetes → Namespaces」：列表、创建、删除（系统命名空间受保护，不可删除）。

所有资源页面共享全局命名空间选择器，切换后自动同步到 Pod / Deployment / Service 等页面。

## Pods

「Kubernetes → Pods」：

- 查看容器状态、镜像、重启次数、资源 request/limit 与实时使用量
- **日志**：打开实时日志流（详见[终端与日志](./terminal-logs.md)）
- **终端**：进入容器交互式终端
- **删除**：重建 Pod
- 创建 Pod：通过粘贴 YAML 创建

## Deployments

「Kubernetes → Deployments」：

- 查看副本数（期望 / 就绪 / 已更新 / 可用）
- **扩缩容**：调整副本数
- **重启**：滚动重启（注入 `kubectl.kubernetes.io/restartedAt` 注解触发）
- 删除、通过 YAML 创建

## Services

「Kubernetes → Services」：类型、ClusterIP、ExternalIP、端口映射、Selector，删除与 YAML 创建。

## ConfigMaps / Secrets

- **ConfigMaps**：Key-Value 编辑器，支持创建/编辑/删除
- **Secrets**：Base64 自动编解码，支持明文/密文切换查看，创建/编辑/删除

::: warning Secret 安全
Secret 详情默认显示解码后内容以便排查，请注意操作环境的安全性。具有只读权限的角色亦可查看 Secret 明文，请按需分配权限。
:::

## Nodes

「Kubernetes → Nodes」：节点状态、IP、操作系统、kubelet 版本、容器运行时、容量与可分配资源、实时使用率、节点条件、标签。详见[仪表盘与监控](./dashboard.md)。

## 通用操作

| 操作 | 入口 |
|---|---|
| 命名空间切换 | 页面顶部命名空间下拉，全局同步 |
| 自动刷新 | 页面右上角「自动刷新」开关（5/10/30/60 秒） |
| 创建资源 | 各页面「创建」按钮，粘贴 YAML 提交 |

## 下一步

- [资源浏览器与 YAML](./resource-browser.md)：管理任意 K8s 资源
