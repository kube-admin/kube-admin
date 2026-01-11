import { reactive } from 'vue'
import { RouteRecordRaw } from 'vue-router'
import { fetchMenus, menuItemData } from '@/apis/user/info'

const components = import.meta.glob('@/views/**/*.vue')

const ruleForm = reactive({
  userId: ''
})

export var generateRoute = (item: menuItemData) => {
  let menu = {
    path: item.path,
    name: item.name,
    component: item.component ? loadView(item.component) : null,
    meta: item.meta,
    children: []
  } as RouteRecordRaw
  
  if (item.children && item.children.length > 0) {
    item.children.forEach((item1) => {
      menu.children?.push(generateRoute(item1))
    })
  }
  
  return menu
}

export const loadView = (view: string | undefined) => {
  if (view === undefined || view === null) return null
  
  // 处理不同的组件路径格式
  let componentPath = ''
  if (view.startsWith('@/views/')) {
    componentPath = view
  } else if (view.startsWith('/src/views/')) {
    componentPath = view
  } else if (view.startsWith('views/')) {
    componentPath = `/src/${view}.vue`
  } else {
    componentPath = `/src/views/${view}.vue`
  }
  
  try {
    return components[componentPath]
  } catch (err) {
    console.error('Failed to load component:', componentPath, err)
    return null
  }
}

// 移除了直接执行的异步代码，改为导出一个函数供在路由守卫中调用
let menus: RouteRecordRaw[] = []

export const initMenus = async () => {
  try {
    //const { data } = await fetchMenus(ruleForm)
    menus = []
    return menus
  } catch (error) {
    console.error('Failed to fetch menus:', error)
    return []
  }
}

export default menus