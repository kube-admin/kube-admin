# 更新日志

完整更新日志见文档站：[docs-site/guide/changelog.md](docs-site/guide/changelog.md)

遵循 [Keep a Changelog](https://keepachangelog.com/) 风格。

## [0.3.0] - 2026-07-12

### Added
- Namespace Terminating 诊断：展开行显示根因 conditions、finalizers、卡住时长；状态列提示图标
- 失效 APIService 抽屉：DiscoveryFailed 根因一键查看集群失效 APIService + 删除
- namespace「所有命名空间」选项（all 哨兵）：select 支持 all 查所有 ns 数据
- URL 状态化守卫：导航自动带 cluster_id/namespace query

### Changed
- namespace 默认改回 default（保留 all 选项）
- APIService 访问改用 kube-aggregator clientset（Client 加 AggregatorClient 字段）

### Fixed
- node 等集群级资源详情 404：namespace=all 被当 namespace 作用域，resource.go 将 all 视为 all-namespaces
- 多集群路径 APIService 接口 panic：Manager.createClient 补 AggregatorClient 构造

## [0.2.0] - 2026-07-11

资源详情页体系、通用 YAML 抽屉、列表工具栏、URL 状态化、多数据库驱动、E2E 测试框架等。详见 [文档站 changelog](docs-site/guide/changelog.md)。
