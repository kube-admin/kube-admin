<template>
  <div class="pods-container">
    <el-card>
      <template #header>
        <ListToolbar title="Pod 列表" :loading="loading" @refresh="fetchPods">
          <template #filters>
            <el-input v-model="searchKeyword" placeholder="搜索名称" clearable style="width: 160px" />
            <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 120px">
              <el-option v-for="s in POD_STATUSES" :key="s" :label="s" :value="s" />
            </el-select>
          </template>
          <el-button type="primary" @click="showCreateDialog">创建 Pod</el-button>
        </ListToolbar>
      </template>

      <el-table :data="filteredPods" style="width: 100%" v-loading="loading">
        <el-table-column label="名称" width="250">
          <template #default="scope">
            <router-link :to="{ path: '/k8s/pods/' + scope.row.name, query: { namespace: scope.row.namespace } }" class="name-link">{{ scope.row.name }}</router-link>
          </template>
        </el-table-column>
        <el-table-column label="所属工作负载" width="240">
          <template #default="scope">
            <router-link v-if="ownerLink(scope.row)" :to="ownerLink(scope.row)" class="name-link">
              {{ scope.row.owner.kind }} / {{ scope.row.owner.name }}
            </router-link>
            <span v-else-if="scope.row.owner">{{ scope.row.owner.kind }} / {{ scope.row.owner.name }}</span>
            <span v-else style="color: #909399">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pod_ip" label="Pod IP" width="150" />
        <el-table-column prop="node_name" label="节点" width="180" />
        <el-table-column label="年龄" width="110">
          <template #default="scope">
            <span :title="scope.row.creation_timestamp">{{ relTime(scope.row.creation_timestamp) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="CPU" width="90">
          <template #default="scope">{{ fmtCpu(scope.row.cpu_usage) }}</template>
        </el-table-column>
        <el-table-column label="内存" width="100">
          <template #default="scope">{{ fmtMem(scope.row.memory_usage) }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="300">
          <template #default="scope">
            <el-button size="small" @click="yamlDrawer?.open(POD_GVR, scope.row.namespace, scope.row.name)">YAML</el-button>
            <el-button size="small" @click="showLogs(scope.row)">日志</el-button>
            <el-button size="small" @click="openTerminal(scope.row)">终端</el-button>
            <el-popconfirm
              title="确定删除这个Pod吗?"
              @confirm="handleDelete(scope.row)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 日志对话框（WebSocket 实时流） -->
    <el-dialog
      v-model="logDialogVisible"
      title="Pod 日志（实时流）"
      width="80%"
      :destroy-on-close="true"
      @close="closeLogStream"
    >
      <div class="log-toolbar">
        <el-select v-model="selectedContainer" placeholder="容器" size="small" style="width: 180px" @change="restartLogStream">
          <el-option v-for="container in currentPodContainers" :key="container.name" :label="container.name" :value="container.name" />
        </el-select>
        <el-switch v-model="logFollow" active-text="实时跟踪" @change="restartLogStream" />
        <el-switch v-model="logPrevious" active-text="上个容器" @change="restartLogStream" />
        <el-input-number v-model="logTailLines" :min="100" :max="100000" :step="500" size="small" style="width: 150px" />
        <el-input v-model="logSearch" placeholder="搜索过滤" size="small" clearable style="width: 180px" />
        <el-button size="small" @click="logPaused = !logPaused">{{ logPaused ? '继续' : '暂停' }}</el-button>
        <el-button size="small" @click="downloadLogs">下载</el-button>
        <el-button size="small" @click="logs = ''">清空</el-button>
      </div>
      <pre ref="logBoxRef" class="log-box">{{ filteredLogs || '(等待日志...)' }}</pre>
    </el-dialog>

    <!-- 创建 Pod 对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建 Pod (YAML)" width="70%">
      <el-input
        v-model="podYaml"
        type="textarea"
        :rows="20"
        placeholder="请输入 Pod 的 YAML 定义..."
        style="font-family: monospace"
      />
      <template #footer>
        <el-button @click="loadPodTemplate">加载模板</el-button>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- YAML 查看/编辑 -->
    <YamlDrawer ref="yamlDrawer" @saved="fetchPods" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getPods, deletePod, getPodLogs, getPodLogsStreamUrl, getNamespaces, createPodFromYaml } from '@/apis/k8s'
import { useRouter } from 'vue-router'
import { useNamespaceStore } from '@/stores/namespace'
import ListToolbar from '@/components/ListToolbar.vue'
import YamlDrawer from '@/components/YamlDrawer.vue'
import type { GVR } from '@/apis/k8s'
import type { AxiosResponse } from 'axios'
import type { Result } from '@/apis/client/request'

const namespaceStore = useNamespaceStore()
const POD_GVR: GVR = { version: 'v1', resource: 'pods' }
const yamlDrawer = ref()

const loading = ref(false)
const pods = ref<any[]>([])
// 列表搜索/筛选（前端过滤）
const POD_STATUSES = ['Running', 'Pending', 'Failed', 'Succeeded']
const searchKeyword = ref('')
const statusFilter = ref('')
const filteredPods = computed(() => pods.value.filter((p: any) =>
  (!searchKeyword.value || (p.name || '').toLowerCase().includes(searchKeyword.value.toLowerCase())) &&
  (!statusFilter.value || p.status === statusFilter.value)
))
const logDialogVisible = ref(false)
const logs = ref('')
const selectedContainer = ref('')
const currentPodContainers = ref<any[]>([])
const currentPod = ref<any>(null)
const createDialogVisible = ref(false)
const creating = ref(false)
const podYaml = ref('')

// 获取命名空间列表
const fetchNamespaces = async () => {
  try {
    const res: AxiosResponse<Result<any[]>> = await getNamespaces()
    namespaceStore.setNamespaces(res.data.data || [])
  } catch (error) {
    ElMessage.error('获取命名空间失败')
  }
}

// 获取Pod列表
const fetchPods = async () => {
  loading.value = true
  try {
    const res: AxiosResponse<Result<any[]>> = await getPods(namespaceStore.currentNamespace)
    pods.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取Pod列表失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    'Running': 'success',
    'Pending': 'warning',
    'Failed': 'danger',
    'Succeeded': 'info'
  }
  return typeMap[status] || 'info'
}

// 计算 Pod 所属工作负载的详情页链接（仅有详情页的 Deployment/StatefulSet/DaemonSet 可点）
const ownerLink = (pod: any): string => {
  const o = pod.owner
  if (!o?.kind || !o?.name) return ''
  const ns = pod.namespace
  switch (o.kind) {
    case 'Deployment':
      return `/k8s/deployments/${o.name}?namespace=${ns}`
    case 'StatefulSet':
      return `/k8s/statefulsets/${o.name}?namespace=${ns}`
    case 'DaemonSet':
      return `/k8s/daemonsets/${o.name}?namespace=${ns}`
    default:
      return ''
  }
}

// 相对时间（年龄）：创建时间 → "3天前"，hover 显示原始时间
const relTime = (ts?: string) => {
  if (!ts) return '-'
  const t = new Date(ts).getTime()
  if (isNaN(t)) return '-'
  const s = Math.floor((Date.now() - t) / 1000)
  if (s < 60) return `${s}秒前`
  if (s < 3600) return `${Math.floor(s / 60)}分钟前`
  if (s < 86400) return `${Math.floor(s / 3600)}小时前`
  return `${Math.floor(s / 86400)}天前`
}

// CPU/内存 Quantity 格式化（与 Dashboard 一致）
const fmtCpu = (q?: string | number) => {
  if (q === undefined || q === null || q === '') return '-'
  const s = String(q)
  if (s.endsWith('n')) return (parseInt(s) / 1e9).toFixed(3)
  if (s.endsWith('u')) return (parseInt(s) / 1e6).toFixed(3)
  if (s.endsWith('m')) return (parseInt(s) / 1e3).toFixed(2)
  return parseFloat(s).toFixed(2)
}
const UNIT_BYTES: Record<string, number> = { Ki: 1024, Mi: 1024 ** 2, Gi: 1024 ** 3, Ti: 1024 ** 4 }
const fmtMem = (q?: string | number) => {
  if (q === undefined || q === null || q === '') return '-'
  const s = String(q)
  const m = s.match(/^(\d+(?:\.\d+)?)(Ki|Mi|Gi|Ti)?$/)
  if (!m) return s
  const bytes = parseFloat(m[1]) * (m[2] ? UNIT_BYTES[m[2]] : 1)
  return (bytes / 1024 ** 2).toFixed(0) + ' Mi'
}

const handleDelete = async (pod: any) => {
  try {
    await deletePod(pod.name, pod.namespace)
    ElMessage.success('删除成功')
    fetchPods()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// ====== 日志流（WebSocket）======
let logWs: WebSocket | null = null
const logFollow = ref(true)
const logPrevious = ref(false)
const logTailLines = ref(1000)
const logSearch = ref('')
const logPaused = ref(false)
const logBoxRef = ref<HTMLPreElement | null>(null)

// 搜索过滤后的日志
const filteredLogs = computed(() => {
  if (!logSearch.value) return logs.value
  return logs.value
    .split('\n')
    .filter((l: string) => l.includes(logSearch.value))
    .join('\n')
})

const showLogs = (pod: any) => {
  currentPod.value = pod
  currentPodContainers.value = pod.containers || []
  if (currentPodContainers.value.length > 0) {
    selectedContainer.value = currentPodContainers.value[0].name
  }
  logs.value = ''
  logDialogVisible.value = true
  nextTick(() => connectLogStream())
}

// 连接日志流 WebSocket
const connectLogStream = () => {
  if (!currentPod.value || !selectedContainer.value) return
  if (logWs) {
    logWs.close()
    logWs = null
  }
  const url = getPodLogsStreamUrl(currentPod.value.name, currentPod.value.namespace, selectedContainer.value, {
    follow: logFollow.value,
    previous: logPrevious.value,
    tailLines: logTailLines.value
  })
  try {
    logWs = new WebSocket(url)
    logWs.onmessage = (ev) => {
      if (logPaused.value) return
      logs.value += ev.data
      scrollToLogBottom()
    }
    logWs.onerror = () => {
      ElMessage.error('日志流连接错误')
    }
  } catch (e: any) {
    ElMessage.error('连接日志流失败: ' + (e?.message || ''))
  }
}

// 参数变更后重新连接
const restartLogStream = () => {
  logs.value = ''
  connectLogStream()
}

const closeLogStream = () => {
  if (logWs) {
    logWs.close()
    logWs = null
  }
}

const scrollToLogBottom = () => {
  nextTick(() => {
    if (logBoxRef.value && logFollow.value) {
      logBoxRef.value.scrollTop = logBoxRef.value.scrollHeight
    }
  })
}

const downloadLogs = () => {
  const blob = new Blob([logs.value], { type: 'text/plain;charset=utf-8' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `${currentPod.value?.name}-${selectedContainer.value}.log`
  a.click()
  URL.revokeObjectURL(a.href)
}

// 组件卸载时关闭日志流，避免连接泄漏
onUnmounted(() => {
  closeLogStream()
})

const router = useRouter()

// 打开 Pod 终端（跳转独立全屏页）
const openTerminal = (pod: any) => {
  const container = pod.containers?.[0]?.name || ''
  // 新标签打开终端：独立于列表页，避免误关列表断终端，且支持多开
  const href = router.resolve({
    path: `/k8s/pods/${pod.name}/terminal`,
    query: { namespace: pod.namespace, container }
  }).href
  window.open(href, '_blank')
}

const showCreateDialog = () => {
  loadPodTemplate()
  createDialogVisible.value = true
}

const loadPodTemplate = () => {
  podYaml.value = `apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  namespace: ${namespaceStore.currentNamespace}
  labels:
    app: my-app
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"`
}

const handleCreate = async () => {
  if (!podYaml.value.trim()) {
    ElMessage.warning('请输入 Pod YAML 定义')
    return
  }
  
  creating.value = true
  try {
    await createPodFromYaml(podYaml.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    fetchPods()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || error.message || '创建失败')
  } finally {
    creating.value = false
  }
}

// 监听命名空间变化
watch(() => namespaceStore.currentNamespace, () => {
  fetchPods()
})

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间和Pod列表
  fetchNamespaces()
  fetchPods()
})

onMounted(() => {
  fetchNamespaces()
  fetchPods()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
    fetchPods()
  })
})
</script>

<style scoped>
.pods-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.log-toolbar {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 10px;
}

.log-box {
  background-color: #1e1e1e;
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  padding: 10px;
  border-radius: 4px;
  height: 60vh;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}
</style>