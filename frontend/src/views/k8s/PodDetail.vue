<template>
  <ResourceDetailLayout
    v-if="pod"
    :title="pod.name"
    kind="Pod"
    :name="pod.name"
    :namespace="pod.namespace"
    :status="pod.status"
    status-type="pod"
    :gvr="{ version: 'v1', resource: 'pods' }"
    @refresh="fetchPod"
  >
    <template #actions>
      <el-button size="small" @click="openTerminal">终端</el-button>
    </template>

    <template #overview>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="状态">{{ pod.status }}</el-descriptions-item>
        <el-descriptions-item label="Pod IP">{{ pod.pod_ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="节点">{{ pod.node_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ pod.namespace }}</el-descriptions-item>
        <el-descriptions-item label="所属">
          <span v-if="pod.owner">
            {{ pod.owner.kind }} /
            <router-link :to="ownerLink" class="owner-link">{{ pod.owner.name }}</router-link>
          </span>
          <span v-else>-（独立 Pod）</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ pod.creation_timestamp }}</el-descriptions-item>
      </el-descriptions>

      <h3 class="section-title">容器</h3>
      <el-table :data="pod.containers" border>
        <el-table-column prop="name" label="名称" min-width="140" />
        <el-table-column prop="image" label="镜像" min-width="220" show-overflow-tooltip />
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.ready ? 'success' : 'warning'" size="small">{{ row.state || (row.ready ? 'Ready' : 'NotReady') }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="restart_count" label="重启" width="70" />
        <el-table-column label="CPU / 内存（实时）" width="190">
          <template #default="{ row }">{{ row.cpu_usage || '-' }} / {{ row.memory_usage || '-' }}</template>
        </el-table-column>
      </el-table>
    </template>
  </ResourceDetailLayout>
  <el-skeleton v-else :rows="8" animated style="padding: 16px" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import ResourceDetailLayout from '@/components/ResourceDetailLayout.vue'
import { getPodDetail } from '@/apis/k8s'

const route = useRoute()
const router = useRouter()

const podName = computed(() => String(route.params.name || ''))
const namespace = computed(() => String(route.query.namespace || 'default'))
const pod = ref<any>(null)

// 所属工作负载的详情链接
const ownerLink = computed(() => {
  const o = pod.value?.owner
  if (!o) return ''
  const map: Record<string, string> = {
    Deployment: `/k8s/deployments/${o.name}`,
    StatefulSet: `/k8s/statefulsets/${o.name}`,
    DaemonSet: `/k8s/daemonsets/${o.name}`
  }
  return (map[o.kind] || '') + `?namespace=${namespace.value}`
})

const fetchPod = async () => {
  try {
    const res: any = await getPodDetail(podName.value, namespace.value)
    pod.value = res.data?.data
  } catch (e) {
    ElMessage.error('获取 Pod 详情失败')
  }
}

const openTerminal = () => {
  const container = pod.value?.containers?.[0]?.name || ''
  // 新标签打开终端，独立于详情页
  const href = router.resolve({ path: `/k8s/pods/${podName.value}/terminal`, query: { namespace: namespace.value, container } }).href
  window.open(href, '_blank')
}

onMounted(fetchPod)
</script>

<style scoped>
.section-title {
  font-size: 14px;
  margin: 16px 0 8px;
}
.owner-link,
:deep(.name-link) {
  color: var(--el-color-primary);
}
</style>
