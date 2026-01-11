import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

export const useNamespaceStore = defineStore('namespace', () => {
  // 当前选中的命名空间
  const currentNamespace = ref<string>(
    localStorage.getItem('currentNamespace') || 'default'
  )

  // 命名空间列表
  const namespaces = ref<any[]>([])

  // 当前集群ID
  const currentClusterId = ref<number | null>(null)

  // 设置当前命名空间
  const setCurrentNamespace = (namespace: string) => {
    currentNamespace.value = namespace
    localStorage.setItem('currentNamespace', namespace)
  }

  // 设置命名空间列表
  const setNamespaces = (list: any[]) => {
    namespaces.value = list
  }

  // 设置当前集群ID
  const setCurrentClusterId = (clusterId: number | null) => {
    currentClusterId.value = clusterId
  }

  // 监听变化，自动保存到 localStorage
  watch(currentNamespace, (newVal) => {
    localStorage.setItem('currentNamespace', newVal)
  })

  return {
    currentNamespace,
    namespaces,
    currentClusterId,
    setCurrentNamespace,
    setNamespaces,
    setCurrentClusterId
  }
})