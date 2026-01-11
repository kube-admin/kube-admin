<template>
  <div class="deployments-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Deployment 列表</span>
          <div>
            <el-select v-model="namespaceStore.currentNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px">
              <el-option
                v-for="ns in namespaceStore.namespaces"
                :key="ns.name"
                :label="ns.name"
                :value="ns.name"
              />
            </el-select>
            <el-button type="primary" @click="showCreateDialog">创建 Deployment</el-button>
          </div>
        </div>
      </template>

      <el-table :data="deployments" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="250" />
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column label="副本数" width="180">
          <template #default="scope">
            {{ scope.row.ready_replicas }}/{{ scope.row.replicas }}
          </template>
        </el-table-column>
        <el-table-column prop="strategy" label="更新策略" width="150" />
        <el-table-column prop="creation_timestamp" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="300">
          <template #default="scope">
            <el-button size="small" @click="showScaleDialog(scope.row)">扩缩容</el-button>
            <el-button size="small" type="warning" @click="handleRestart(scope.row)">
              重启
            </el-button>
            <el-popconfirm
              title="确定删除这个Deployment吗?"
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

    <!-- 扩缩容对话框 -->
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

    <!-- 创建 Deployment 对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建 Deployment (YAML)" width="70%">
      <el-input
        v-model="deploymentYaml"
        type="textarea"
        :rows="20"
        placeholder="请输入 Deployment 的 YAML 定义..."
        style="font-family: monospace"
      />
      <template #footer>
        <el-button @click="loadDeploymentTemplate">加载模板</el-button>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getDeployments,
  deleteDeployment,
  scaleDeployment,
  restartDeployment,
  getNamespaces
} from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'

const namespaceStore = useNamespaceStore()

const loading = ref(false)
const deployments = ref<any[]>([])
const scaleDialogVisible = ref(false)
const scaleReplicas = ref(1)
const currentDeployment = ref<any>(null)
const createDialogVisible = ref(false)
const creating = ref(false)
const deploymentYaml = ref('')

// 获取命名空间列表
const fetchNamespaces = async () => {
  try {
    const res = await getNamespaces()
    namespaceStore.setNamespaces(res.data.data || [])
  } catch (error) {
    ElMessage.error('获取命名空间失败')
  }
}

// 获取Deployment列表
const fetchDeployments = async () => {
  loading.value = true
  try {
    const res = await getDeployments(namespaceStore.currentNamespace)
    deployments.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取Deployment列表失败')
  } finally {
    loading.value = false
  }
}

const handleDelete = async (deployment: any) => {
  try {
    await deleteDeployment(deployment.name, deployment.namespace)
    ElMessage.success('删除成功')
    fetchDeployments()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const showScaleDialog = (deployment: any) => {
  currentDeployment.value = deployment
  scaleReplicas.value = deployment.replicas
  scaleDialogVisible.value = true
}

const handleScale = async () => {
  if (!currentDeployment.value) return
  try {
    await scaleDeployment(
      currentDeployment.value.name,
      scaleReplicas.value,
      currentDeployment.value.namespace
    )
    ElMessage.success('扩缩容成功')
    scaleDialogVisible.value = false
    fetchDeployments()
  } catch (error) {
    ElMessage.error('扩缩容失败')
  }
}

const handleRestart = async (deployment: any) => {
  try {
    await restartDeployment(deployment.name, deployment.namespace)
    ElMessage.success('重启成功')
    fetchDeployments()
  } catch (error) {
    ElMessage.error('重启失败')
  }
}

const showCreateDialog = () => {
  loadDeploymentTemplate()
  createDialogVisible.value = true
}

const loadDeploymentTemplate = () => {
  deploymentYaml.value = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
  namespace: ${namespaceStore.currentNamespace}
  labels:
    app: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
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
  if (!deploymentYaml.value.trim()) {
    ElMessage.warning('请输入 Deployment YAML 定义')
    return
  }
  
  creating.value = true
  try {
    // 注意：这里需要修改 createDeploymentFromYaml 函数的导入
    ElMessage.warning('创建 Deployment 功能待实现')
    // await createDeploymentFromYaml(deploymentYaml.value)
    createDialogVisible.value = false
    fetchDeployments()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || error.message || '创建失败')
  } finally {
    creating.value = false
  }
}

// 监听命名空间变化
watch(() => namespaceStore.currentNamespace, () => {
  fetchDeployments()
})

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间和Deployment列表
  fetchNamespaces()
  fetchDeployments()
})

onMounted(() => {
  fetchNamespaces()
  fetchDeployments()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
    fetchDeployments()
  })
})
</script>

<style scoped>
.deployments-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>