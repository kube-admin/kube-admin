import request from '@/apis/client/request'
import { useNamespaceStore } from '@/stores/namespace'

// 获取当前集群ID的辅助函数
let namespaceStoreInstance: ReturnType<typeof useNamespaceStore> | null = null
const getCurrentClusterId = () => {
  if (!namespaceStoreInstance) {
    namespaceStoreInstance = useNamespaceStore()
  }
  return namespaceStoreInstance.currentClusterId
}

// Dashboard APIs
export const getDashboardStats = () => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/dashboard/stats', { params })
}

// Namespace APIs
export const getNamespaces = () => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/namespaces', { params })
}

export const createNamespace = (data: { name: string }) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post('/api/v1/namespaces', data, { params })
}

export const deleteNamespace = (name: string) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.delete(`/api/v1/namespaces/${name}`, { params })
}

// Node APIs
export const getNodes = () => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/nodes', { params })
}

export const getNodeDetail = (name: string) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/nodes/${name}`, { params })
}

// 获取Node SSH WebSocket连接URL
export const getNodeSSHUrl = (name: string) => {
  const clusterId = getCurrentClusterId()
  // 构建完整的API路径
  let apiUrl = `/api/v1/nodes/${name}/ssh`
  if (clusterId) {
    apiUrl += `?cluster_id=${clusterId}`
  }
  
  // 添加认证token
  const token = localStorage.getItem('token')
  if (token) {
    if (clusterId) {
      apiUrl += `&token=${token}`
    } else {
      apiUrl += `?token=${token}`
    }
  }
  
  // 在开发环境中，使用代理路径
  if (process.env.NODE_ENV === 'development') {
    // 开发环境使用相对路径，让Vite代理处理
    return apiUrl
  } else {
    // 生产环境使用完整URL
    const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
    const host = window.location.host
    return `${protocol}${host}${apiUrl}`
  }
}

// Pod APIs
export const getPods = (namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/pods', { params })
}

export const getPodDetail = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/pods/${name}`, { params })
}

export const deletePod = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.delete(`/api/v1/pods/${name}`, { params })
}

export const getPodLogs = (name: string, namespace: string = 'default', container: string = '', tailLines: number = 100) => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace, container, tail_lines: tailLines }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/pods/${name}/logs`, { params })
}

// 在Pod中执行命令
export const execPodCommand = (name: string, namespace: string = 'default', container: string = '', command: string[] = []) => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace, container }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post(`/api/v1/pods/${name}/exec`, { command }, { params })
}

// 获取Pod终端WebSocket连接URL
export const getPodTerminalUrl = (name: string, namespace: string = 'default', container: string = '') => {
  const clusterId = getCurrentClusterId()
  // 构建完整的API路径
  let apiUrl = `/api/v1/pods/${name}/terminal?namespace=${namespace}&container=${container}`
  if (clusterId) {
    apiUrl += `&cluster_id=${clusterId}`
  }
  
  // 添加认证token
  const token = localStorage.getItem('token')
  if (token) {
    apiUrl += `&token=${token}`
  }
  
  // 在开发环境中，使用代理路径
  if (process.env.NODE_ENV === 'development') {
    // 开发环境使用相对路径，让Vite代理处理
    return apiUrl
  } else {
    // 生产环境使用完整URL
    const protocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
    const host = window.location.host
    return `${protocol}${host}${apiUrl}`
  }
}

export const createPodFromYaml = (yaml: string) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post('/api/v1/pods/yaml', { yaml }, { params })
}

// Deployment APIs
export const getDeployments = (namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/deployments', { params })
}

export const getDeploymentDetail = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/deployments/${name}`, { params })
}

export const deleteDeployment = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.delete(`/api/v1/deployments/${name}`, { params })
}

export const scaleDeployment = (name: string, replicas: number, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace, replicas }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.put(`/api/v1/deployments/${name}/scale`, null, { params })
}

export const restartDeployment = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.put(`/api/v1/deployments/${name}/restart`, null, { params })
}

export const createDeploymentFromYaml = (yaml: string) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post('/api/v1/deployments/yaml', { yaml }, { params })
}

// ConfigMap APIs
export const getConfigMaps = (namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/configmaps', { params })
}

export const getConfigMapDetail = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/configmaps/${name}`, { params })
}

export const createConfigMap = (data: { namespace: string; name: string; data: Record<string, string> }) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post('/api/v1/configmaps', data, { params })
}

export const updateConfigMap = (name: string, namespace: string, data: Record<string, string>) => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.put(`/api/v1/configmaps/${name}`, { data }, { params })
}

export const deleteConfigMap = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.delete(`/api/v1/configmaps/${name}`, { params })
}

// Service APIs
export const getServices = (namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/services', { params })
}

export const getServiceDetail = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/services/${name}`, { params })
}

export const deleteService = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.delete(`/api/v1/services/${name}`, { params })
}

export const createServiceFromYaml = (yaml: string) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post('/api/v1/services/yaml', { yaml }, { params })
}

// Secret APIs
export const getSecrets = (namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get('/api/v1/secrets', { params })
}

export const getSecretDetail = (name: string, namespace: string = 'default', decode: boolean = false) => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace, decode }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.get(`/api/v1/secrets/${name}`, { params })
}

export const createSecret = (data: { namespace: string; name: string; type: string; data: Record<string, string> }) => {
  const clusterId = getCurrentClusterId()
  const params: any = {}
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.post('/api/v1/secrets', data, { params })
}

export const updateSecret = (name: string, namespace: string, data: Record<string, string>) => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.put(`/api/v1/secrets/${name}`, { data }, { params })
}

export const deleteSecret = (name: string, namespace: string = 'default') => {
  const clusterId = getCurrentClusterId()
  const params: any = { namespace }
  if (clusterId) {
    params.cluster_id = clusterId
  }
  return request.delete(`/api/v1/secrets/${name}`, { params })
}