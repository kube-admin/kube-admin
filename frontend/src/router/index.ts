import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/Index.vue'
import nprogress from '@/utils/nprogress'
import menus from './menus'
import { initMenus } from './routers'

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

router.afterEach(() => {
  nprogress.done()
})

export default router