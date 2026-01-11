// 全局工具函数

/**
 * 格式化时间戳
 */
export const formatTimestamp = (timestamp: string): string => {
  if (!timestamp) return '-'
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

/**
 * 格式化内存大小
 */
export const formatMemory = (memory: string): string => {
  if (!memory) return '-'
  
  const match = memory.match(/(\d+)Ki/)
  if (match) {
    const kb = parseInt(match[1])
    const gb = (kb / 1024 / 1024).toFixed(2)
    return `${gb} GB`
  }
  
  const miMatch = memory.match(/(\d+)Mi/)
  if (miMatch) {
    const mi = parseInt(miMatch[1])
    const gb = (mi / 1024).toFixed(2)
    return `${gb} GB`
  }
  
  const giMatch = memory.match(/(\d+)Gi/)
  if (giMatch) {
    return `${giMatch[1]} GB`
  }
  
  return memory
}

/**
 * 格式化 CPU
 */
export const formatCPU = (cpu: string): string => {
  if (!cpu) return '-'
  
  const match = cpu.match(/(\d+)m/)
  if (match) {
    return `${parseInt(match[1]) / 1000} Core`
  }
  
  return `${cpu} Core`
}

/**
 * 获取相对时间
 */
export const getRelativeTime = (timestamp: string): string => {
  if (!timestamp) return '-'
  
  const now = new Date()
  const past = new Date(timestamp)
  const diff = now.getTime() - past.getTime()
  
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  
  if (days > 0) return `${days}天前`
  if (hours > 0) return `${hours}小时前`
  if (minutes > 0) return `${minutes}分钟前`
  return `${seconds}秒前`
}

/**
 * 复制到剪贴板
 */
export const copyToClipboard = async (text: string): Promise<boolean> => {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (error) {
    // 降级方案
    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    const success = document.execCommand('copy')
    document.body.removeChild(textarea)
    return success
  }
}

/**
 * 下载文件
 */
export const downloadFile = (content: string, filename: string) => {
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

/**
 * 防抖函数
 */
export const debounce = <T extends (...args: any[]) => any>(
  func: T,
  wait: number
): ((...args: Parameters<T>) => void) => {
  let timeout: ReturnType<typeof setTimeout> | null = null
  
  return function(...args: Parameters<T>) {
    if (timeout) clearTimeout(timeout)
    timeout = setTimeout(() => {
      func(...args)
    }, wait)
  }
}

/**
 * 节流函数
 */
export const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle: boolean
  
  return function(...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

/**
 * 解析 YAML
 */
export const parseYAML = (yaml: string): any => {
  // 简单的 YAML 解析 (生产环境应使用 js-yaml 库)
  try {
    return JSON.parse(yaml)
  } catch {
    return null
  }
}

/**
 * 对象转 YAML 格式字符串
 */
export const objectToYAML = (obj: any, indent: number = 0): string => {
  const spaces = ' '.repeat(indent)
  let yaml = ''
  
  for (const [key, value] of Object.entries(obj)) {
    if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      yaml += `${spaces}${key}:\n${objectToYAML(value, indent + 2)}`
    } else if (Array.isArray(value)) {
      yaml += `${spaces}${key}:\n`
      value.forEach(item => {
        if (typeof item === 'object') {
          yaml += `${spaces}- \n${objectToYAML(item, indent + 4)}`
        } else {
          yaml += `${spaces}- ${item}\n`
        }
      })
    } else {
      yaml += `${spaces}${key}: ${value}\n`
    }
  }
  
  return yaml
}

/**
 * 验证命名空间名称
 */
export const validateNamespaceName = (name: string): boolean => {
  const regex = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/
  return regex.test(name)
}

/**
 * 验证资源名称
 */
export const validateResourceName = (name: string): boolean => {
  const regex = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/
  return regex.test(name)
}

/**
 * 高亮搜索关键词
 */
export const highlightKeyword = (text: string, keyword: string): string => {
  if (!keyword) return text
  const regex = new RegExp(`(${keyword})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

/**
 * 生成随机 ID
 */
export const generateId = (): string => {
  return Math.random().toString(36).substring(2, 15) + 
         Math.random().toString(36).substring(2, 15)
}
