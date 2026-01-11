# kube-admin Frontend

## 运行

```bash
# 安装依赖
npm install
# 或
pnpm install

# 开发
npm run dev

# 构建
npm run build

# 预览构建结果
npm run preview
```

## 技术栈

- Vue 3
- TypeScript
- Element Plus
- Vue Router
- Pinia
- Axios
- Vite

## 目录结构

```
frontend/
├── src/
│   ├── apis/           # API 调用
│   │   ├── client/
│   │   │   └── request.ts
│   │   └── k8s.ts
│   ├── assets/         # 静态资源
│   ├── components/     # 公共组件
│   ├── layout/         # 布局组件
│   ├── router/         # 路由配置
│   │   ├── index.ts
│   │   └── menus.ts
│   ├── stores/         # 状态管理
│   ├── views/          # 页面组件
│   │   └── k8s/       # K8s 管理页面
│   │       ├── Pods.vue
│   │       └── Deployments.vue
│   ├── App.vue
│   └── main.ts
├── index.html
├── package.json
└── vite.config.ts
```

## 开发说明

API 代理配置在 `vite.config.ts` 中,开发环境会自动代理 `/api` 请求到后端服务。
