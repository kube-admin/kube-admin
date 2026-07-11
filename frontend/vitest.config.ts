import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  test: {
    environment: 'jsdom',
    globals: true,
    // 排除 Playwright E2E 用例（由 `npm run test:e2e` 独立运行）
    exclude: ['**/node_modules/**', '**/dist/**', 'e2e/**', 'playwright.config.ts']
  }
})
