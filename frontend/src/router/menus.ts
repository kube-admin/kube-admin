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
      },
      {
        path: '/k8s/events',
        name: 'k8sEvents',
        component: () => import('@/views/k8s/Events.vue'),
        meta: {
          title: 'Events',
          showMenu: true
        }
      },
      {
        path: '/k8s/resource-group',
        name: 'k8sResourceGroup',
        meta: {
          title: '资源管理',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/statefulsets',
            name: 'k8sStatefulSets',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'StatefulSets', showMenu: true, gvr: { group: 'apps', version: 'v1', resource: 'statefulsets' } }
          },
          {
            path: '/k8s/daemonsets',
            name: 'k8sDaemonSets',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'DaemonSets', showMenu: true, gvr: { group: 'apps', version: 'v1', resource: 'daemonsets' } }
          },
          {
            path: '/k8s/replicasets',
            name: 'k8sReplicaSets',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'ReplicaSets', showMenu: true, gvr: { group: 'apps', version: 'v1', resource: 'replicasets' } }
          },
          {
            path: '/k8s/ingresses',
            name: 'k8sIngresses',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'Ingresses', showMenu: true, gvr: { group: 'networking.k8s.io', version: 'v1', resource: 'ingresses' } }
          },
          {
            path: '/k8s/networkpolicies',
            name: 'k8sNetworkPolicies',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'NetworkPolicies', showMenu: true, gvr: { group: 'networking.k8s.io', version: 'v1', resource: 'networkpolicies' } }
          },
          {
            path: '/k8s/pvcs',
            name: 'k8sPVCs',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'PersistentVolumeClaims', showMenu: true, gvr: { version: 'v1', resource: 'persistentvolumeclaims' } }
          },
          {
            path: '/k8s/pvs',
            name: 'k8sPVs',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'PersistentVolumes', showMenu: true, gvr: { version: 'v1', resource: 'persistentvolumes' } }
          },
          {
            path: '/k8s/storageclasses',
            name: 'k8sStorageClasses',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'StorageClasses', showMenu: true, gvr: { group: 'storage.k8s.io', version: 'v1', resource: 'storageclasses' } }
          },
          {
            path: '/k8s/hpas',
            name: 'k8sHPAs',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'HPA 自动扩缩容', showMenu: true, gvr: { group: 'autoscaling', version: 'v2', resource: 'horizontalpodautoscalers' } }
          },
          {
            path: '/k8s/serviceaccounts',
            name: 'k8sServiceAccounts',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'ServiceAccounts', showMenu: true, gvr: { version: 'v1', resource: 'serviceaccounts' } }
          },
          {
            path: '/k8s/roles',
            name: 'k8sRoles',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'Roles', showMenu: true, gvr: { group: 'rbac.authorization.k8s.io', version: 'v1', resource: 'roles' } }
          },
          {
            path: '/k8s/rolebindings',
            name: 'k8sRoleBindings',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'RoleBindings', showMenu: true, gvr: { group: 'rbac.authorization.k8s.io', version: 'v1', resource: 'rolebindings' } }
          }
        ]
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