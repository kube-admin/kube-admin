<template>
  <div class="resources-container">
    <el-card>
      <template #header>
        <ListToolbar :title="currentType?.label || '资源管理'" :loading="loading" @refresh="fetchList">
          <template #filters>
            <el-input v-model="searchKeyword" placeholder="搜索名称" clearable style="width: 160px" />
          </template>
          <el-button type="primary" @click="openApplyDialog">应用 YAML</el-button>
        </ListToolbar>
      </template>

      <el-table :data="filteredItems" v-loading="loading" border style="width: 100%">
        <el-table-column prop="metadata.name" label="名称" min-width="200" show-overflow-tooltip>
          <template #default="scope">
            <router-link v-if="detailPath(scope.row)" :to="detailPath(scope.row)" class="name-link">{{ scope.row.metadata.name }}</router-link>
            <span v-else>{{ scope.row.metadata.name }}</span>
          </template>
        </el-table-column>
        <el-table-column v-if="currentType?.namespaced" prop="metadata.namespace" label="命名空间" width="150" />
        <el-table-column label="创建时间" width="200">
          <template #default="scope">{{ formatTime(scope.row.metadata?.creationTimestamp) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="scope">
            <el-button v-if="isWorkload" size="small" @click="showScaleDialog(scope.row)">扩缩容</el-button>
            <el-button v-if="isWorkload" size="small" type="warning" @click="handleRestart(scope.row)">重启</el-button>
            <el-button size="small" @click="yamlDrawer?.open(currentType?.gvr, scope.row.metadata?.namespace || '', scope.row.metadata.name)">YAML</el-button>
            <el-button size="small" type="danger" @click="removeItem(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="!loading && items.length === 0" class="empty-tip">暂无资源</div>
    </el-card>

    <!-- YAML 查看/编辑 -->
    <YamlDrawer ref="yamlDrawer" @saved="fetchList" />

    <!-- 应用 YAML 对话框 -->
    <el-dialog v-model="applyDialogVisible" title="应用 YAML（创建或更新任意资源）" width="60%">
      <YamlEditor v-model="applyYamlContent" height="50vh" />
      <div class="drawer-footer">
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="doApply" :loading="applying">应用</el-button>
      </div>
    </el-dialog>

    <!-- 扩缩容对话框（workload） -->
    <el-dialog v-model="scaleDialogVisible" title="扩缩容" width="400px">
      <el-form>
        <el-form-item label="副本数">
          <el-input-number v-model="scaleReplicas" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="scaleDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleScale">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listResources, getResource, deleteResource, applyResource, scaleResource, restartResource } from '@/apis/k8s'
import type { GVR } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'
import ListToolbar from '@/components/ListToolbar.vue'
import YamlEditor from '@/components/YamlEditor.vue'
import YamlDrawer from '@/components/YamlDrawer.vue'
import yaml from 'js-yaml'

interface ResourceType {
  key: string
  label: string
  gvr: GVR
  namespaced: boolean
}

// 支持的资源类型（覆盖工作负载/网络/存储/自动扩缩容/RBAC）
const RESOURCE_TYPES: ResourceType[] = [
  { key: 'statefulsets', label: 'StatefulSets', gvr: { group: 'apps', version: 'v1', resource: 'statefulsets' }, namespaced: true },
  { key: 'daemonsets', label: 'DaemonSets', gvr: { group: 'apps', version: 'v1', resource: 'daemonsets' }, namespaced: true },
  { key: 'replicasets', label: 'ReplicaSets', gvr: { group: 'apps', version: 'v1', resource: 'replicasets' }, namespaced: true },
  { key: 'jobs', label: 'Jobs', gvr: { group: 'batch', version: 'v1', resource: 'jobs' }, namespaced: true },
  { key: 'cronjobs', label: 'CronJobs', gvr: { group: 'batch', version: 'v1', resource: 'cronjobs' }, namespaced: true },
  { key: 'ingresses', label: 'Ingresses', gvr: { group: 'networking.k8s.io', version: 'v1', resource: 'ingresses' }, namespaced: true },
  { key: 'networkpolicies', label: 'NetworkPolicies', gvr: { group: 'networking.k8s.io', version: 'v1', resource: 'networkpolicies' }, namespaced: true },
  { key: 'persistentvolumes', label: 'PersistentVolumes', gvr: { version: 'v1', resource: 'persistentvolumes' }, namespaced: false },
  { key: 'persistentvolumeclaims', label: 'PersistentVolumeClaims', gvr: { version: 'v1', resource: 'persistentvolumeclaims' }, namespaced: true },
  { key: 'storageclasses', label: 'StorageClasses', gvr: { group: 'storage.k8s.io', version: 'v1', resource: 'storageclasses' }, namespaced: false },
  { key: 'horizontalpodautoscalers', label: 'HorizontalPodAutoscalers', gvr: { group: 'autoscaling', version: 'v2', resource: 'horizontalpodautoscalers' }, namespaced: true },
  { key: 'serviceaccounts', label: 'ServiceAccounts', gvr: { version: 'v1', resource: 'serviceaccounts' }, namespaced: true },
  { key: 'roles', label: 'Roles', gvr: { group: 'rbac.authorization.k8s.io', version: 'v1', resource: 'roles' }, namespaced: true },
  { key: 'rolebindings', label: 'RoleBindings', gvr: { group: 'rbac.authorization.k8s.io', version: 'v1', resource: 'rolebindings' }, namespaced: true },
  { key: 'clusterroles', label: 'ClusterRoles', gvr: { group: 'rbac.authorization.k8s.io', version: 'v1', resource: 'clusterroles' }, namespaced: false },
  { key: 'clusterrolebindings', label: 'ClusterRoleBindings', gvr: { group: 'rbac.authorization.k8s.io', version: 'v1', resource: 'clusterrolebindings' }, namespaced: false }
]

