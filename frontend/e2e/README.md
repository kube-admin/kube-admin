# E2E 测试（Playwright）

定位：兜底"集成型回归"--跨路由/store/组件/API 的链路（如 URL 状态化恢复），这类缺陷单测难以复现。

## 运行前置

E2E 在真实浏览器里跑完整链路，需要本地开发栈就绪：

1. **后端**（:8080）：`cd backend && PORT=8080 go run ./cmd/main.go`
2. **前端**（:3000）：`cd frontend && npm run dev`（playwright.config.ts 会复用已运行的 vite）
3. **可达集群**：默认 `cluster_id=1`（"本地"）可列出命名空间，且含 `kube-system`
4. **默认账号**：`admin/admin123`（可用 `E2E_USER`/`E2E_PASS` 环境变量覆盖）

## 运行

```bash
cd frontend
npm run test:e2e              # 跑全部
npx playwright test --headed  # 可视化调试
npx playwright test --ui      # 交互式 UI 模式
```

首次运行需安装浏览器：`npx playwright install chromium`（二进制装在用户级缓存，跨项目共享）。

## 用例

- `login.spec.ts`：登录 -> 落地首页
- `url-state.spec.ts`：URL 状态化恢复（分享链接打开即恢复集群/命名空间）、拉取失败不覆写、切换写回 URL
  - 保护历史回归：App.vue setup 在路由解析前运行导致 URL->store 恢复失效（被 localStorage 残留假阳性掩盖）

## 为什么不在 CI 跑

CI（GitHub-hosted runner）无真实 k8s 集群，E2E 跑通也不代表线上可用；起 kind 集群会显著拉长 CI 时间。当前阶段 E2E 为本地/手动验证工具，待 v1.0 前引入 kind-based CI job。
