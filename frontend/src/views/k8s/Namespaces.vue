<template>
  <div class="namespaces-container">
    <el-card>
      <template #header>
        <ListToolbar title="Namespace 列表" :loading="loading" @refresh="fetchNamespaces">
          <template #filters>
            <el-input v-model="searchKeyword" placeholder="搜索名称" clearable style="width: 160px" />
          </template>
          <el-button type="primary" @click="showCreateDialog">创建 Namespace</el-button>
        </ListToolbar>
      </template>

      <el-table :data="filteredNamespaces" style="width: 100%" v-loading="loading">
        <el-table-column type="expand">
          <template #default="scope">
            <div style="padding: 8px 16px 16px 48px">
              <el-descriptions :column="2" border size="small" style="margin-bottom: 12px">
                <el-descriptions-item label="状态">
                  <el-tag :type="getStatusType(scope.row.status)" size="small">{{ scope.row.status }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="卡住时长">
                  <span v-if="scope.row.deletionTimestamp">
                    {{ formatStuckDuration(scope.row.deletionTimestamp) }}
                    <span style="color:#909399;margin-left:6px">(since {{ scope.row.deletionTimestamp }})</span>
                  </span>
                  <span v-else style="color:#909399">-</span>
                </el-descriptions-item>
                <el-descriptions-item label="Finalizers" :span="2">
                  <el-tag v-for="f in (scope.row.finalizers || [])" :key="f" size="small" style="margin-right:6px">{{ f }}</el-tag>
                  <span v-if="!scope.row.finalizers || scope.row.finalizers.length === 0" style="color:#909399">无</span>
                </el-descriptions-item>
              </el-descriptions>

              <div style="margin-bottom:8px;font-weight:600">根因 Conditions（status=True）</div>
              <el-table :data="rootCauses(scope.row.conditions)" size="small" border style="margin-bottom:12px">
                <el-table-column prop="type" label="类型" width="280">
                  <template #default="c">
                    <span style="color:#E6A23C;font-weight:600">{{ c.row.type }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="reason" label="Reason" width="220" />
                <el-table-column prop="message" label="Message" />
              </el-table>

              <div style="margin-top:8px">
                <el-button
                  v-if="hasRootCause(scope.row.conditions, 'NamespaceDeletionDiscoveryFailure')"
                  type="warning" size="small"
                  @click="openAPIServiceDrawer()"
                >查看失效 APIService -></el-button>
                <el-button
                  v-if="hasRootCause(scope.row.conditions, 'NamespaceFinalizersRemaining')"
                  type="warning" size="small"
                  @click="handleFinalize(scope.row)"
                >强制清理 finalizers</el-button>
                <el-alert
                  v-if="hasContentRemaining(scope.row.conditions)"
                  type="info" :closable="false" show-icon
                  title="根因为残留资源（ContentRemaining/ContentFailure），界面不提供自动清理；请检查该 namespace 下是否有子资源 finalizer 卡住，手动处理后再删除。"
                  style="margin-top:8px"
                />
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" width="300" />
        <el-table-column prop="status" label="状态" width="150">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
            <el-tooltip
              v-if="scope.row.status === 'Terminating'"
              content="展开查看根因"
              placement="top"
            >
              <el-icon style="margin-left: 4px; color: #E6A23C; vertical-align: middle"><Warning /></el-icon>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="age" label="年龄" width="150" />
        <el-table-column label="资源统计" width="300">
          <template #default="scope">
            <el-tag size="small" style="margin-right: 5px">Pods: {{ scope.row.pods || 0 }}</el-tag>
            <el-tag size="small" type="success" style="margin-right: 5px">Deployments: {{ scope.row.deployments || 0 }}</el-tag>
            <el-tag size="small" type="warning" style="margin-right: 5px">Services: {{ scope.row.services || 0 }}</el-tag>
            <el-tag size="small" type="info">ConfigMaps: {{ scope.row.configmaps || 0 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="320">
          <template #default="scope">
            <el-button size="small" @click="yamlDrawer?.open(NAMESPACE_GVR, '', scope.row.name)">YAML</el-button>
            <el-popconfirm
              title="确定删除这个Namespace吗?这将删除该命名空间下的所有资源!"
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

    <!-- 创建 Namespace 对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建 Namespace" width="400px">
      <el-form :model="namespaceForm" :rules="rules" ref="namespaceFormRef">
        <el-form-item label="名称" prop="name">
          <el-input v-model="namespaceForm.name" placeholder="请输入命名空间名称"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <!-- YAML 查看/编辑 -->
    <YamlDrawer ref="yamlDrawer" @saved="fetchNamespaces" />

    <!-- 失效 APIService 抽屉（诊断 namespace DiscoveryFailed） -->
    <el-drawer v-model="apiserviceDrawerVisible" title="失效 APIService（集群级）" size="680px">
      <el-alert
        type="warning" :closable="false" show-icon
        title="删除 APIService 是集群级危险操作，会影响所有依赖此 API 的客户端。仅当后端服务已不存在或确认无用时执行。"
        style="margin-bottom:12px"
      />
      <el-table :data="unavailableAPIServices" v-loading="apiserviceLoading" size="small" border empty-text="集群当前无失效 APIService，被卡的 namespace 可能是 finalizer 或残留资源根因">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="service" label="后端服务" />
        <el-table-column prop="status" label="状态" width="180" />
        <el-table-column prop="age" label="年龄" width="80" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-popconfirm
              :title="`确定删除 APIService「${scope.row.name}」？集群级操作，不可逆！`"
              confirm-button-text="删除" cancel-button-text="取消"
              @confirm="handleDeleteAPIService(scope.row.name)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Warning } from '@element-plus/icons-vue'
import { getNamespaces, createNamespace, deleteNamespace, finalizeNamespace, getPods, getDeployments, getServices, getConfigMaps, getUnavailableAPIServices, deleteAPIService } from '@/apis/k8s'
import type { GVR } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'
import ListToolbar from '@/components/ListToolbar.vue'
import YamlDrawer from '@/components/YamlDrawer.vue'

const namespaceStore = useNamespaceStore()
const NAMESPACE_GVR: GVR = { version: 'v1', resource: 'namespaces' }
const yamlDrawer = ref()

const loading = ref(false)
const namespaces = ref<any[]>([])
const searchKeyword = ref('')
const filteredNamespaces = computed(() => namespaces.value.filter((r: any) =>
  !searchKeyword.value || (r.name || '').toLowerCase().includes(searchKeyword.value.toLowerCase())
))
const createDialogVisible = ref(false)
const creating = ref(false)

// 表单数据
const namespaceForm = reactive({
  name: ''
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入命名空间名称', trigger: 'blur' }]
}

// 表单引用
const namespaceFormRef = ref()

// 获取状态标签类型
const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    'Active': 'success',
    'Terminating': 'warning'
  }
  return typeMap[status] || 'info'
}

// 根因 conditions：只保留 status==='True' 的
const rootCauses = (conditions: any[]) => (conditions || []).filter((c: any) => c.status === 'True')

const hasRootCause = (conditions: any[], type: string) =>
  (conditions || []).some((c: any) => c.type === type && c.status === 'True')

const hasContentRemaining = (conditions: any[]) =>
  hasRootCause(conditions, 'NamespaceDeletionContentFailure') ||
  hasRootCause(conditions, 'NamespaceContentRemaining')

// 由 deletionTimestamp 计算卡住时长
const formatStuckDuration = (ts: string) => {
  const ms = Date.now() - new Date(ts).getTime()
  if (ms < 60000) return '< 1m'
  const min = Math.floor(ms / 60000)
  if (min < 60) return `${min}m`
  const h = Math.floor(min / 60)
  if (h < 24) return `${h}h ${min % 60}m`
  const d = Math.floor(h / 24)
  return `${d}d ${h % 24}h`
}

// 获取命名空间列表及资源统计
const fetchNamespaces = async () => {
  loading.value = true
  try {
    // 获取基础命名空间信息
    const res = await getNamespaces()
    const namespaceList = res.data.data || []
    
    // 为每个命名空间获取资源统计信息
    const namespacesWithStats = await Promise.all(namespaceList.map(async (ns: any) => {
      try {
        // 并行获取各种资源的数量
        const [podsRes, deploymentsRes, servicesRes, configMapsRes] = await Promise.all([
          getPods(ns.name).catch(() => ({ data: { data: [] } })),
          getDeployments(ns.name).catch(() => ({ data: { data: [] } })),
          getServices(ns.name).catch(() => ({ data: { data: [] } })),
          getConfigMaps(ns.name).catch(() => ({ data: { data: [] } }))
        ])
        
        return {
          ...ns,
          pods: podsRes.data.data?.length || 0,
          deployments: deploymentsRes.data.data?.length || 0,
          services: servicesRes.data.data?.length || 0,
          configmaps: configMapsRes.data.data?.length || 0
        }
      } catch (error) {
        // 如果获取某个命名空间的资源失败，返回基础信息
        return {
          ...ns,
          pods: 0,
          deployments: 0,
          services: 0,
          configmaps: 0
        }
      }
    }))
    
    namespaces.value = namespacesWithStats
    
    // 更新全局命名空间列表
    namespaceStore.setNamespaces(namespacesWithStats)
  } catch (error) {
    ElMessage.error('获取命名空间列表失败')
  } finally {
    loading.value = false
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  namespaceForm.name = ''
  createDialogVisible.value = true
}

// 处理创建
const handleCreate = async () => {
  if (!namespaceFormRef.value) return
  
  await namespaceFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    
    creating.value = true
    try {
      await createNamespace({ name: namespaceForm.name })
      ElMessage.success('创建成功')
      createDialogVisible.value = false
      fetchNamespaces()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '创建失败')
    } finally {
      creating.value = false
    }
  })
}

// 处理删除
const handleDelete = async (namespace: any) => {
  if (namespace.name === 'default' || namespace.name === 'kube-system' || namespace.name === 'kube-public') {
    ElMessage.warning('系统命名空间不能删除')
    return
  }
  
  try {
    await deleteNamespace(namespace.name)
    ElMessage.success('删除成功')
    fetchNamespaces()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '删除失败')
  }
}