const route = useRoute()
const currentTypeKey = ref('')
const currentType = computed(() => RESOURCE_TYPES.find((t) => t.key === currentTypeKey.value))

// 资源详情页链接（仅 StatefulSet/DaemonSet/Service 首批有详情页，其余保持纯文本）
const detailPath = (row: any): string => {
  const res = currentType.value?.gvr?.resource
  if (!res) return ''
  const ns = row.metadata?.namespace
  const name = row.metadata?.name
  if (res === 'statefulsets') return `/k8s/statefulsets/${name}?namespace=${ns}`
  if (res === 'daemonsets') return `/k8s/daemonsets/${name}?namespace=${ns}`
  if (res === 'services') return `/k8s/services/${name}?namespace=${ns}`
  return ''
}
const items = ref<any[]>([])
const searchKeyword = ref('')
const filteredItems = computed(() => items.value.filter((r: any) =>
  !searchKeyword.value || (r.metadata?.name || '').toLowerCase().includes(searchKeyword.value.toLowerCase())
))
const namespaceStore = useNamespaceStore()
// 命名空间统一用 Header 全局选择器（store），本页不再自带命名空间选择器
const namespace = computed(() => namespaceStore.currentNamespace || '')
const loading = ref(false)

const yamlDrawer = ref()
const applyDialogVisible = ref(false)
const applyYamlContent = ref('')
const applying = ref(false)

// 当前类型是否为支持扩缩容/重启的 workload（StatefulSet/DaemonSet/ReplicaSet）
const WORKLOAD_RESOURCES = ['statefulsets', 'daemonsets', 'replicasets', 'deployments']
const isWorkload = computed(() => WORKLOAD_RESOURCES.includes(currentType.value?.gvr.resource || ''))

// 扩缩容状态
const scaleDialogVisible = ref(false)
const scaleReplicas = ref(1)
const currentResource = ref<any>(null)

const formatTime = (ts?: string) => {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN')
}

// 从路由 meta.gvr 初始化当前资源类型
const initFromRoute = () => {
  const metaGvr = route.meta?.gvr as Partial<GVR> | undefined
  if (metaGvr && metaGvr.resource) {
    const found = RESOURCE_TYPES.find(
      (t) =>
        t.gvr.resource === metaGvr.resource &&
        (t.gvr.group || '') === (metaGvr.group || '') &&
        t.gvr.version === metaGvr.version
    )
    if (found) {
      currentTypeKey.value = found.key
      return
    }
  }
  if (!currentTypeKey.value) currentTypeKey.value = RESOURCE_TYPES[0].key
}

const fetchList = async () => {
  if (!currentType.value) return
  loading.value = true
  try {
    const ns = currentType.value.namespaced ? namespace.value : ''
    const res: any = await listResources(currentType.value.gvr, ns)
    items.value = res.data?.data?.items || []
  } catch (e: any) {
    ElMessage.error(e?.message || '获取列表失败')
    items.value = []
  } finally {
    loading.value = false
  }
}

const onTypeChange = () => fetchList()

const removeItem = (row: any) => {
  ElMessageBox.confirm(`确定删除 ${row.metadata.name}？此操作不可恢复。`, '确认删除', {
    type: 'warning'
  })
    .then(async () => {
      if (!currentType.value) return
      try {
        await deleteResource(currentType.value.gvr, row.metadata?.namespace || '', row.metadata.name)
        ElMessage.success('删除成功')
        fetchList()
      } catch (e: any) {
        ElMessage.error(e?.message || '删除失败')
      }
    })
    .catch(() => {})
}

const openApplyDialog = () => {
  applyYamlContent.value = ''
  applyDialogVisible.value = true
}

const doApply = async () => {
  if (!applyYamlContent.value.trim()) {
    ElMessage.warning('YAML 不能为空')
    return
  }
  applying.value = true
  try {
    await applyResource(applyYamlContent.value)
    ElMessage.success('应用成功')
    applyDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e?.message || '应用失败')
  } finally {
    applying.value = false
  }
}

watch(() => route.meta, () => {
  initFromRoute()
  fetchList()
})

// 扩缩容
const showScaleDialog = (row: any) => {
  currentResource.value = row
  scaleReplicas.value = (row.spec?.replicas ?? 1) as number
  scaleDialogVisible.value = true
}
const handleScale = async () => {
  if (!currentResource.value || !currentType.value) return
  try {
    await scaleResource(currentType.value.gvr, currentResource.value.metadata.namespace || '', currentResource.value.metadata.name, scaleReplicas.value)
    ElMessage.success('扩缩容成功')
    scaleDialogVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e?.message || '扩缩容失败')
  }
}
// 滚动重启
const handleRestart = async (row: any) => {
  if (!currentType.value) return
  try {
    await restartResource(currentType.value.gvr, row.metadata.namespace || '', row.metadata.name)
    ElMessage.success('重启成功')
    fetchList()
  } catch (e: any) {
    ElMessage.error(e?.message || '重启失败')
  }
}

// Header 全局命名空间变化时联动刷新列表
watch(() => namespaceStore.currentNamespace, () => fetchList())

onMounted(() => {
  initFromRoute()
  fetchList()
})
</script>

<style scoped>
.resources-container {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}
.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.empty-tip {
  text-align: center;
  color: #909399;
  padding: 20px;
}
.drawer-footer {
  margin-top: 10px;
  text-align: right;
}
</style>
