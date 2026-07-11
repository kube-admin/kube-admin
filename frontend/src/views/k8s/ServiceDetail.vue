<template>
  <ResourceDetailLayout
    v-if="svc"
    :title="svc.name"
    kind="Service"
    :name="svc.name"
    :namespace="svc.namespace"
    :status="(svc.type || '').toLowerCase()"
    status-type="service"
    :gvr="{ version: 'v1', resource: 'services' }"
    @refresh="fetchSvc"
  >
    <template #overview>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="类型">{{ svc.type }}</el-descriptions-item>
        <el-descriptions-item label="Cluster IP">{{ svc.cluster_ip }}</el-descriptions-item>
        <el-descriptions-item label="External IP">{{ (svc.external_ip || []).join(', ') || '-' }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ svc.namespace }}</el-descriptions-item>
        <el-descriptions-item label="Selector" :span="2">{{ formatSelector(svc.selector) }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ svc.creation_timestamp }}</el-descriptions-item>
      </el-descriptions>

      <h3 class="section-title">端口</h3>
      <el-table :data="svc.ports" border>
        <el-table-column prop="name" label="名称" min-width="120" />
        <el-table-column prop="protocol" label="协议" width="80" />
        <el-table-column prop="port" label="端口" width="80" />
        <el-table-column prop="target_port" label="目标端口" width="100" />
        <el-table-column prop="node_port" label="NodePort" width="100" />
      </el-table>
    </template>
  </ResourceDetailLayout>
  <el-skeleton v-else :rows="6" animated style="padding: 16px" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import ResourceDetailLayout from '@/components/ResourceDetailLayout.vue'
import { getServiceDetail } from '@/apis/k8s'

const route = useRoute()
const name = computed(() => String(route.params.name || ''))
const namespace = computed(() => String(route.query.namespace || 'default'))
const svc = ref<any>(null)

const formatSelector = (s: any) => {
  if (!s) return '-'
  return Object.entries(s).map(([k, v]) => `${k}=${v}`).join(', ')
}

const fetchSvc = async () => {
  try {
    const res: any = await getServiceDetail(name.value, namespace.value)
    svc.value = res.data?.data
  } catch (e) {
    ElMessage.error('获取 Service 详情失败')
  }
}

onMounted(fetchSvc)
</script>

<style scoped>
.section-title {
  font-size: 14px;
  margin: 16px 0 8px;
}
</style>
