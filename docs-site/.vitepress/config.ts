import { defineConfig } from 'vitepress'

// 部署到子路径时通过 DOCS_BASE 覆盖，例如 /kube-admin/
const base = process.env.DOCS_BASE ?? '/'

export default defineConfig({
  title: 'Kube Admin',
  description: '基于 Vue 3 + Go + Kubernetes 的开源多集群管理平台',
  lang: 'zh-CN',
  lastUpdated: true,
  cleanUrls: false,
  base,
  ignoreDeadLinks: true,
  head: [
    ['meta', { name: 'theme-color', content: '#326CE5' }],
    ['link', { rel: 'icon', type: 'image/svg+xml', href: `${base}logo.svg` }]
  ],
  themeConfig: {
    logo: '/logo.svg',
    search: { provider: 'local' },
    nav: [
      { text: '首页', link: '/' },
      { text: '使用指南', link: '/guide/overview' },
      { text: '部署', link: '/guide/deploy' },
      { text: 'API', link: '/guide/api' }
    ],
    sidebar: {
      '/guide/': [
        {
          text: '开始',
          items: [
            { text: '项目介绍', link: '/guide/overview' },
            { text: '快速上手', link: '/guide/getting-started' }
          ]
        },
        {
          text: '集群与监控',
          items: [
            { text: '多集群管理', link: '/guide/clusters' },
            { text: '仪表盘与监控', link: '/guide/dashboard' }
          ]
        },
        {
          text: '资源管理',
          items: [
            { text: '核心资源', link: '/guide/workloads' },
            { text: '资源浏览器与 YAML', link: '/guide/resource-browser' },
            { text: '终端与日志', link: '/guide/terminal-logs' },
            { text: '事件', link: '/guide/events' }
          ]
        },
        {
          text: '安全与系统',
          items: [
            { text: '用户 / 权限 / 审计', link: '/guide/security' }
          ]
        },
        {
          text: '参考',
          items: [
            { text: '部署指南', link: '/guide/deploy' },
            { text: 'API 速查', link: '/guide/api' },
            { text: '常见问题', link: '/guide/faq' },
            { text: '更新日志', link: '/guide/changelog' }
          ]
        }
      ]
    },
    footer: {
      message: '基于 MIT 协议发布',
      copyright: 'Copyright © 2026 Kube Admin'
    },
    outline: { level: [2, 3] },
    docFooter: { prev: '上一页', next: '下一页' },
    lastUpdatedText: '最后更新'
  }
})
