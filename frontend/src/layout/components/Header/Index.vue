<template>
  <el-menu
    :default-active="activeIndex"
    class="el-menu-header"
    mode="horizontal"
    :ellipsis="false"
    :router="true"
    @select="handleSelect"
  >
    <div class="el-collapse-icon">
      <a @click="toggleCollapse()">
        <el-icon v-if="isCollapse">
          <Expand />
        </el-icon>
        <el-icon v-else>
          <Fold />
        </el-icon>
      </a>
    </div>
    <div class="breadcrumb">
      <Breadcrumb />
    </div>
    <div class="flex-grow" />
    
    <!-- 集群选择器 -->
    <div class="cluster-selector" v-if="clusters.length > 0">
      <el-select 
        v-model="currentClusterId" 
        placeholder="请选择集群" 
        @change="handleClusterChange"
      >
        <el-option key="" label="请选择集群" value="" />
        <el-option
          v-for="cluster in clusters"
          :key="cluster.id"
          :label="cluster.name"
          :value="cluster.id"
        />
      </el-select>
    </div>
    
    <!-- 命名空间选择器 -->
    <div class="namespace-selector" v-if="clusters.length > 0 && currentClusterId !== ''">
      <el-select 
        v-model="currentNamespace" 
        placeholder="请选择命名空间" 
        @change="handleNamespaceChange"
        :disabled="!currentClusterId"
      >
        <el-option key="" label="请选择命名空间" value="" />
        <el-option
          v-for="namespace in namespaces"
          :key="namespace.name"
          :label="namespace.name"
          :value="namespace.name"
        />
      </el-select>
    </div>
    
    <div class="dark-icon" @click="toggleDark()">
      <el-icon>
        <Moon v-if="isDark" />
        <Sunny v-else />
      </el-icon>
    </div>
    <el-sub-menu index="/">
      <template #title>
        <el-icon>
          <Avatar />
        </el-icon>
        {{ userInfo.username }}
      </template>
      <el-menu-item index="/login" @click="logout">退出</el-menu-item>
    </el-sub-menu>
  </el-menu>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { toggleCollapse, isCollapse } from '@/stores/collapse'
import { toggleDark, isDark } from '@/stores/dark'
import Breadcrumb from '../Breadcrumb/Index.vue'
import { listClusters } from '@/apis/k8s/clusters'
import { getNamespaces } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'

// 从 localStorage 获取用户信息
const userInfo = ref({
  username: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).username : 'Admin'
})

const router = useRouter()
const activeIndex = ref('1')

// 集群相关状态
const clusters = ref<any[]>([])
const currentClusterId = ref<number | ''>('')

// 命名空间相关状态
const namespaces = ref<any[]>([])
const currentNamespace = ref<string>('')

// 获取命名空间 store
const namespaceStore = useNamespaceStore()

const handleSelect = (key: string, keyPath: string[]) => {
  console.log(key, keyPath)
}

// 获取集群列表
const fetchClusters = async () => {
  try {
    const response: any = await listClusters()
    clusters.value = response.data?.data || []
    
    // 如果有集群数据，设置默认选中第一个
    if (clusters.value.length > 0) {
      // 检查是否有之前保存的集群选择
      const savedCluster = localStorage.getItem('currentCluster')
      if (savedCluster) {
        const cluster = JSON.parse(savedCluster)
        // 检查保存的集群是否存在于当前集群列表中
        const existingCluster = clusters.value.find(c => c.id === cluster.id)
        if (existingCluster) {
          currentClusterId.value = cluster.id
          // 更新 store 中的集群 ID
          namespaceStore.setCurrentClusterId(cluster.id)
        } else {
          // 如果保存的集群不存在，清空选择
          currentClusterId.value = ''
          namespaceStore.setCurrentClusterId(null)
        }
      } else {
        // 默认不自动选择第一个集群，让用户手动选择
        currentClusterId.value = ''
        namespaceStore.setCurrentClusterId(null)
      }
    } else {
      // 如果没有集群数据，清空选择
      currentClusterId.value = ''
      namespaceStore.setCurrentClusterId(null)
    }
  } catch (error) {
    console.error('获取集群列表失败:', error)
    // 出错时清空集群选择
    currentClusterId.value = ''
    namespaceStore.setCurrentClusterId(null)
  }
}

