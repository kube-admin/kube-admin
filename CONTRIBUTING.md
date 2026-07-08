# 贡献指南

感谢你关注 Kube Admin！欢迎通过 Issue 与 Pull Request 参与贡献。

## 开发环境

```bash
git clone <repo-url>
cd kube-admin

# 后端
cd backend && go mod tidy && go run cmd/main.go

# 前端
cd frontend && npm install && npm run dev
```

## 提交规范

- Commit 信息遵循 [Conventional Commits](https://www.conventionalcommits.org/)，例如：
  - `feat(pod): 支持日志 WebSocket 流式`
  - `fix(auth): 修复 token 解析`
  - `docs: 更新 README`
  - `refactor(service): 抽取通用资源 service`
- Go 代码需通过 `gofmt` 与 `go vet`；新增逻辑请补充单元测试。
- Vue/TS 代码需通过 `npm run lint`（vue-tsc 类型检查）；新增组件/工具函数鼓励补充 vitest 测试。

## 新增 K8s 资源管理

优先复用**通用资源 service**（`backend/internal/service/resource.go`，基于 dynamic client）：

- 列表 / 详情 / 删除 / YAML apply 已通过通用接口 `/api/v1/resources` 支持，前端在 `Resources.vue` 中以 GVR 切换即可覆盖新资源，无需重复编码。
- 仅当需要强类型字段聚合（如监控、状态计算）时，才新增专属 service。

## CI

所有 PR 会自动运行 GitHub Actions：后端 `go build / vet / test`，前端 `lint / test / build`，均需通过。

## 行为准则

请保持友善与专业，尊重每一位贡献者。
