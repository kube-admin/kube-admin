<template>
  <div class="clusters-container">
    <el-card>
      <template #header>
        <ListToolbar title="Clusters 列表" :loading="loading" @refresh="fetchClusters">
          <template #filters>
            <el-input v-model="searchKeyword" placeholder="搜索名称" clearable style="width: 160px" />
          </template>
          <el-button v-permission="'admin'" type="primary" @click="showCreateDialog">创建 Clusters</el-button>
        </ListToolbar>
      </template>

      <!-- 集群列表 -->
      <el-table :data="filteredClusters" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="150"></el-table-column>
        <el-table-column prop="description" label="描述"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'">
              {{ scope.row.status === 'active' ? '活跃' : '非活跃' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300">
          <template #default="scope">
            <el-button size="small" @click="switchToCluster(scope.row)">切换</el-button>
            <el-button size="small" @click="testConnectionHandler(scope.row)">测试连接</el-button>
            <el-button size="small" @click="editCluster(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteClusterConfirm(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑集群对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="clusterForm" :rules="rules" ref="clusterFormRef" label-width="120px">
        <el-form-item label="集群名称" prop="name">
          <el-input v-model="clusterForm.name" placeholder="请输入集群名称"></el-input>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="clusterForm.description" type="textarea" placeholder="请输入集群描述"></el-input>
        </el-form-item>
        
        <!-- 连接方式说明 -->
        <el-alert
          title="连接方式说明：您可以选择以下任一方式连接集群：1. 提供 kubeconfig 文件内容（推荐） 2. 提供 kubeconfig 文件路径 3. 提供服务器地址和 Token"
          type="info"
          show-icon
          :closable="false"
          style="margin-bottom: 20px;"
        ></el-alert>
        
        <el-form-item label="Config文件内容" prop="config_content">
          <el-input
            v-model="clusterForm.config_content"
            type="textarea"
            :rows="10"
            :placeholder="configContentPlaceholder"
            style="font-family: monospace"
          ></el-input>
        </el-form-item>

        <el-form-item label="Config文件路径" prop="config_path">
          <el-input v-model="clusterForm.config_path" placeholder="可选: kubeconfig文件路径"></el-input>
        </el-form-item>

        <el-form-item label="服务器地址" prop="server_url">
          <el-input v-model="clusterForm.server_url" placeholder="例如: https://kubernetes.default.svc" :disabled="isConnectionMethodDisabled"></el-input>
        </el-form-item>
        <el-form-item label="Token" prop="token">
          <el-input v-model="clusterForm.token" type="password" :placeholder="tokenPlaceholder" :disabled="isConnectionMethodDisabled"></el-input>
        </el-form-item>
        
        <el-alert
          title="注意：如果提供了Config文件内容，则优先使用内容进行连接；否则使用Config文件路径；如果两者都未提供，则使用服务器地址和Token方式进行连接"
          type="warning"
          show-icon
          :closable="false"
        ></el-alert>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitCluster" :loading="submitting">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 测试连接结果对话框 -->
    <el-dialog v-model="testResultVisible" title="测试连接结果" width="400px">
      <div v-if="testResult.success" class="test-success">
        <el-icon class="success-icon"><SuccessFilled /></el-icon>
        <p>{{ testResult.message }}</p>
        <p v-if="testResult.version">版本: {{ testResult.version }}</p>
      </div>
      <div v-else class="test-failure">
        <el-icon class="failure-icon"><CircleCloseFilled /></el-icon>
        <p>{{ testResult.message }}</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="testResultVisible = false">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { SuccessFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import {
  listClusters,
  getCluster,
  createCluster,
  updateCluster,
  deleteCluster,
  testConnectionById
} from '@/apis/k8s/clusters'
import ListToolbar from '@/components/ListToolbar.vue'

// 数据状态
const clusters = ref<any[]>([])
const searchKeyword = ref('')
const filteredClusters = computed(() => clusters.value.filter((r: any) =>
  !searchKeyword.value || (r.name || '').toLowerCase().includes(searchKeyword.value.toLowerCase())
))
const loading = ref(false)
const submitting = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const dialogTitle = ref('添加集群')
const editingClusterId = ref<number | null>(null)
const editingCluster = ref<any>(null) // 编辑中的集群对象，用于判断已配置状态

// 测试连接结果对话框
const testResultVisible = ref(false)
const testResult = ref({
  success: false,
  message: '',
  version: ''
})

// 表单数据
const clusterForm = reactive({
  name: '',
  description: '',
  server_url: '',
  token: '',
  config_path: '',
  config_content: '' // 新增：配置文件内容
})

// 计算属性：判断连接方式是否被禁用
const isConnectionMethodDisabled = computed(() => {
  return !!clusterForm.config_content
})

// 动态 placeholder：编辑时提示已配置字段留空保持不变
const configContentPlaceholder = computed(() => {
  if (editingClusterId.value && editingCluster.value?.has_config_content) {
    return '已配置，留空表示不修改'
  }
  return '可选: kubeconfig文件内容，如果填写此项则优先使用内容进行连接'
})

const tokenPlaceholder = computed(() => {
  if (editingClusterId.value && editingCluster.value?.has_token) {
    return '已配置，留空表示不修改'
  }
  return '请输入访问Token'
})

// 表单验证规则
const rules = reactive({
  name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
  server_url: [{ required: false, message: '请输入服务器地址', trigger: 'blur' }],
  token: [{ required: false, message: '请输入Token', trigger: 'blur' }]
})

// 表单引用
const clusterFormRef = ref()

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 获取集群列表
const fetchClusters = async () => {
  loading.value = true
  try {
    const response: any = await listClusters()
    // 修复数据处理逻辑，统一使用 response.data.data 格式
    clusters.value = response.data?.data || []
  } catch (error) {
    ElMessage.error('获取集群列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 显示创建对话框
const showCreateDialog = () => {
  dialogTitle.value = '添加集群'
  editingClusterId.value = null
  editingCluster.value = null
  // 重置表单
  clusterForm.name = ''
  clusterForm.description = ''
  clusterForm.server_url = ''
  clusterForm.token = ''
  clusterForm.config_path = ''
  clusterForm.config_content = ''
  dialogVisible.value = true
}

// 编辑集群
const editCluster = (cluster: any) => {
  dialogTitle.value = '编辑集群'
  editingClusterId.value = cluster.id
  editingCluster.value = cluster
  // 填充表单数据（敏感字段不回填，留空表示不修改）
  clusterForm.name = cluster.name
  clusterForm.description = cluster.description
  clusterForm.server_url = cluster.server_url
  clusterForm.token = ''
  clusterForm.config_path = cluster.config_path || ''
  clusterForm.config_content = ''
  dialogVisible.value = true
}

// 提交集群信息
const submitCluster = async () => {
  if (!clusterFormRef.value) return
  
  // 创建时强制要求连接方式；编辑时敏感字段留空表示保留，不强制
  const requireConn = !editingClusterId.value && !clusterForm.config_content && !clusterForm.config_path
  const dynamicRules = {
    name: [{ required: true, message: '请输入集群名称', trigger: 'blur' }],
    server_url: [{ required: requireConn, message: '请输入服务器地址', trigger: 'blur' }],
    token: [{ required: requireConn, message: '请输入Token', trigger: 'blur' }]
  }

  // 更新表单验证规则
  Object.assign(rules, dynamicRules)

  await clusterFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    
    submitting.value = true
    try {
      if (editingClusterId.value) {
        // 更新集群
        await updateCluster(editingClusterId.value, clusterForm)
        ElMessage.success('集群更新成功')
      } else {
        // 创建集群
        await createCluster(clusterForm)
        ElMessage.success('集群创建成功')
      }
      dialogVisible.value = false
      fetchClusters() // 刷新列表
      // 通知 Header 等组件集群列表已变更，刷新集群下拉
      window.dispatchEvent(new CustomEvent('clustersChanged'))
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除集群
const deleteClusterConfirm = (cluster: any) => {
  ElMessageBox.confirm(
    `确定要删除集群 "${cluster.name}" 吗？此操作不可恢复。`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteCluster(cluster.id)
      ElMessage.success('集群删除成功')
      fetchClusters() // 刷新列表
      // 通知 Header 等组件集群列表已变更，刷新集群下拉
      window.dispatchEvent(new CustomEvent('clustersChanged'))
    } catch (error: any) {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }).catch(() => {
    // 用户取消删除
  })
}

// 切换到指定集群
const switchToCluster = (cluster: any) => {
  // 保存当前集群到 localStorage
  localStorage.setItem('currentCluster', JSON.stringify(cluster))
  ElMessage.success(`已切换到集群: ${cluster.name}`)
  
  // 触发全局事件通知其他组件集群已切换
  window.dispatchEvent(new CustomEvent('clusterChanged', { detail: cluster }))
}

// 测试连接（已保存集群基于ID测试，后端使用解密后的凭据）
const testConnectionHandler = async (cluster: any) => {
  try {
    const response: any = await testConnectionById(cluster.id)
    testResult.value = response.data?.data || {
      success: false,
      message: '测试连接失败',
      version: ''
    }
    testResultVisible.value = true
  } catch (error: any) {
    testResult.value = {
      success: false,
      message: error.response?.data?.message || '测试连接失败',
      version: ''
    }
    testResultVisible.value = true
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchClusters()
})

// 监听配置内容变化，动态调整表单验证规则
watch(() => clusterForm.config_content, (newVal) => {
  // 当配置内容存在时，清空服务器地址和Token输入框
  if (newVal) {
    clusterForm.server_url = ''
    clusterForm.token = ''
  }
})
</script>

<style scoped>
.clusters-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}

.test-success {
  text-align: center;
  color: #67c23a;
}

.test-success .success-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.test-failure {
  text-align: center;
  color: #f56c6c;
}

.test-failure .failure-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.dialog-footer {
  text-align: right;
}
</style>