// 获取命名空间列表
const fetchNamespaces = async () => {
  // 如果集群ID为空，不尝试获取命名空间
  if (!currentClusterId.value) {
    namespaces.value = []
    currentNamespace.value = ''
    namespaceStore.setNamespaces([])
    namespaceStore.setCurrentNamespace('')
    return
  }
  
  try {
    const response: any = await getNamespaces()
    namespaces.value = response.data?.data || []
    
    // 设置默认命名空间
    const savedNamespace = localStorage.getItem('currentNamespace')
    if (savedNamespace && namespaces.value.some((ns: any) => ns.name === savedNamespace)) {
      currentNamespace.value = savedNamespace
    } else if (namespaces.value.length > 0) {
      // 默认选择第一个命名空间
      currentNamespace.value = namespaces.value[0].name
      localStorage.setItem('currentNamespace', currentNamespace.value)
    } else {
      // 如果没有命名空间数据，清空选择
      currentNamespace.value = ''
      localStorage.removeItem('currentNamespace')
    }
    
    // 更新 store 中的命名空间列表
    namespaceStore.setNamespaces(namespaces.value)
    // 更新 store 中的当前命名空间
    namespaceStore.setCurrentNamespace(currentNamespace.value)
  } catch (error) {
    console.error('获取命名空间列表失败:', error)
    // 出错时清空命名空间选择
    namespaces.value = []
    currentNamespace.value = ''
    namespaceStore.setNamespaces([])
    namespaceStore.setCurrentNamespace('')
  }
}

// 集群切换处理
const handleClusterChange = (clusterId: number | '') => {
  if (clusterId === '') {
    // 如果选择了"请选择集群"，清空命名空间列表和状态
    namespaces.value = []
    currentNamespace.value = ''
    
    // 清空 localStorage 中的集群和命名空间信息
    localStorage.removeItem('currentCluster')
    localStorage.removeItem('currentNamespace')
    
    // 更新 store 中的集群 ID 为空
    namespaceStore.setCurrentClusterId(null)
    namespaceStore.setCurrentNamespace('')
    namespaceStore.setNamespaces([])
    
    // 触发全局事件通知其他组件集群已清空
    window.dispatchEvent(new CustomEvent('clusterChanged', { detail: { id: null, name: '' } }))
  } else {
    const cluster = clusters.value.find(c => c.id === clusterId)
    if (cluster) {
      // 保存当前集群到 localStorage
      localStorage.setItem('currentCluster', JSON.stringify(cluster))
      
      // 更新 store 中的集群 ID
      namespaceStore.setCurrentClusterId(clusterId)
      
      // 触发全局事件通知其他组件集群已切换
      window.dispatchEvent(new CustomEvent('clusterChanged', { detail: cluster }))
      
      // 重新获取命名空间列表
      fetchNamespaces()
    }
  }
}

// 命名空间切换处理
const handleNamespaceChange = (namespace: string) => {
  // 保存当前命名空间到 localStorage
  localStorage.setItem('currentNamespace', namespace)
  
  // 更新 store 中的当前命名空间
  namespaceStore.setCurrentNamespace(namespace)
  
  // 触发全局事件通知其他组件命名空间已切换
  window.dispatchEvent(new CustomEvent('namespaceChanged', { detail: namespace }))
}

// 退出登录
const logout = () => {
  // 清除本地存储的用户信息
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  localStorage.removeItem('currentCluster') // 同时清除集群信息
  localStorage.removeItem('currentNamespace') // 同时清除命名空间信息
  // 跳转到登录页
  router.push('/login')
}

// 组件挂载时获取集群列表和命名空间列表
onMounted(async () => {
  await fetchClusters()
  await fetchNamespaces()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    currentClusterId.value = cluster.id
  })
})

// 监听 store 中的集群 ID 变化
watch(() => namespaceStore.currentClusterId, (newClusterId) => {
  // 处理类型转换，当 newClusterId 为 null 时转换为空字符串
  const formattedClusterId = newClusterId || ''
  if (formattedClusterId !== currentClusterId.value) {
    currentClusterId.value = formattedClusterId as number | ''
  }
})

// 监听 store 中的当前命名空间变化
watch(() => namespaceStore.currentNamespace, (newNamespace) => {
  if (newNamespace !== currentNamespace.value) {
    currentNamespace.value = newNamespace
  }
})

// 监听当前集群 ID 的变化，当它发生变化时重新获取命名空间
watch(currentClusterId, async (newClusterId, oldClusterId) => {
  if (newClusterId !== oldClusterId) {
    await fetchNamespaces()
  }
})
</script>

<style scoped>
.el-menu-header a {
  font-size: 14px;
  font-weight: 500;
}

.el-collapse-icon a {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  font-size: 16px;
  margin: 0 0 0 15px;
}

.breadcrumb {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  font-size: 20px;
  margin: 0 0 0 15px;
}

.dark-icon {
  font-size: 20px;
  margin: 15px 15px 0 15px;
  cursor: pointer;
}

.flex-grow {
  flex-grow: 1;
}

.cluster-selector {
  margin: 0 15px;
  display: flex;
  align-items: center;
}

.namespace-selector {
  margin: 0 15px;
  display: flex;
  align-items: center;
}

.cluster-selector .el-select,
.namespace-selector .el-select {
  width: 150px;
}
</style>