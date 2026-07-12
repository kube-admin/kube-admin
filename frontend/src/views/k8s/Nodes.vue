<template>
  <div class="nodes-container">
    <el-card>
      <template #header>
        <ListToolbar title="Node 列表" :loading="loading" @refresh="fetchNodes">
          <template #filters>
            <el-input v-model="searchKeyword" placeholder="搜索名称" clearable style="width: 160px" />
          </template>
        </ListToolbar>
      </template>

      <el-table :data="filteredNodes" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="250">
          <template #default="scope">
            <router-link
              :to="`/k8s/resource-detail/nodes/${scope.row.name}?group=&version=v1`"
              class="name-link"
            >{{ scope.row.name }}</router-link>
          </template>
        </el-table-column>
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
        <el-table-column label="CPU 使用率" width="180">
          <template #default="scope">
            <el-progress
              :percentage="Math.round(scope.row.usage?.cpu_percent || 0)"
              :color="getProgressColor(scope.row.usage?.cpu_percent || 0)"
              :stroke-width="10"
            />
            <div style="font-size: 12px; color: #909399; margin-top: 2px">
              {{ scope.row.usage?.cpu_used || '-' }} / {{ scope.row.allocatable?.cpu }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="内存使用率" width="180">
          <template #default="scope">
            <el-progress
              :percentage="Math.round(scope.row.usage?.memory_percent || 0)"
              :color="getProgressColor(scope.row.usage?.memory_percent || 0)"
              :stroke-width="10"
            />
            <div style="font-size: 12px; color: #909399; margin-top: 2px">
              {{ formatMemory(scope.row.usage?.memory_used) }} / {{ formatMemory(scope.row.allocatable?.memory) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="creation_timestamp" label="加入时间" width="180" />
        <el-table-column label="操作" fixed="right" width="180">
          <template #default="scope">
            <el-button size="small" @click="yamlDrawer?.open(NODE_GVR, '', scope.row.name)">YAML</el-button>
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

        <el-tab-pane label="资源使用" name="resources">
          <el-row :gutter="20">
            <el-col :span="8">
              <el-card>
                <template #header>CPU 实时使用率</template>
                <el-progress
                  :percentage="getCPUUsagePercent()"
                  :color="getProgressColor(getCPUUsagePercent())"
                />
                <div style="margin-top: 10px; text-align: center">
                  <div>已用: {{ currentNode?.usage?.cpu_used || '-' }}</div>
                  <div>可分配: {{ currentNode?.allocatable?.cpu }}</div>
                  <div>总容量: {{ currentNode?.capacity?.cpu }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <template #header>内存实时使用率</template>
                <el-progress
                  :percentage="getMemoryUsagePercent()"
                  :color="getProgressColor(getMemoryUsagePercent())"
                />
                <div style="margin-top: 10px; text-align: center">
                  <div>已用: {{ formatMemory(currentNode?.usage?.memory_used) }}</div>
                  <div>可分配: {{ formatMemory(currentNode?.allocatable?.memory) }}</div>
                  <div>总容量: {{ formatMemory(currentNode?.capacity?.memory) }}</div>
                </div>
              </el-card>
            </el-col>
            <el-col :span="8">
              <el-card>
                <template #header>Pods</template>
                <div style="margin-top: 10px; text-align: center">
                  <div>总容量: {{ currentNode?.capacity?.pods }}</div>
                  <div>可分配: {{ currentNode?.allocatable?.pods }}</div>
                </div>
              </el-card>
            </el-col>
          </el-row>
          <el-alert
            v-if="!currentNode?.usage?.cpu_used && !currentNode?.usage?.memory_used"
            title="未检测到实时使用率数据，请确认集群已安装 metrics-server"
            type="warning"
            :closable="false"
            style="margin-top: 15px"
          />
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

    <!-- YAML 查看/编辑 -->
    <YamlDrawer ref="yamlDrawer" @saved="fetchNodes" />


  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getNodes, getNodeDetail } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'
import ListToolbar from '@/components/ListToolbar.vue'
import YamlDrawer from '@/components/YamlDrawer.vue'
import type { GVR } from '@/apis/k8s'

const loading = ref(false)
const nodes = ref<any[]>([])
const searchKeyword = ref('')
const filteredNodes = computed(() => nodes.value.filter((r: any) =>
  !searchKeyword.value || (r.name || '').toLowerCase().includes(searchKeyword.value.toLowerCase())
))
const detailDialogVisible = ref(false)
const currentNode = ref<any>(null)
const activeTab = ref('basic')

// 获取命名空间 store
const namespaceStore = useNamespaceStore()
const NODE_GVR: GVR = { version: 'v1', resource: 'nodes' }
const yamlDrawer = ref()



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
  // 支持 Ki/Mi/Gi/Ti 后缀与纯数字（字节），统一换算为 GB
  const match = memory.match(/^(\d+(?:\.\d+)?)(Ki|Mi|Gi|Ti|Pi|Ei)?$/)
  if (!match) return memory
  const num = parseFloat(match[1])
  const unit = match[2] || ''
  const toGB: Record<string, number> = {
    '': num / 1024 / 1024 / 1024,
    Ki: num / 1024 / 1024,
    Mi: num / 1024,
    Gi: num,
    Ti: num * 1024,
    Pi: num * 1024 * 1024,
    Ei: num * 1024 * 1024 * 1024
  }
  return `${toGB[unit].toFixed(2)} GB`
}

const getCPUUsagePercent = () => {
  // 优先使用 metrics-server 实时使用率
  return Math.round(currentNode.value?.usage?.cpu_percent || 0)
}

const getMemoryUsagePercent = () => {
  return Math.round(currentNode.value?.usage?.memory_percent || 0)
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
.name-link {
  color: var(--el-color-primary);
  text-decoration: none;
}
.name-link:hover {
  text-decoration: underline;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}


</style>