<template>
  <div class="secrets-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Secret 列表</span>
          <div>
            <el-select v-model="namespaceStore.currentNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px">
              <el-option
                v-for="ns in namespaces"
                :key="ns.name"
                :label="ns.name"
                :value="ns.name"
              />
            </el-select>
            <el-button type="primary" @click="showCreateDialog">创建 Secret</el-button>
          </div>
        </div>
      </template>

      <el-table :data="secrets" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="250" />
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column prop="type" label="类型" width="200" />
        <el-table-column label="数据项" width="100">
          <template #default="scope">
            <el-tag>{{ Object.keys(scope.row.data || {}).length }} 项</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="creation_timestamp" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="250">
          <template #default="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看</el-button>
            <el-button size="small" type="warning" @click="editSecret(scope.row)">编辑</el-button>
            <el-popconfirm
              title="确定删除这个Secret吗?"
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

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="60%">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="命名空间">
          <el-select v-model="form.namespace" :disabled="isEdit">
            <el-option
              v-for="ns in namespaces"
              :key="ns.name"
              :label="ns.name"
              :value="ns.name"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="类型" v-if="!isEdit">
          <el-select v-model="form.type">
            <el-option label="Opaque" value="Opaque" />
            <el-option label="kubernetes.io/tls" value="kubernetes.io/tls" />
            <el-option label="kubernetes.io/dockerconfigjson" value="kubernetes.io/dockerconfigjson" />
          </el-select>
        </el-form-item>
        <el-form-item label="数据">
          <div style="width: 100%">
            <div v-for="(item, index) in form.dataList" :key="index" style="margin-bottom: 10px">
              <el-row :gutter="10">
                <el-col :span="8">
                  <el-input v-model="item.key" placeholder="Key" />
                </el-col>
                <el-col :span="14">
                  <el-input 
                    v-model="item.value" 
                    type="textarea" 
                    :rows="2" 
                    placeholder="Value (明文,系统会自动Base64编码)" 
                    show-password
                  />
                </el-col>
                <el-col :span="2">
                  <el-button type="danger" @click="removeDataItem(index)">删除</el-button>
                </el-col>
              </el-row>
            </div>
            <el-button @click="addDataItem">添加数据项</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 查看详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="Secret 详情" width="60%">
      <el-alert
        title="安全提示"
        type="warning"
        :closable="false"
        style="margin-bottom: 20px"
      >
        Secret 数据已加密显示,点击"显示明文"按钮可查看原始值
      </el-alert>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="名称">{{ currentSecret?.name }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ currentSecret?.namespace }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ currentSecret?.type }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentSecret?.creation_timestamp }}</el-descriptions-item>
      </el-descriptions>
      <el-divider />
      <div style="margin-top: 20px">
        <div style="margin-bottom: 10px">
          <el-button size="small" @click="toggleDecode">
            {{ showDecoded ? '隐藏明文' : '显示明文' }}
          </el-button>
        </div>
        <h4>数据内容:</h4>
        <el-table :data="getDataList(currentSecret?.data)" border>
          <el-table-column prop="key" label="Key" width="200" />
          <el-table-column prop="value" label="Value">
            <template #default="scope">
              <div v-if="showDecoded">
                <pre style="margin: 0">{{ decodeBase64(scope.row.value) }}</pre>
              </div>
              <div v-else>
                <el-tag>{{ scope.row.value.substring(0, 20) }}...</el-tag>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getSecrets,
  getSecretDetail,
  createSecret,
  updateSecret,
  deleteSecret,
  getNamespaces
} from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'

const loading = ref(false)
const secrets = ref<any[]>([])
const namespaces = ref<any[]>([])
const currentNamespace = ref('default')
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const dialogTitle = ref('创建 Secret')
const currentSecret = ref<any>(null)
const showDecoded = ref(false)

// 获取命名空间 store
const namespaceStore = useNamespaceStore()

const form = ref({
  name: '',
  namespace: 'default',
  type: 'Opaque',
  dataList: [{ key: '', value: '' }]
})

const fetchNamespaces = async () => {
  try {
    const res = await getNamespaces()
    namespaces.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取命名空间失败')
  }
}

const fetchSecrets = async () => {
  loading.value = true
  try {
    const res = await getSecrets(namespaceStore.currentNamespace)
    secrets.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取Secret列表失败')
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  dialogTitle.value = '创建 Secret'
  form.value = {
    name: '',
    namespace: currentNamespace.value,
    type: 'Opaque',
    dataList: [{ key: '', value: '' }]
  }
  dialogVisible.value = true
}

const editSecret = async (secret: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑 Secret'
  
  try {
    // 获取解码后的Secret数据
    const res = await getSecretDetail(secret.name, secret.namespace, true)
    const secretData = res.data.data
    
    const dataList = Object.entries(secretData.data || {}).map(([key, value]) => ({
      key,
      value: value as string
    }))
    
    form.value = {
      name: secret.name,
      namespace: secret.namespace,
      type: secret.type,
      dataList: dataList.length > 0 ? dataList : [{ key: '', value: '' }]
    }
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取Secret详情失败')
  }
}

const addDataItem = () => {
  form.value.dataList.push({ key: '', value: '' })
}

const removeDataItem = (index: number) => {
  form.value.dataList.splice(index, 1)
}

const handleSubmit = async () => {
  const data: Record<string, string> = {}
  form.value.dataList.forEach(item => {
    if (item.key) {
      data[item.key] = item.value
    }
  })

  try {
    if (isEdit.value) {
      await updateSecret(form.value.name, form.value.namespace, data)
      ElMessage.success('更新成功')
    } else {
      await createSecret({
        namespace: form.value.namespace,
        name: form.value.name,
        type: form.value.type,
        data
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchSecrets()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

const viewDetail = async (secret: any) => {
  try {
    const res = await getSecretDetail(secret.name, secret.namespace, false)
    currentSecret.value = res.data.data
    showDecoded.value = false
    detailDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取Secret详情失败')
  }
}

const toggleDecode = async () => {
  if (!currentSecret.value) return
  
  try {
    const res = await getSecretDetail(
      currentSecret.value.name,
      currentSecret.value.namespace,
      !showDecoded.value
    )
    currentSecret.value = res.data.data
    showDecoded.value = !showDecoded.value
  } catch (error) {
    ElMessage.error('切换显示模式失败')
  }
}

const getDataList = (data: any) => {
  if (!data) return []
  return Object.entries(data).map(([key, value]) => ({ key, value }))
}

const decodeBase64 = (str: string) => {
  try {
    return atob(str)
  } catch (e) {
    return str
  }
}

const handleDelete = async (secret: any) => {
  try {
    await deleteSecret(secret.name, secret.namespace)
    ElMessage.success('删除成功')
    fetchSecrets()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 监听命名空间变化
watch(() => namespaceStore.currentNamespace, () => {
  fetchSecrets()
})

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间和Secret列表
  fetchNamespaces()
  fetchSecrets()
})

onMounted(() => {
  fetchNamespaces()
  fetchSecrets()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
    fetchSecrets()
  })
})
</script>

<style scoped>
.secrets-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
