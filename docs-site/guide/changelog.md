# 更新日志

遵循 [Keep a Changelog](https://keepachangelog.com/) 风格。

## [0.2.0] - 2026-07-11

### Added
- **资源详情页体系**：Pod / Deployment / StatefulSet / DaemonSet / Service 详情页（概览 + YAML + Events 标签），列表↔详情双向链接（Pod 显示所属工作负载、Deployment 列关联 Pods）
- **通用 YAML 抽屉**（`<YamlDrawer>`）：任意资源在列表页直接查看/编辑 YAML（apply 保存），统一"查看/编辑"操作
- **列表工具栏**（`<ListToolbar>`）：统一所有列表页 header（标题 + 刷新 + 自动刷新 + 过滤 slot），消除风格不一致
- **列表搜索/筛选**：Pod（名称+状态）、Service（名称+类型）、其余资源名称模糊搜索（前端过滤）
- **workload 扩缩容/滚动重启**：StatefulSet / DaemonSet / ReplicaSet（通用 `ResourceService.Scale/Restart`，dynamic client）
- **Pod 列表 CPU/内存列**：汇总各容器 metrics，与 Dashboard 同源
- **URL 状态化**：Header 集群/命名空间同步到 URL，分享链接打开即恢复对应上下文
- **多数据库驱动**：`DB_DRIVER` 支持 sqlite（默认）/ mysql / postgres，`DB_DSN` 配置连接串；`AutoMigrate` 跨库通用，多副本部署可切外部 MySQL/PostgreSQL
- **E2E 测试框架**：Playwright 接入（`npm run test:e2e`），覆盖登录与 URL 状态化恢复链路，兜底集成型回归；本地运行，CI 接入待 kind 集群

### Changed
- **Pod 终端改为新标签打开**：独立于列表页，避免误关列表断终端；支持多开
- **Dashboard CPU/内存格式化**：k8s Quantity → 人可读（核 / Gi）
- **列表页命名空间选择器统一到 Header 全局**：去掉 5 个页内冗余选择器，联动刷新；Resources 页类型选择器移除（由路由决定）

### Fixed
- **metrics-server 在受限网络部署**：hostNetwork 绕开损坏的 calico pod 网络 + 阿里云源经本地 registry（hub.wang.dd:5000）分发 + 手动 EndpointSlice
- **monaco-editor `factory.create` 报错**：`optimizeDeps.exclude` + `MonacoEnvironment.getWorker` 配置
- **Header 添加集群下拉不刷新**：`clustersChanged` 事件总线
- **详情页 StatusTag 换行**：穿透覆盖 el-tag__content 固定宽度
- **URL 状态化恢复失效**：`main.ts` 未 `await router.isReady()` 即挂载，`App.vue` setup 执行时初始路由尚未解析、`route.query` 为空，URL->store 恢复一次性读取到空值从未生效（此前被 localStorage 残留 `currentCluster` 假阳性掩盖）。改为响应式 `watch(route.query)` 在路由解析后恢复；`fetchNamespaces` 拉取失败或列表为空时保留已恢复值，不再覆写。新增 Playwright E2E（`npm run test:e2e`）覆盖该回归