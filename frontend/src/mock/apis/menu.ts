import type { MockMethod } from 'vite-plugin-mock'

const menus = [
  {
    path: '/',
    name: 'home',
    component: 'HomeView',
    meta: {
      title: '首页',
      icon: 'menu',
      showMenu: true
    }
  },
  {
    path: '/k8s',
    name: 'k8s',
    meta: {
      title: 'Kubernetes',
      icon: 'Platform',
      showMenu: true
    },
    children: [
      {
        path: '/k8s/dashboard',
        name: 'k8sDashboard',
        component: 'k8s/Dashboard',
        meta: {
          title: 'Dashboard',
          showMenu: true
        }
      },
      {
        path: '/k8s/clusters',
        name: 'k8sClusters',
        component: 'k8s/Clusters',
        meta: {
          title: '集群管理',
          showMenu: true
        }
      },
      {
        path: '/k8s/namespaces',
        name: 'k8sNamespaces',
        component: 'k8s/Namespaces',
        meta: {
          title: 'Namespaces',
          showMenu: true
        }
      },
      {
        path: '/k8s/nodes',
        name: 'k8sNodes',
        component: 'k8s/Nodes',
        meta: {
          title: 'Nodes',
          showMenu: true
        }
      },
      {
        path: '/k8s/pods',
        name: 'k8sPods',
        component: 'k8s/Pods',
        meta: {
          title: 'Pods',
          showMenu: true
        }
      },
      {
        path: '/k8s/deployments',
        name: 'k8sDeployments',
        component: 'k8s/Deployments',
        meta: {
          title: 'Deployments',
          showMenu: true
        }
      },
      {
        path: '/k8s/services',
        name: 'k8sServices',
        component: 'k8s/Services',
        meta: {
          title: 'Services',
          showMenu: true
        }
      },
      {
        path: '/k8s/configmaps',
        name: 'k8sConfigMaps',
        component: 'k8s/ConfigMaps',
        meta: {
          title: 'ConfigMaps',
          showMenu: true
        }
      },
      {
        path: '/k8s/secrets',
        name: 'k8sSecrets',
        component: 'k8s/Secrets',
        meta: {
          title: 'Secrets',
          showMenu: true
        }
      }
    ]
  }
]

export default [
  {
    url: '/api/system/menus', // 注意，这里只能是string格式
    method: 'get',
    response: (req) => {
      return {
        code: 0,
        data: menus
      }
    }
  }
] as MockMethod[]