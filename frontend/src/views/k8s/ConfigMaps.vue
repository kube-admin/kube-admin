<template>
  <div class="configmaps-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>ConfigMap 列表</span>
          <div>
            <el-select v-model="namespaceStore.currentNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px">
              <el-option
                v-for="ns in namespaces"
                :key="ns.name"
                :label="ns.name"
                :value="ns.name"
              />
            </el-select>
            <el-button type="primary" @click="showCreateDialog">创建 ConfigMap</el-button>
          </div>
        </div>
      </template>

      <el-table :data="configMaps" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="250" />
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column label="数据项" width="100">
          <template #default="scope">
            <el-tag>{{ Object.keys(scope.row.data || {}).length }} 项</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="creation_timestamp" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="250">
          <template #default="scope">
            <el-button size="small" @click="viewDetail(scope.row)">查看</el-button>
            <el-button size="small" type="warning" @click="editConfigMap(scope.row)">编辑</el-button>
            <el-popconfirm
              title="确定删除这个ConfigMap吗?"
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
        <el-form-item label="数据">
          <div style="width: 100%">
            <div v-for="(item, index) in form.dataList" :key="index" style="margin-bottom: 10px">
              <el-row :gutter="10">
                <el-col :span="8">
                  <el-input v-model="item.key" placeholder="Key" />
                </el-col>
                <el-col :span="14">
                  <el-input v-model="item.value" type="textarea" :rows="2" placeholder="Value" />
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
    <el-dialog v-model="detailDialogVisible" title="ConfigMap 详情" width="60%">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="名称">{{ currentConfigMap?.name }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ currentConfigMap?.namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentConfigMap?.creation_timestamp }}</el-descriptions-item>
      </el-descriptions>
      <el-divider />
      <div style="margin-top: 20px">
        <h4>数据内容:</h4>
        <el-table :data="getDataList(currentConfigMap?.data)" border>
          <el-table-column prop="key" label="Key" width="200" />
          <el-table-column prop="value" label="Value">
            <template #default="scope">
              <pre style="margin: 0">{{ scope.row.value }}</pre>
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
  getConfigMaps,
  createConfigMap,
  updateConfigMap,
  deleteConfigMap,
  getNamespaces
} from '@/apis/k8s'
import type { AxiosResponse } from 'axios'
import type { Result } from '@/apis/client/request'
import { useNamespaceStore } from '@/stores/namespace'

const loading = ref(false)
const configMaps = ref<any[]>([])
const namespaces = ref<any[]>([])
const currentNamespace = ref('default')
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const isEdit = ref(false)
const dialogTitle = ref('创建 ConfigMap')
const currentConfigMap = ref<any>(null)

const form = ref({
  name: '',
  namespace: 'default',
  dataList: [{ key: '', value: '' }]
})

const namespaceStore = useNamespaceStore()

const fetchNamespaces = async () => {
  try {
    const res: AxiosResponse<Result<any[]>> = await getNamespaces()
    namespaces.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取命名空间失败')
  }
}

const fetchConfigMaps = async () => {
  loading.value = true
  try {
    const res: AxiosResponse<Result<any[]>> = await getConfigMaps(namespaceStore.currentNamespace)
    configMaps.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取ConfigMap列表失败')
  } finally {
    loading.value = false
  }
}

const showCreateDialog = () => {
  isEdit.value = false
  dialogTitle.value = '创建 ConfigMap'
  form.value = {
    name: '',
    namespace: currentNamespace.value,
    dataList: [{ key: '', value: '' }]
  }
  dialogVisible.value = true
}

const editConfigMap = (configMap: any) => {
  isEdit.value = true
  dialogTitle.value = '编辑 ConfigMap'
  const dataList = Object.entries(configMap.data || {}).map(([key, value]) => ({
    key,
    value: value as string
  }))
  form.value = {
    name: configMap.name,
    namespace: configMap.namespace,
    dataList: dataList.length > 0 ? dataList : [{ key: '', value: '' }]
  }
  dialogVisible.value = true
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
      await updateConfigMap(form.value.name, form.value.namespace, data)
      ElMessage.success('更新成功')
    } else {
      await createConfigMap({
        namespace: form.value.namespace,
        name: form.value.name,
        data
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchConfigMaps()
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  }
}

const viewDetail = (configMap: any) => {
  currentConfigMap.value = configMap
  detailDialogVisible.value = true
}

const getDataList = (data: any) => {
  if (!data) return []
  return Object.entries(data).map(([key, value]) => ({ key, value }))
}

const handleDelete = async (configMap: any) => {
  try {
    await deleteConfigMap(configMap.name, configMap.namespace)
    ElMessage.success('删除成功')
    fetchConfigMaps()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 监听命名空间变化
watch(() => namespaceStore.currentNamespace, () => {
  fetchConfigMaps()
})

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间和ConfigMap列表
  fetchNamespaces()
  fetchConfigMaps()
})

onMounted(() => {
  fetchNamespaces()
  fetchConfigMaps()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
    fetchConfigMaps()
  })
})
</script>

<style scoped>
.configmaps-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
