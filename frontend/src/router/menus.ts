const menus = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue'),
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
        component: () => import('@/views/k8s/Dashboard.vue'),
        meta: {
          title: 'Dashboard',
          showMenu: true
        }
      },
      {
        path: '/k8s/clusters',
        name: 'k8sClusters',
        component: () => import('@/views/k8s/Clusters.vue'),
        meta: {
          title: 'Clusters',
          showMenu: true
        }
      },
      {
        path: '/k8s/namespaces',
        name: 'k8sNamespaces',
        component: () => import('@/views/k8s/Namespaces.vue'),
        meta: {
          title: 'Namespaces',
          showMenu: true
        }
      },
      {
        path: '/k8s/nodes',
        name: 'k8sNodes',
        component: () => import('@/views/k8s/Nodes.vue'),
        meta: {
          title: 'Nodes',
          showMenu: true
        }
      },
      {
        path: '/k8s/pods',
        name: 'k8sPods',
        component: () => import('@/views/k8s/Pods.vue'),
        meta: {
          title: 'Pods',
          showMenu: true
        }
      },
      {
        path: '/k8s/deployments',
        name: 'k8sDeployments',
        component: () => import('@/views/k8s/Deployments.vue'),
        meta: {
          title: 'Deployments',
          showMenu: true
        }
      },
      {
        path: '/k8s/services',
        name: 'k8sServices',
        component: () => import('@/views/k8s/Services.vue'),
        meta: {
          title: 'Services',
          showMenu: true
        }
      },
      {
        path: '/k8s/configmaps',
        name: 'k8sConfigMaps',
        component: () => import('@/views/k8s/ConfigMaps.vue'),
        meta: {
          title: 'ConfigMaps',
          showMenu: true
        }
      },
      {
        path: '/k8s/secrets',
        name: 'k8sSecrets',
        component: () => import('@/views/k8s/Secrets.vue'),
        meta: {
          title: 'Secrets',
          showMenu: true
        }
      }
    ]
  },
  {
    path: '/system',
    name: 'system',
    meta: {
      title: '系统设置',
      icon: 'Setting',
      showMenu: true
    },
    children: [
      {
        path: '/system/users',
        name: 'systemUsers',
        component: () => import('@/views/system/Users.vue'),
        meta: {
          title: '用户管理',
          showMenu: true
        }
      }
    ]
  }
]
export default menus