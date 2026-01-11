<template>
  <div class="namespaces-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Namespace 列表</span>
          <el-button type="primary" @click="showCreateDialog">创建 Namespace</el-button>
        </div>
      </template>

      <el-table :data="namespaces" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="300" />
        <el-table-column prop="status" label="状态" width="150">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
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
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="scope">
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getNamespaces, createNamespace, deleteNamespace, getPods, getDeployments, getServices, getConfigMaps } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'

const namespaceStore = useNamespaceStore()

const loading = ref(false)
const namespaces = ref<any[]>([])
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
}
</style>