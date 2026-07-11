import { test as base, request, expect, type Page } from '@playwright/test'

/**
 * 复用登录态的测试夹具。
 *
 * 通过 API 登录拿到 token/user（模块级缓存，全文件只登一次），
 * 再用 addInitScript 在每次导航前注入 localStorage——等价于"已登录用户打开分享链接"。
 */
const BASE = 'http://localhost:3000'
const USERNAME = process.env.E2E_USER ?? 'admin'
const PASSWORD = process.env.E2E_PASS ?? 'admin123'

type Creds = { token: string; userJson: string }

let _creds: Creds | null = null

async function creds(): Promise<Creds> {
  if (_creds) return _creds
  const ctx = await request.newContext({ baseURL: BASE })
  const r = await ctx.post('/api/v1/auth/login', {
    data: { username: USERNAME, password: PASSWORD },
  })
  expect(r.ok(), `登录请求失败: ${r.status()} ${r.statusText()}`).toBeTruthy()
  const body = await r.json()
  expect(body?.data?.token, '响应缺少 token').toBeTruthy()
  _creds = { token: body.data.token, userJson: JSON.stringify(body.data.user) }
  await ctx.dispose()
  return _creds
}

/** authedPage：已注入登录态的 page，用于需要鉴权的用例。 */
export const test = base.extend<{ authedPage: Page }>({
  authedPage: async ({ page }, use) => {
    const c = await creds()
    await page.addInitScript((c: Creds) => {
      localStorage.setItem('token', c.token)
      localStorage.setItem('user', c.userJson)
    }, c)
    await use(page)
  },
})

export { expect }