// 强制清理 finalizers：用于解除 namespace 卡在 Terminating。
// 危险操作（会立即真正删除 ns 及其下所有资源），用 MessageBox 二次确认 + 输入 ns 名。
const handleFinalize = async (namespace: any) => {
  try {
    await ElMessageBox.confirm(
      `将清空「${namespace.name}」的 finalizers 并强制删除该 namespace 及其下所有资源，操作不可逆！\n这是给卡在 Terminating 的 namespace 的兜底手段，请确认 controller 已不在或 finalizer 确实无法完成。`,
      '危险操作：强制清理 finalizers',
      { type: 'error', confirmButtonText: '确认强制清理', cancelButtonText: '取消' }
    )
  } catch {
    return // 用户取消
  }
  try {
    await finalizeNamespace(namespace.name)
    ElMessage.success('已清理 finalizers，namespace 将被真正删除')
    fetchNamespaces()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '强制清理失败')
  }
}

// 失效 APIService 抽屉
const apiserviceDrawerVisible = ref(false)
const apiserviceLoading = ref(false)
const unavailableAPIServices = ref<any[]>([])

const fetchUnavailableAPIServices = async () => {
  apiserviceLoading.value = true
  try {
    const res = await getUnavailableAPIServices()
    unavailableAPIServices.value = res.data.data || []
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '获取失效 APIService 失败')
  } finally {
    apiserviceLoading.value = false
  }
}

const openAPIServiceDrawer = async () => {
  apiserviceDrawerVisible.value = true
  await fetchUnavailableAPIServices()
}

const handleDeleteAPIService = async (name: string) => {
  try {
    await deleteAPIService(name)
    ElMessage.success('已删除 APIService')
    await fetchUnavailableAPIServices()
    // 删除失效 APIService 后，被卡的 namespace 会继续完成删除
    fetchNamespaces()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '删除 APIService 失败')
  }
}

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间列表
  fetchNamespaces()
})

onMounted(() => {
  fetchNamespaces()

  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
  })
})
</script>

<style scoped>
.namespaces-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}
</style>