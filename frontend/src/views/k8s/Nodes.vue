<template>
  <div class="nodes-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Node 列表</span>
          <el-button type="primary" @click="fetchNodes" :loading="loading">刷新</el-button>
        </div>
      </template>

      <el-table :data="nodes" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="250" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="internal_ip" label="内部IP" width="150" />
        <el-table-column prop="os_image" label="操作系统" width="200" />
        <el-table-column prop="kubelet_version" label="Kubelet版本" width="150" />
        <el-table-column label="CPU" width="120">
          <template #default="scope">
            <div>
              <div style="font-size: 12px; color: #909399">容量: {{ scope.row.capacity?.cpu }}</div>
              <div style="font-size: 12px; color: #67C23A">可用: {{ scope.row.allocatable?.cpu }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="内存" width="150">
          <template #default="scope">
            <div>
              <div style="font-size: 12px; color: #909399">容量: {{ formatMemory(scope.row.capacity?.memory) }}</div>
              <div style="font-size: 12px; color: #67C23A">可用: {{ formatMemory(scope.row.allocatable?.memory) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="creation_timestamp" label="加入时间" width="180" />
        <el-table-column label="操作" fixed="right" width="120">
          <template #default="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 节点详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="Node 详情" width="70%">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="名称">{{ currentNode?.name }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(currentNode?.status)">
                {{ currentNode?.status }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="内部IP">{{ currentNode?.internal_ip }}</el-descriptions-item>
            <el-descriptions-item label="Kubelet版本">{{ currentNode?.kubelet_version }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{ currentNode?.os_image }}</el-descriptions-item>
            <el-descriptions-item label="容器运行时">{{ currentNode?.container_runtime }}</el-descriptions-item>
            <el-descriptions-item label="加入时间" :span="2">{{ currentNode?.creation_timestamp }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>

        <el-tab-pane label="资源容量" name="resources">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card>
                <template #header>CPU</template>
                <el-progress 
                  :percentage="getCPUUsagePercent()" 
                  :color="getProgressColor(getCPUUsagePercent())"
                />
                <div style="margin-top: 10px; text-align: center">
                  <div>总容量: {{ currentNode?.capacity?.cpu }}</div>
                  <div>可分配: {{ currentNode?.allocatable?.cpu }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <template #header>内存</template>
                <el-progress 
                  :percentage="getMemoryUsagePercent()" 
                  :color="getProgressColor(getMemoryUsagePercent())"
                />
                <div style="margin-top: 10px; text-align: center">
                  <div>总容量: {{ formatMemory(currentNode?.capacity?.memory) }}</div>
                  <div>可分配: {{ formatMemory(currentNode?.allocatable?.memory) }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <template #header>Pods</template>
                <el-progress 
                  :percentage="50" 
                  :color="getProgressColor(50)"
                />
                <div style="margin-top: 10px; text-align: center">
                  <div>总容量: {{ currentNode?.capacity?.pods }}</div>
                  <div>可分配: {{ currentNode?.allocatable?.pods }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-tab-pane>

        <el-tab-pane label="节点条件" name="conditions">
          <el-table :data="currentNode?.conditions" border>
            <el-table-column prop="type" label="类型" width="200" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'True' ? 'success' : 'danger'">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="reason" label="原因" width="150" />
            <el-table-column prop="message" label="消息" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="标签" name="labels">
          <el-table :data="getLabelsList(currentNode?.labels)" border>
            <el-table-column prop="key" label="Key" />
            <el-table-column prop="value" label="Value" />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>


  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getNodes, getNodeDetail } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'

const loading = ref(false)
const nodes = ref<any[]>([])
const detailDialogVisible = ref(false)
const currentNode = ref<any>(null)
const activeTab = ref('basic')

// 获取命名空间 store
const namespaceStore = useNamespaceStore()



const fetchNodes = async () => {
  loading.value = true
  try {
    const res = await getNodes()
    nodes.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取Node列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetail = async (node: any) => {
  try {
    const res = await getNodeDetail(node.name)
    currentNode.value = res.data.data
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取Node详情失败')
  }
}



const getStatusType = (status: string) => {
  return status === 'Ready' ? 'success' : 'danger'
}

const formatMemory = (memory: string) => {
  if (!memory) return '-'
  const match = memory.match(/(\d+)Ki/)
  if (match) {
    const kb = parseInt(match[1])
    const gb = (kb / 1024 / 1024).toFixed(2)
    return `${gb} GB`
  }
  return memory
}

const getCPUUsagePercent = () => {
  if (!currentNode.value) return 0
  const capacity = parseFloat(currentNode.value.capacity?.cpu || '0')
  const allocatable = parseFloat(currentNode.value.allocatable?.cpu || '0')
  if (capacity === 0) return 0
  return Math.round((allocatable / capacity) * 100)
}

const getMemoryUsagePercent = () => {
  if (!currentNode.value) return 0
  const capacity = parseMemory(currentNode.value.capacity?.memory || '0')
  const allocatable = parseMemory(currentNode.value.allocatable?.memory || '0')
  if (capacity === 0) return 0
  return Math.round((allocatable / capacity) * 100)
}

const parseMemory = (memory: string) => {
  const match = memory.match(/(\d+)Ki/)
  return match ? parseInt(match[1]) : 0
}

const getProgressColor = (percent: number) => {
  if (percent >= 80) return '#67C23A'
  if (percent >= 50) return '#E6A23C'
  return '#F56C6C'
}

const getLabelsList = (labels: any) => {
  if (!labels) return []
  return Object.entries(labels).map(([key, value]) => ({ key, value }))
}

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取节点列表
  fetchNodes()
})

onMounted(() => {
  fetchNodes()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    // 集群变化时重新获取节点列表
    fetchNodes()
  })
})

onBeforeUnmount(() => {
  // 移除集群变更事件监听
  window.removeEventListener('clusterChanged', () => {
    fetchNodes()
  })
})
</script>

<style scoped>
.nodes-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}


</style>