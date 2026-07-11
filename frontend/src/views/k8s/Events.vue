<template>
  <div class="events-container">
    <el-card>
      <template #header>
        <ListToolbar title="Events 事件" :loading="loading" @refresh="fetchEvents">
          <el-select
            v-model="namespace"
            placeholder="命名空间"
            size="small"
            style="width: 180px"
            @change="fetchEvents"
          >
            <el-option label="所有命名空间" value="" />
            <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
          </el-select>
          <el-select
            v-model="typeFilter"
            placeholder="类型"
            size="small"
            style="width: 130px"
          >
            <el-option label="全部" value="" />
            <el-option label="Normal" value="Normal" />
            <el-option label="Warning" value="Warning" />
          </el-select>
        </ListToolbar>
      </template>

      <el-table :data="filteredEvents" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'Warning' ? 'danger' : 'info'" size="small">
              {{ scope.row.type || 'Normal' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="namespace" label="命名空间" width="140" />
        <el-table-column prop="involved_object" label="关联对象" width="220" show-overflow-tooltip />
        <el-table-column prop="reason" label="原因" width="160" show-overflow-tooltip />
        <el-table-column prop="message" label="消息" show-overflow-tooltip />
        <el-table-column prop="count" label="次数" width="80" align="center" />
        <el-table-column prop="last_timestamp" label="最近时间" width="180" />
        <el-table-column prop="source" label="来源" width="160" show-overflow-tooltip />
      </el-table>

      <div v-if="!loading && filteredEvents.length === 0" class="empty-tip">
        暂无事件
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getEvents, getNamespaces } from '@/apis/k8s'
import ListToolbar from '@/components/ListToolbar.vue'

interface EventItem {
  name: string
  namespace: string
  type: string
  reason: string
  message: string
  involved_object: string
  count: number
  first_timestamp: string
  last_timestamp: string
  source: string
}

const events = ref<EventItem[]>([])
const namespaces = ref<string[]>([])
const namespace = ref('')
const typeFilter = ref('')
const loading = ref(false)

// 类型过滤在前端完成，避免频繁请求后端
const filteredEvents = computed(() => {
  if (!typeFilter.value) return events.value
  return events.value.filter((e) => (e.type || 'Normal') === typeFilter.value)
})

const fetchEvents = async () => {
  loading.value = true
  try {
    const res: any = await getEvents(namespace.value)
    events.value = res.data?.data || []
  } catch (e: any) {
    ElMessage.error(e?.message || '获取事件失败')
  } finally {
    loading.value = false
  }
}

const fetchNamespaces = async () => {
  try {
    const res: any = await getNamespaces()
    const list = res.data?.data || []
    namespaces.value = list.map((n: any) => n.name).filter(Boolean)
  } catch (e) {
    // 命名空间获取失败不阻塞事件查询
  }
}

onMounted(() => {
  fetchNamespaces()
  fetchEvents()
})
</script>

<style scoped>
.events-container {
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
</style>
