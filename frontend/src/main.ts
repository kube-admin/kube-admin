import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import App from './App.vue'
import router from './router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import locale from 'element-plus/es/locale/lang/zh-cn'
import i18n from './locales'
import { permission } from './directives/permission'

import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import '@/assets/main.scss'

// monaco-editor web worker 配置：
// 未配置时 monaco 会回退到主线程加载 worker 代码，导致编辑大 YAML 时卡 UI（控制台 2 条 warning）。
// YAML 无专用 worker，统一用基础 editor.worker。
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker'
;(self as any).MonacoEnvironment = {
  getWorker() {
    return new EditorWorker()
  }
}

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.use(ElementPlus, { locale: locale })
app.use(createPinia())
app.use(router)
app.use(i18n)

// 自定义指令：按角色控制元素显隐
app.directive('permission', permission)

app.mount('#app')
