import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/Index.vue'
import nprogress from '@/utils/nprogress'
import menus from './menus'
import { initMenus } from './routers'
import { useNamespaceStore } from '@/stores/namespace'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/',
      name: 'layout',
      component: Layout,
      children: menus
    },
    {
      // Pod 终端全屏页：顶级路由，脱离 layout（无侧边栏），纯全屏交互
      path: '/k8s/pods/:name/terminal',
      name: 'podTerminal',
      component: () => import('../views/k8s/PodTerminal.vue'),
      meta: {}
    },
    {
      path: '/:catchAll(.*)',
      redirect: '/404',
      meta: {}
    }
  ]
})

// 添加路由守卫
router.beforeEach(async (to, from, next) => {
  nprogress.start()
  
  // 如果访问的是登录页，直接放行
  if (to.path === '/login') {
    next()
    return
  }
  
  // 检查是否有 token
  const token = localStorage.getItem('token')
  if (!token) {
    // 没有 token，跳转到登录页
    next('/login')
    return
  }
  
  // 如果有 token 但还没有初始化菜单，初始化菜单
  if (router.options.routes.length <= 3) { // 只有默认的三个路由
    try {
      const dynamicMenus = await initMenus()
      // 添加动态路由
      dynamicMenus.forEach(route => {
        router.addRoute('layout', route)
      })
    } catch (error) {
      console.error('Failed to initialize menus:', error)
    }
  }
  
  next()
})

// URL 状态化守卫：菜单跳转等不带 query 时，从 store 补 cluster_id/namespace，
// 使任何导航后的 URL 都能还原集群/命名空间上下文（分享链接可复现）。
// watch(store) 只在 store 变化时写 URL，导航不改 store，故需守卫兜底。
router.beforeEach((to, from, next) => {
  if (to.path === '/login' || to.path === '/404') {
    next()
    return
  }
  const nsStore = useNamespaceStore()
  const query: Record<string, any> = { ...to.query }
  let changed = false
  if (nsStore.currentClusterId && !query.cluster_id) {
    query.cluster_id = String(nsStore.currentClusterId)
    changed = true
  }
  if (nsStore.currentNamespace && !query.namespace) {
    query.namespace = nsStore.currentNamespace
    changed = true
  }
  if (changed) {
    next({ path: to.path, query, replace: true })
    return
  }
  next()
})

router.afterEach(() => {
  nprogress.done()
})

export default router