<template>
  <div class="resources-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ currentType?.label || '资源管理' }}</span>
          <div class="header-actions">
            <el-select v-model="currentTypeKey" size="small" style="width: 200px" @change="onTypeChange">
              <el-option v-for="t in RESOURCE_TYPES" :key="t.key" :label="t.label" :value="t.key" />
            </el-select>
            <el-select
              v-model="namespace"
              size="small"
              style="width: 170px"
              :disabled="!currentType?.namespaced"
              @change="fetchList"
            >
              <el-option label="所有命名空间" value="" />
              <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
            </el-select>
            <el-button size="small" @click="fetchList" :loading="loading">刷新</el-button>
            <el-button size="small" type="primary" @click="openApplyDialog">应用 YAML</el-button>
            <AutoRefresh :interval="30" @refresh="fetchList" />
          </div>
        </div>
      </template>

      <el-table :data="items" v-loading="loading" border style="width: 100%">
        <el-table-column prop="metadata.name" label="名称" min-width="200" show-overflow-tooltip />
        <el-table-column v-if="currentType?.namespaced" prop="metadata.namespace" label="命名空间" width="150" />
        <el-table-column label="创建时间" width="200">
          <template #default="scope">{{ formatTime(scope.row.metadata?.creationTimestamp) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewYaml(scope.row)">YAML</el-button>
            <el-button size="small" type="danger" @click="removeItem(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="!loading && items.length === 0" class="empty-tip">暂无资源</div>
    </el-card>

    <!-- YAML 查看/编辑抽屉 -->
    <el-drawer v-model="yamlDrawerVisible" :title="yamlDrawerTitle" size="60%">
      <YamlEditor v-model="yamlContent" height="70vh" />
      <div class="drawer-footer">
        <el-button @click="yamlDrawerVisible = false">取消</el-button>
        <el-button type="primary" @click="applyFromDrawer" :loading="applying">应用修改</el-button>
      </div>
    </el-drawer>

    <!-- 应用 YAML 对话框 -->
    <el-dialog v-model="applyDialogVisible" title="应用 YAML（创建或更新任意资源）" width="60%">
      <YamlEditor v-model="applyYamlContent" height="50vh" />
      <div class="drawer-footer">
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="doApply" :loading="applying">应用</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listResources, getResource, deleteResource, applyResource, getNamespaces } from '@/apis/k8s'
import type { GVR } from '@/apis/k8s'
import AutoRefresh from '@/components/AutoRefresh.vue'
import YamlEditor from '@/components/YamlEditor.vue'
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
const items = ref<any[]>([])
const namespaces = ref<string[]>([])
const namespace = ref('')
const loading = ref(false)

const yamlDrawerVisible = ref(false)
const yamlDrawerTitle = ref('YAML')
const yamlContent = ref('')
const applyDialogVisible = ref(false)
const applyYamlContent = ref('')
const applying = ref(false)

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

const fetchNamespaces = async () => {
  try {
    const res: any = await getNamespaces()
    namespaces.value = (res.data?.data || []).map((n: any) => n.name).filter(Boolean)
  } catch (e) {
    // 忽略
  }
}

const onTypeChange = () => fetchList()

const viewYaml = async (row: any) => {
  if (!currentType.value) return
  try {
    const ns = row.metadata?.namespace || ''
    const res: any = await getResource(currentType.value.gvr, ns, row.metadata.name)
    const obj = res.data?.data || {}
    yamlContent.value = yaml.dump(obj, { lineWidth: -1, noRefs: true, quotingType: '"' })
    yamlDrawerTitle.value = `YAML - ${row.metadata.name}`
    yamlDrawerVisible.value = true
  } catch (e: any) {
    ElMessage.error(e?.message || '获取详情失败')
  }
}

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

const applyFromDrawer = async () => {
  applying.value = true
  try {
    await applyResource(yamlContent.value)
    ElMessage.success('应用成功')
    yamlDrawerVisible.value = false
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

onMounted(() => {
  initFromRoute()
  fetchNamespaces()
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
