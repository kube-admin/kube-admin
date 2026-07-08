import { createI18n } from 'vue-i18n'

// 国际化字典。当前覆盖通用与菜单文案，业务页面文案可渐进迁移至 $t()。
const messages = {
  zh: {
    common: {
      create: '创建',
      edit: '编辑',
      delete: '删除',
      save: '保存',
      cancel: '取消',
      confirm: '确定',
      refresh: '刷新',
      search: '搜索',
      detail: '详情',
      success: '成功',
      failed: '失败',
      loading: '加载中...',
      empty: '暂无数据',
      tip: '提示',
      confirmDelete: '确定删除？此操作不可恢复。'
    },
    menu: {
      home: '首页',
      kubernetes: 'Kubernetes',
      dashboard: '仪表盘',
      clusters: '集群',
      namespaces: '命名空间',
      nodes: '节点',
      pods: 'Pods',
      deployments: 'Deployments',
      services: 'Services',
      configmaps: 'ConfigMaps',
      secrets: 'Secrets',
      events: '事件',
      resourceGroup: '资源管理',
      system: '系统设置',
      users: '用户管理'
    },
    auth: {
      login: '登录',
      logout: '退出登录',
      username: '用户名',
      password: '密码',
      title: 'K8s 管理系统'
    }
  },
  en: {
    common: {
      create: 'Create',
      edit: 'Edit',
      delete: 'Delete',
      save: 'Save',
      cancel: 'Cancel',
      confirm: 'OK',
      refresh: 'Refresh',
      search: 'Search',
      detail: 'Detail',
      success: 'Success',
      failed: 'Failed',
      loading: 'Loading...',
      empty: 'No data',
      tip: 'Tip',
      confirmDelete: 'Confirm delete? This cannot be undone.'
    },
    menu: {
      home: 'Home',
      kubernetes: 'Kubernetes',
      dashboard: 'Dashboard',
      clusters: 'Clusters',
      namespaces: 'Namespaces',
      nodes: 'Nodes',
      pods: 'Pods',
      deployments: 'Deployments',
      services: 'Services',
      configmaps: 'ConfigMaps',
      secrets: 'Secrets',
      events: 'Events',
      resourceGroup: 'Resources',
      system: 'System',
      users: 'Users'
    },
    auth: {
      login: 'Login',
      logout: 'Logout',
      username: 'Username',
      password: 'Password',
      title: 'K8s Admin'
    }
  }
}

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('lang') || 'zh',
  fallbackLocale: 'en',
  messages
})

export default i18n
