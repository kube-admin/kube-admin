<template>
  <div class="services-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Service 列表</span>
          <div>
            <el-select v-model="namespaceStore.currentNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px">
              <el-option
                v-for="ns in namespaceStore.namespaces"
                :key="ns.name"
                :label="ns.name"
                :value="ns.name"
              />
            </el-select>
            <el-button type="primary" @click="showCreateDialog">创建 Service</el-button>
          </div>
        </div>
      </template>

      <el-table :data="services" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="200" />
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column prop="type" label="类型" width="150">
          <template #default="scope">
            <el-tag :type="getServiceTypeTag(scope.row.type)">
              {{ scope.row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cluster_ip" label="Cluster IP" width="150" />
        <el-table-column label="端口" width="250">
          <template #default="scope">
            <div v-if="scope.row.ports && scope.row.ports.length > 0">
              <el-tag 
                v-for="(port, idx) in scope.row.ports.slice(0, 2)" 
                :key="idx"
                size="small"
                style="margin-right: 5px"
              >
                {{ port.port }}{{ port.protocol ? '/' + port.protocol : '' }}
              </el-tag>
              <el-tag v-if="scope.row.ports.length > 2" size="small" type="info">
                +{{ scope.row.ports.length - 2 }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="Selector" width="200">
          <template #default="scope">
            <div v-if="scope.row.selector && Object.keys(scope.row.selector).length > 0">
              <el-tag 
                v-for="(value, key) in getFirstSelector(scope.row.selector)" 
                :key="key"
                size="small"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="creation_timestamp" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看详情</el-button>
            <el-popconfirm
              title="确定删除这个Service吗?"
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="Service 详情" width="70%">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="名称">{{ currentService?.name }}</el-descriptions-item>
            <el-descriptions-item label="命名空间">{{ currentService?.namespace }}</el-descriptions-item>
            <el-descriptions-item label="类型">
              <el-tag :type="getServiceTypeTag(currentService?.type)">
                {{ currentService?.type }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="Cluster IP">{{ currentService?.cluster_ip }}</el-descriptions-item>
            <el-descriptions-item label="外部IP" :span="2">
              <span v-if="currentService?.external_ip && currentService.external_ip.length > 0">
                {{ currentService.external_ip.join(', ') }}
              </span>
              <span v-else>-</span>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间" :span="2">
              {{ currentService?.creation_timestamp }}
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>

        <el-tab-pane label="端口配置" name="ports">
          <el-table :data="currentService?.ports" border>
            <el-table-column prop="name" label="名称" width="150" />
            <el-table-column prop="protocol" label="协议" width="100" />
            <el-table-column prop="port" label="端口" width="100" />
            <el-table-column prop="target_port" label="目标端口" width="120" />
            <el-table-column label="NodePort" width="100">
              <template #default="scope">
                {{ scope.row.node_port || '-' }}
              </template>
            </el-table-column>
            <el-table-column label="完整配置">
              <template #default="scope">
                <el-tag size="small">
                  {{ scope.row.port }}:{{ scope.row.target_port }}/{{ scope.row.protocol }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="Selector" name="selector">
          <el-table :data="getSelectorList(currentService?.selector)" border>
            <el-table-column prop="key" label="Key" />
            <el-table-column prop="value" label="Value" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="标签" name="labels">
          <el-table :data="getLabelsList(currentService?.labels)" border>
            <el-table-column prop="key" label="Key" />
            <el-table-column prop="value" label="Value" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="注解" name="annotations">
          <el-table :data="getAnnotationsList(currentService?.annotations)" border>
            <el-table-column prop="key" label="Key" width="300" />
            <el-table-column prop="value" label="Value">
              <template #default="scope">
                <div style="word-break: break-all">{{ scope.row.value }}</div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>

    <!-- 创建 Service 对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建 Service (YAML)" width="70%">
      <el-input
        v-model="serviceYaml"
        type="textarea"
        :rows="20"
        placeholder="请输入 Service 的 YAML 定义..."
        style="font-family: monospace"
      />
      <template #footer>
        <el-button @click="loadServiceTemplate">加载模板</el-button>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getServices, getServiceDetail, deleteService, getNamespaces, createServiceFromYaml } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'

const namespaceStore = useNamespaceStore()

const loading = ref(false)
const services = ref<any[]>([])
const detailDialogVisible = ref(false)
const currentService = ref<any>(null)
const activeTab = ref('basic')
const createDialogVisible = ref(false)
const creating = ref(false)
const serviceYaml = ref('')

// 获取命名空间列表
const fetchNamespaces = async () => {
  try {
    const res = await getNamespaces()
    namespaceStore.setNamespaces(res.data.data || [])
  } catch (error) {
    ElMessage.error('获取命名空间失败')
  }
}

// 获取Service列表
const fetchServices = async () => {
  loading.value = true
  try {
    const res = await getServices(namespaceStore.currentNamespace)
    services.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取Service列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetail = async (service: any) => {
  try {
    const res = await getServiceDetail(service.name, service.namespace)
    currentService.value = res.data.data
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取Service详情失败')
  }
}

const handleDelete = async (service: any) => {
  try {
    await deleteService(service.name, service.namespace)
    ElMessage.success('删除成功')
    fetchServices()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const getServiceTypeTag = (type: string) => {
  const typeMap: Record<string, string> = {
    'ClusterIP': 'success',
    'NodePort': 'warning',
    'LoadBalancer': 'primary',
    'ExternalName': 'info'
  }
  return typeMap[type] || ''
}

const getFirstSelector = (selector: any) => {
  if (!selector) return {}
  const entries = Object.entries(selector)
  return entries.length > 0 ? { [entries[0][0]]: entries[0][1] } : {}
}

const getSelectorList = (selector: any) => {
  if (!selector) return []
  return Object.entries(selector).map(([key, value]) => ({ key, value }))
}

const getLabelsList = (labels: any) => {
  if (!labels) return []
  return Object.entries(labels).map(([key, value]) => ({ key, value }))
}

const getAnnotationsList = (annotations: any) => {
  if (!annotations) return []
  return Object.entries(annotations).map(([key, value]) => ({ key, value }))
}

const showCreateDialog = () => {
  loadServiceTemplate()
  createDialogVisible.value = true
}

const loadServiceTemplate = () => {
  serviceYaml.value = `apiVersion: v1
kind: Service
metadata:
  name: my-service
  namespace: ${namespaceStore.currentNamespace}
  labels:
    app: my-app
spec:
  selector:
    app: my-app
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
  type: ClusterIP`
}

const handleCreate = async () => {
  if (!serviceYaml.value.trim()) {
    ElMessage.warning('请输入 Service YAML 定义')
    return
  }
  
  creating.value = true
  try {
    await createServiceFromYaml(serviceYaml.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    fetchServices()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || error.message || '创建失败')
  } finally {
    creating.value = false
  }
}

// 监听命名空间变化
watch(() => namespaceStore.currentNamespace, () => {
  fetchServices()
})

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间和服务列表
  fetchNamespaces()
  fetchServices()
})

onMounted(() => {
  fetchNamespaces()
  fetchServices()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
    fetchServices()
  })
})
</script>

<style scoped>
.services-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>