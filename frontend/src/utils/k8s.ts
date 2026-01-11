// K8s 资源相关工具函数

/**
 * 获取 Pod 状态类型
 */
export const getPodStatusType = (status: string): string => {
  const statusMap: Record<string, string> = {
    'Running': 'success',
    'Pending': 'warning',
    'Failed': 'danger',
    'Succeeded': 'info',
    'Unknown': 'info',
    'CrashLoopBackOff': 'danger',
    'Error': 'danger',
    'Completed': 'success'
  }
  return statusMap[status] || 'info'
}

/**
 * 获取 Deployment 状态
 */
export const getDeploymentStatus = (deployment: any): string => {
  if (!deployment) return 'Unknown'
  
  const { replicas = 0, ready_replicas = 0, available_replicas = 0 } = deployment
  
  if (ready_replicas === replicas && available_replicas === replicas) {
    return 'Ready'
  } else if (ready_replicas < replicas) {
    return 'Updating'
  } else {
    return 'Unavailable'
  }
}

/**
 * 获取 Node 状态类型
 */
export const getNodeStatusType = (status: string): string => {
  return status === 'Ready' ? 'success' : 'danger'
}

/**
 * 获取 Service 类型标签类型
 */
export const getServiceTypeTag = (type: string): string => {
  const typeMap: Record<string, string> = {
    'ClusterIP': 'success',
    'NodePort': 'warning',
    'LoadBalancer': 'primary',
    'ExternalName': 'info'
  }
  return typeMap[type] || ''
}

/**
 * 解析容器状态
 */
export const parseContainerStatus = (containerStatus: any): string => {
  if (!containerStatus) return 'Unknown'
  
  if (containerStatus.state?.running) {
    return 'Running'
  } else if (containerStatus.state?.waiting) {
    return containerStatus.state.waiting.reason || 'Waiting'
  } else if (containerStatus.state?.terminated) {
    return containerStatus.state.terminated.reason || 'Terminated'
  }
  
  return 'Unknown'
}

/**
 * 格式化容器资源限制
 */
export const formatResourceLimits = (limits: any): string => {
  if (!limits) return '-'
  
  const parts: string[] = []
  
  if (limits.cpu) {
    parts.push(`CPU: ${limits.cpu}`)
  }
  
  if (limits.memory) {
    parts.push(`Memory: ${limits.memory}`)
  }
  
  return parts.join(', ') || '-'
}

/**
 * 计算 Pod 运行时长
 */
export const calculatePodAge = (creationTimestamp: string): string => {
  if (!creationTimestamp) return '-'
  
  const now = new Date()
  const created = new Date(creationTimestamp)
  const diff = now.getTime() - created.getTime()
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  
  if (days > 0) return `${days}d${hours}h`
  if (hours > 0) return `${hours}h${minutes}m`
  return `${minutes}m`
}

/**
 * 解析镜像名称和标签
 */
export const parseImageNameTag = (image: string): { name: string; tag: string } => {
  if (!image) return { name: '', tag: '' }
  
  const parts = image.split(':')
  return {
    name: parts[0] || '',
    tag: parts[1] || 'latest'
  }
}

/**
 * 格式化标签选择器
 */
export const formatLabelSelector = (selector: any): string => {
  if (!selector) return '-'
  
  const labels = Object.entries(selector).map(([key, value]) => `${key}=${value}`)
  return labels.join(', ') || '-'
}

/**
 * 检查资源是否健康
 */
export const isResourceHealthy = (resource: any, type: string): boolean => {
  if (!resource) return false
  
  switch (type) {
    case 'pod':
      return resource.status === 'Running'
    case 'deployment':
      return resource.ready_replicas === resource.replicas
    case 'node':
      return resource.status === 'Ready'
    default:
      return true
  }
}

/**
 * 获取资源图标
 */
export const getResourceIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    'pod': 'Box',
    'deployment': 'Setting',
    'service': 'Connection',
    'configmap': 'Document',
    'secret': 'Lock',
    'node': 'Platform',
    'namespace': 'FolderOpened'
  }
  return iconMap[type.toLowerCase()] || 'Document'
}

/**
 * 格式化端口信息
 */
export const formatPortInfo = (port: any): string => {
  if (!port) return '-'
  
  const { port: portNum, target_port, protocol } = port
  return `${portNum}:${target_port}/${protocol || 'TCP'}`
}

/**
 * 判断是否为系统命名空间
 */
export const isSystemNamespace = (namespace: string): boolean => {
  const systemNamespaces = [
    'default',
    'kube-system',
    'kube-public',
    'kube-node-lease'
  ]
  return systemNamespaces.includes(namespace)
}

/**
 * 获取资源配额使用百分比
 */
export const getQuotaPercentage = (used: number, limit: number): number => {
  if (!limit || limit === 0) return 0
  return Math.round((used / limit) * 100)
}

/**
 * 格式化条件信息
 */
export const formatCondition = (condition: any): string => {
  if (!condition) return '-'
  
  const { type, status, reason, message } = condition
  return `${type}: ${status}${reason ? ` (${reason})` : ''}${message ? ` - ${message}` : ''}`
}
