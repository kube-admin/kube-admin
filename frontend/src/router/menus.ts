// 菜单信息架构（IA）：
// - 「集群管理」提为平台级一级菜单（集群接入 CRUD，与具体集群内资源解耦）
// - 「Kubernetes」下资源按职能域分组（工作负载/网络/配置/存储/安全/自动伸缩/集群资源），
//   废弃原按使用频率划分的「资源管理」口袋，消除同职能资源被劈成两半的问题
// - 「可观测」独立一级，容纳 Events，为未来 Logs/监控预留
// 注意：所有资源页 path/name/component 保持不变，仅重组分组树，避免外部跳转断裂
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
  // 平台级：集群接入管理（CRUD/测试连接），不属于某个集群内部
  {
    path: '/k8s/clusters',
    name: 'k8sClusters',
    component: () => import('@/views/k8s/Clusters.vue'),
    meta: {
      title: '集群管理',
      icon: 'Connection',
      showMenu: true
    }
  },
  // 当前选中集群内的资源，按职能域分组
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
      // 工作负载
      {
        path: '/k8s/workloads',
        name: 'k8sWorkloads',
        meta: {
          title: '工作负载',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/pods',
            name: 'k8sPods',
            component: () => import('@/views/k8s/Pods.vue'),
            meta: { title: 'Pods', showMenu: true }
          },
          {
            path: '/k8s/deployments',
            name: 'k8sDeployments',
            component: () => import('@/views/k8s/Deployments.vue'),
            meta: { title: 'Deployments', showMenu: true }
          },
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
          }
        ]
      },
      // 网络
      {
        path: '/k8s/network',
        name: 'k8sNetwork',
        meta: {
          title: '网络',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/services',
            name: 'k8sServices',
            component: () => import('@/views/k8s/Services.vue'),
            meta: { title: 'Services', showMenu: true }
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
          }
        ]
      },
      // 配置
      {
        path: '/k8s/config',
        name: 'k8sConfig',
        meta: {
          title: '配置',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/configmaps',
            name: 'k8sConfigMaps',
            component: () => import('@/views/k8s/ConfigMaps.vue'),
            meta: { title: 'ConfigMaps', showMenu: true }
          },
          {
            path: '/k8s/secrets',
            name: 'k8sSecrets',
            component: () => import('@/views/k8s/Secrets.vue'),
            meta: { title: 'Secrets', showMenu: true }
          }
        ]
      },
      // 存储
      {
        path: '/k8s/storage',
        name: 'k8sStorage',
        meta: {
          title: '存储',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/pvs',
            name: 'k8sPVs',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'PersistentVolumes', showMenu: true, gvr: { version: 'v1', resource: 'persistentvolumes' } }
          },
          {
            path: '/k8s/pvcs',
            name: 'k8sPVCs',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'PersistentVolumeClaims', showMenu: true, gvr: { version: 'v1', resource: 'persistentvolumeclaims' } }
          },
          {
            path: '/k8s/storageclasses',
            name: 'k8sStorageClasses',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'StorageClasses', showMenu: true, gvr: { group: 'storage.k8s.io', version: 'v1', resource: 'storageclasses' } }
          }
        ]
      },
      // 安全 (RBAC)
      {
        path: '/k8s/security',
        name: 'k8sSecurity',
        meta: {
          title: '安全',
          showMenu: true
        },
        children: [
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
      },
      // 自动伸缩
      {
        path: '/k8s/autoscaling',
        name: 'k8sAutoscaling',
        meta: {
          title: '自动伸缩',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/hpas',
            name: 'k8sHPAs',
            component: () => import('@/views/k8s/Resources.vue'),
            meta: { title: 'HPA', showMenu: true, gvr: { group: 'autoscaling', version: 'v2', resource: 'horizontalpodautoscalers' } }
          }
        ]
      },
      // 集群级资源
      {
        path: '/k8s/cluster-resources',
        name: 'k8sClusterResources',
        meta: {
          title: '集群资源',
          showMenu: true
        },
        children: [
          {
            path: '/k8s/nodes',
            name: 'k8sNodes',
            component: () => import('@/views/k8s/Nodes.vue'),
            meta: { title: 'Nodes', showMenu: true }
          },
          {
            path: '/k8s/namespaces',
            name: 'k8sNamespaces',
            component: () => import('@/views/k8s/Namespaces.vue'),
            meta: { title: 'Namespaces', showMenu: true }
          }
        ]
      }
    ]
  },
  // 可观测：独立一级，为未来 Logs/监控/Prometheus 预留
  {
    path: '/observability',
    name: 'observability',
    meta: {
      title: '可观测',
      icon: 'Monitor',
      showMenu: true
    },
    children: [
      {
        path: '/k8s/events',
        name: 'k8sEvents',
        component: () => import('@/views/k8s/Events.vue'),
        meta: {
          title: 'Events',
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
  },
  // ===== 资源详情页（隐藏菜单，layout 直接 children，带侧边栏）=====
  {
    path: '/k8s/pods/:name',
    name: 'podDetail',
    component: () => import('@/views/k8s/PodDetail.vue'),
    meta: { showMenu: false }
  },
  {
    path: '/k8s/deployments/:name',
    name: 'deploymentDetail',
    component: () => import('@/views/k8s/DeploymentDetail.vue'),
    meta: { showMenu: false }
  },
  {
    path: '/k8s/statefulsets/:name',
    name: 'statefulsetDetail',
    component: () => import('@/views/k8s/WorkloadDetail.vue'),
    meta: { showMenu: false }
  },
  {
    path: '/k8s/daemonsets/:name',
    name: 'daemonsetDetail',
    component: () => import('@/views/k8s/WorkloadDetail.vue'),
    meta: { showMenu: false }
  },
  {
    path: '/k8s/replicasets/:name',
    name: 'replicasetDetail',
    component: () => import('@/views/k8s/WorkloadDetail.vue'),
    meta: { showMenu: false }
  },
  {
    path: '/k8s/services/:name',
    name: 'serviceDetail',
    component: () => import('@/views/k8s/ServiceDetail.vue'),
    meta: { showMenu: false }
  }
]
export default menus
