# 命名空间管理

## 列表与统计

「Kubernetes → 集群资源 → Namespaces」页面展示当前集群所有命名空间，包含状态、年龄、Pods / Deployments / Services / ConfigMaps 资源统计。

支持创建、删除、查看/编辑 YAML。

## Terminating 卡死诊断

namespace 删除时可能卡在 `Terminating` 状态。展开行可查看根因：

- **卡住时长**：基于 `deletionTimestamp` 计算
- **Finalizers**：阻塞删除的 finalizer 列表
- **根因 Conditions**（`status=True` 的 condition）：
  - `NamespaceDeletionDiscoveryFailure`：API 发现失败（常见为失效 APIService，如 calico-aggregator）
  - `NamespaceFinalizersRemaining`：自定义 finalizer 卡住
  - `NamespaceContentRemaining` / `NamespaceDeletionContentFailure`：残留资源

状态列 `Terminating` 标签旁有橙色提示图标，引导展开查看根因。

### 根因与处理动作

| 根因 | 处理动作 |
|---|---|
| DiscoveryFailed（失效 APIService） | 「查看失效 APIService」打开抽屉，删除失效 APIService（治本，删除后 namespace 自动完成删除） |
| FinalizersRemaining | 「强制清理 finalizers」清空 finalizers 强制删除（危险操作，二次确认） |
| ContentRemaining / ContentFailure | 仅提示，需手动检查残留资源 |

## 失效 APIService 抽屉

DiscoveryFailed 根因下，点「查看失效 APIService」打开抽屉：

- 列出集群所有 `Available != True` 的 APIService（名称、后端服务、状态 reason、年龄）
- 单个删除（集群级危险操作，`el-popconfirm` 二次确认）
- 删除失效 APIService 后，被卡的 namespace 会自动完成删除

::: tip 治本 vs 治标
- **治本**：删除失效 APIService（或修复其后端服务），discovery 恢复，所有被卡的 namespace 自动完成删除
- **治标**：强制清理 finalizers，仅删除当前 namespace，集群 APIService 问题仍在，其他 namespace 删除仍会卡
:::

## 所有命名空间

Header namespace select 支持「所有命名空间」选项：

- 选择后，Pods / Deployments / Services 等列表页显示所有命名空间的数据
- 默认显示 `default` 命名空间；URL/selector 为空串时恢复 default

## URL 状态化

集群与命名空间选择同步到 URL query（`cluster_id` + `namespace`）：

- 点菜单跳转任何页面，URL 自动带当前集群/命名空间（路由守卫兜底注入）
- 复制 URL 分享，打开即恢复对应集群/命名空间上下文
- `App.vue` 监听 `route.query` 响应式恢复，`watch(store)` 写回 URL
