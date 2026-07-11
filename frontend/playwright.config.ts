import { defineConfig, devices } from '@playwright/test'

/**
 * Playwright 本地 E2E 配置。
 *
 * 定位：兜底"集成型回归"——跨路由/store/组件/API 的链路（如 URL 状态化恢复），
 * 这类缺陷单测难以复现。仅在本地开发栈（vite :3000 + backend :8080 + 可达集群）上运行。
 *
 * 不入 CI：CI 无真实集群，E2E 跑通也不代表线上可用；待引入 kind 集群后再开 CI E2E job。
 *
 * 前置：本地已起 `npm run dev`（:3000）与后端（:8080），且集群 cluster_id=1 可达。
 */
export default defineConfig({
  testDir: './e2e',
  // 集成测试涉及共享后端/集群状态，串行执行避免竞态
  fullyParallel: false,
  workers: 1,
  retries: 0,
  reporter: 'list',
  use: {
    baseURL: 'http://localhost:3000',
    trace: 'on-first-retry',
    actionTimeout: 10_000,
    navigationTimeout: 15_000,
    // 登录态注入见 e2e/fixtures.ts
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:3000',
    // 复用已运行的 vite，避免重复起服务；本地无 vite 时自动拉起
    reuseExistingServer: true,
    timeout: 60_000,
  },
})
