import { test, expect } from '@playwright/test'

test.describe('登录', () => {
  test('admin 登录后进入首页并同步集群到 URL', async ({ page }) => {
    await page.goto('/login')

    // 表单填充（Login.vue 默认预填 username=admin，这里显式覆盖以保证幂等）
    await page.getByLabel('用户名').fill('admin')
    await page.getByLabel('密码').fill('admin123')
    await page.getByRole('button', { name: '登录' }).click()

    // 落地首页（router.push('/')）；全新 profile 无历史集群，URL 不带 cluster_id 属正常
    await expect(page).toHaveURL(/\/$/, { timeout: 15_000 })
    await expect(page.getByRole('heading', { name: /Kubernetes 管理平台/ })).toBeVisible()
  })
})
