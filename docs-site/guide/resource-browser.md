# 资源浏览器与 YAML 工作台

对于核心资源之外的所有 K8s 资源，Kube Admin 提供了**通用资源浏览器**与**YAML 工作台**，基于 dynamic client 统一管理，无需为每种资源单独开发。

## 资源浏览器

进入「Kubernetes → 资源管理」下的各菜单项，或在任一资源页顶部的资源类型下拉中切换。已内置覆盖：

- **工作负载**：StatefulSets / DaemonSets / ReplicaSets
- **网络**：Ingresses / NetworkPolicies
- **存储**：PersistentVolumeClaims / PersistentVolumes / StorageClasses
- **自动扩缩容**：HorizontalPodAutoscalers
- **访问控制**：ServiceAccounts / Roles / RoleBindings（ClusterRoles / ClusterRoleBindings 可在下拉切换）

每个资源类型支持：

- 列表查看（名称 / 命名空间 / 创建时间）
- 查看并编辑 YAML
- 删除
- 命名空间过滤（命名空间级资源）

::: tip 全资源统一能力
资源浏览器对所有 GVR（Group/Version/Resource）一视同仁——任何 CRD 或内置资源，只要集群识别，都可在下拉切换的 `key` 中加入并管理。前端定义集中在 `Resources.vue` 的 `RESOURCE_TYPES`。
:::

## YAML 工作台

点击「应用 YAML」打开编辑器，可粘贴任意资源 YAML 并提交：

- **创建或更新**：后端按 GVK 解析、通过 RESTMapper 定位 GVR，已存在则更新（保留 `resourceVersion`），不存在则创建
- **语法高亮与编辑**：Monaco Editor，支持自动补全与格式化
- **修改现有资源**：在列表点「YAML」打开详情，编辑后「应用修改」

## 适用场景

| 场景 | 推荐方式 |
|---|---|
| 日常扩缩容 / 重启 | 专门页面（Deployments 等） |
| 创建一个不常用资源（如 NetworkPolicy） | YAML 工作台 |
| 快速查看/编辑某 CRD 实例 | 资源浏览器 |
| 批量套用编排 | 建议使用 `kubectl`，本工作台面向单资源 |

## 安全提示

- YAML apply 与删除属于写操作，需 `admin` / `operator` 角色。
- 所有写操作会记录到审计日志（详见[用户/权限/审计](./security.md)）。

## 下一步

- [终端与日志](./terminal-logs.md)
- [事件](./events.md)
