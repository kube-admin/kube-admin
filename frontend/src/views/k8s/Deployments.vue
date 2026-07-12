<template>
  <div class="deployments-container">
    <el-card>
      <template #header>
        <ListToolbar title="Deployment 列表" :loading="loading" @refresh="fetchDeployments">
          <template #filters>
            <el-input v-model="searchKeyword" placeholder="搜索名称" clearable style="width: 160px" />
          </template>
          <el-button type="primary" @click="showCreateDialog">创建 Deployment</el-button>
        </ListToolbar>
      </template>

      <el-table :data="filteredDeployments" style="width: 100%" v-loading="loading">
        <el-table-column label="名称" width="250">
          <template #default="scope">
            <router-link :to="{ path: '/k8s/deployments/' + scope.row.name, query: { namespace: scope.row.namespace } }" class="name-link">{{ scope.row.name }}</router-link>
          </template>
        </el-table-column>
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column label="副本数" width="180">
          <template #default="scope">
            {{ scope.row.ready_replicas }}/{{ scope.row.replicas }}
          </template>
        </el-table-column>
        <el-table-column prop="strategy" label="更新策略" width="150" />
        <el-table-column prop="creation_timestamp" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="360">
          <template #default="scope">
            <el-button size="small" @click="yamlDrawer?.open(DEPLOY_GVR, scope.row.namespace, scope.row.name)">YAML</el-button>
            <el-button size="small" @click="showScaleDialog(scope.row)">扩缩容</el-button>
            <el-button size="small" @click="openTerminal(scope.row)">终端</el-button>
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

    <!-- YAML 查看/编辑 -->
    <YamlDrawer ref="yamlDrawer" @saved="fetchDeployments" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getDeployments,
  deleteDeployment,
  scaleDeployment,
  restartDeployment,
  createDeploymentFromYaml,
  getNamespaces,
  getPods
} from '@/apis/k8s'
import { useRouter } from 'vue-router'
import { useNamespaceStore } from '@/stores/namespace'
import ListToolbar from '@/components/ListToolbar.vue'
import YamlDrawer from '@/components/YamlDrawer.vue'
import type { GVR } from '@/apis/k8s'

const namespaceStore = useNamespaceStore()
const DEPLOY_GVR: GVR = { group: 'apps', version: 'v1', resource: 'deployments' }
const yamlDrawer = ref()
const router = useRouter()

// 打开该 Deployment 下首个 Running Pod 的终端（跳转全屏页）
const openTerminal = async (deploy: any) => {
  try {
    const res: any = await getPods(deploy.namespace)
    const target = (res.data?.data || []).find(
      (p: any) => p.name.startsWith(deploy.name + '-') && p.status === 'Running'
    )
    if (!target) {
      ElMessage.warning('未找到该 Deployment 下 Running 的 Pod')
      return
    }
    const href = router.resolve({
      path: `/k8s/pods/${target.name}/terminal`,
      query: { namespace: deploy.namespace, container: target.containers?.[0]?.name || '' }
    }).href
    window.open(href, '_blank')
  } catch (e) {
    ElMessage.error('获取 Pod 列表失败')
  }
}

const loading = ref(false)
const deployments = ref<any[]>([])
const searchKeyword = ref('')
const filteredDeployments = computed(() => deployments.value.filter((d: any) =>
  !searchKeyword.value || (d.name || '').toLowerCase().includes(searchKeyword.value.toLowerCase())
))
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
    const res = await getDeployments(namespaceStore.effectiveNamespace)
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
    await createDeploymentFromYaml(deploymentYaml.value)
    ElMessage.success('创建成功')
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
  font-size: 16px;
  font-weight: 600;
}
</style>