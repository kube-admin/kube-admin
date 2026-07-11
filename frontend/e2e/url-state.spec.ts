import { test, expect } from './fixtures'

/**
 * URL 状态化回归保护。
 *
 * 历史 bug：App.vue.onMounted 晚于子组件 Header 挂载，URL 中的集群/命名空间被
 * fetchNamespaces 回退逻辑覆盖（namespace=kube-system 被改写为空或列表首个）。
 * 修复：URL→store 恢复提前到 setup 同步执行；拉取失败/空列表时不覆写已恢复值。
 *
 * 前置：cluster_id=1 可达且含 kube-system 命名空间。
 */
test.describe('URL 状态化', () => {
  test('分享链接恢复集群与命名空间上下文', async ({ authedPage: page }) => {
    await page.goto('/k8s/pods?cluster_id=1&namespace=kube-system')

    // 命名空间选择器恢复为 kube-system（而非被回退到列表首个）
    const nsSelector = page.locator('.namespace-selector .el-select')
    await expect(nsSelector).toContainText('kube-system', { timeout: 15_000 })

    // Pods 表格加载 kube-system 的资源（非空且命名空间列命中）
    const firstRow = page.locator('.el-table__row').first()
    await expect(firstRow).toContainText('kube-system', { timeout: 15_000 })

    // URL 未被覆写为空
    await expect(page).toHaveURL(/namespace=kube-system/)
  })

  test('命名空间拉取失败时不覆写 URL 恢复值', async ({ authedPage: page }) => {
    // 模拟命名空间接口失败：拦截返回 500
    await page.route('**/api/v1/namespaces**', (r) => r.fulfill({ status: 500, body: '{}' }))

    await page.goto('/k8s/pods?cluster_id=1&namespace=kube-system')

    // URL 仍保留 namespace=kube-system，不被清空（fetchNamespaces 失败保留逻辑）
    await expect(page).toHaveURL(/namespace=kube-system/)
  })

  test('切换命名空间写回 URL', async ({ authedPage: page }) => {
    await page.goto('/k8s/pods?cluster_id=1&namespace=kube-system')
    await expect(page.locator('.namespace-selector .el-select')).toContainText('kube-system', { timeout: 15_000 })

    // 打开下拉并选 default
    await page.locator('.namespace-selector .el-select').click()
    await page.getByRole('option', { name: 'default' }).click()

    // store→URL 同步：URL 出现 namespace=default
    await expect(page).toHaveURL(/namespace=default/, { timeout: 10_000 })
  })
})
