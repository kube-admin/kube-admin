# 前端开发指南

## 开发环境设置

### 前置要求
- Node.js 16+
- npm 或 pnpm

### 安装依赖
```bash
cd frontend
npm install
# 或
pnpm install
```

### 开发模式
```bash
npm run dev
```
访问: http://localhost:3000

### 生产构建
```bash
npm run build
```

## 项目结构

```
frontend/
├── src/
│   ├── apis/              # API 调用
│   │   ├── client/
│   │   │   └── request.ts # Axios 封装
│   │   └── k8s.ts        # K8s API
│   ├── assets/           # 静态资源
│   ├── components/       # 公共组件
│   ├── layout/          # 布局组件
│   ├── router/          # 路由配置
│   │   ├── index.ts
│   │   └── menus.ts
│   ├── stores/          # Pinia 状态管理
│   ├── utils/           # 工具函数
│   ├── views/           # 页面组件
│   │   └── k8s/        # K8s 资源管理页面
│   ├── App.vue
│   └── main.ts
├── public/              # 公共资源
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 添加新页面

### 1. 创建 Vue 组件

在 `src/views/k8s/` 下创建新组件:

```vue
<template>
  <div class="resource-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>资源列表</span>
          <el-button type="primary" @click="handleCreate">创建</el-button>
        </div>
      </template>
      
      <el-table :data="resources" v-loading="loading">
        <!-- 表格列定义 -->
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const resources = ref<any[]>([])

const fetchResources = async () => {
  loading.value = true
  try {
    // API 调用
  } catch (error) {
    ElMessage.error('获取资源列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchResources()
})
</script>

<style scoped>
.resource-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
```

### 2. 添加 API 调用

在 `src/apis/k8s.ts` 中添加:

```typescript
export const getResources = (namespace: string = 'default') => {
  return request.get('/api/v1/resources', { params: { namespace } })
}

export const createResource = (data: any) => {
  return request.post('/api/v1/resources', data)
}
```

### 3. 注册路由

在 `src/router/menus.ts` 中添加:

```typescript
{
  path: '/k8s/resources',
  name: 'k8sResources',
  component: () => import('@/views/k8s/Resources.vue'),
  meta: {
    title: 'Resources',
    showMenu: true
  }
}
```

## 组件开发规范

### 1. 命名规范
- **组件文件**: PascalCase (如 `ConfigMaps.vue`)
- **变量/函数**: camelCase (如 `fetchPods`)
- **常量**: UPPER_SNAKE_CASE (如 `API_BASE_URL`)

### 2. 组件结构
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 1. 导入
import { ref, onMounted } from 'vue'

// 2. 类型定义
interface Resource {
  name: string
  namespace: string
}

// 3. 响应式数据
const loading = ref(false)
const resources = ref<Resource[]>([])

// 4. 函数定义
const fetchData = async () => {
  // ...
}

// 5. 生命周期钩子
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
/* 组件样式 */
</style>
```

### 3. TypeScript 类型
```typescript
// 推荐: 定义明确的类型
interface PodInfo {
  name: string
  namespace: string
  status: string
}

const pods = ref<PodInfo[]>([])

// 避免: 使用 any
const pods = ref<any[]>([])
```

## API 调用最佳实践

### 1. 错误处理
```typescript
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getResources()
    resources.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取数据失败')
    console.error('Error:', error)
  } finally {
    loading.value = false
  }
}
```

### 2. Loading 状态
```vue
<el-table :data="resources" v-loading="loading">
  <!-- 表格内容 -->
</el-table>
```

### 3. 操作确认
```vue
<el-popconfirm
  title="确定删除吗?"
  @confirm="handleDelete(item)"
>
  <template #reference>
    <el-button type="danger" size="small">删除</el-button>
  </template>
</el-popconfirm>
```

## 样式规范

### 1. 使用 scoped 样式
```vue
<style scoped>
.container {
  padding: 20px;
}
</style>
```

### 2. 使用 Element Plus 变量
```scss
<style scoped lang="scss">
.custom-button {
  color: var(--el-color-primary);
  background-color: var(--el-bg-color);
}
</style>
```

### 3. 响应式设计
```scss
@media (max-width: 768px) {
  .container {
    padding: 10px;
  }
}
```

## 状态管理

使用 Pinia 进行全局状态管理:

```typescript
// stores/k8s.ts
import { defineStore } from 'pinia'

export const useK8sStore = defineStore('k8s', {
  state: () => ({
    currentNamespace: 'default',
    clusters: []
  }),
  
  actions: {
    setNamespace(namespace: string) {
      this.currentNamespace = namespace
    }
  }
})
```

使用:
```typescript
import { useK8sStore } from '@/stores/k8s'

const k8sStore = useK8sStore()
k8sStore.setNamespace('kube-system')
```

## 调试技巧

### 1. Vue DevTools
安装 Vue DevTools 浏览器扩展

### 2. 网络请求调试
在浏览器开发者工具的 Network 标签中查看 API 请求

### 3. 控制台日志
```typescript
console.log('Data:', data)
console.error('Error:', error)
console.table(resources.value)
```

## 性能优化

### 1. 组件懒加载
```typescript
{
  path: '/k8s/pods',
  component: () => import('@/views/k8s/Pods.vue')
}
```

### 2. 计算属性缓存
```typescript
import { computed } from 'vue'

const filteredPods = computed(() => {
  return pods.value.filter(pod => pod.status === 'Running')
})
```

### 3. 防抖
```typescript
import { debounce } from 'lodash-es'

const handleSearch = debounce((query: string) => {
  // 搜索逻辑
}, 300)
```

## 常见问题

### 1. CORS 错误
确保 `vite.config.ts` 中配置了代理:
```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

### 2. 401 未授权
检查 Token 是否正确设置:
```typescript
// src/apis/client/request.ts
const token = localStorage.getItem('token')
if (token) {
  config.headers!.Authorization = `Bearer ${token}`
}
```

### 3. 类型错误
检查 TypeScript 类型定义是否正确

## 测试

### 单元测试 (TODO)
```bash
npm run test:unit
```

### E2E 测试 (TODO)
```bash
npm run test:e2e
```

## 部署

### 构建生产版本
```bash
npm run build
```

### 预览构建结果
```bash
npm run preview
```

### Docker 部署
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 资源链接

- [Vue 3 文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
- [TypeScript 文档](https://www.typescriptlang.org/)
- [Vite 文档](https://vitejs.dev/)
