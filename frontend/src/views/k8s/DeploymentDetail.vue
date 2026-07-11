<template>
  <ResourceDetailLayout
    v-if="deploy"
    :title="deploy.name"
    kind="Deployment"
    :name="deploy.name"
    :namespace="deploy.namespace"
    :status="deploy.ready_replicas >= deploy.replicas ? 'ready' : 'updating'"
    status-type="deployment"
    :gvr="{ group: 'apps', version: 'v1', resource: 'deployments' }"
    @refresh="fetchDeploy"
  >
    <template #overview>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="副本">{{ deploy.ready_replicas }} / {{ deploy.replicas }}（ready / desired）</el-descriptions-item>
        <el-descriptions-item label="已更新">{{ deploy.updated_replicas }}</el-descriptions-item>
        <el-descriptions-item label="可用">{{ deploy.available_replicas }}</el-descriptions-item>
        <el-descriptions-item label="更新策略">{{ deploy.strategy }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ deploy.namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ deploy.creation_timestamp }}</el-descriptions-item>
      </el-descriptions>

      <h3 class="section-title">关联 Pods</h3>
      <el-table :data="pods" border v-loading="podsLoading">
        <el-table-column label="名称" min-width="260">
          <template #default="{ row }">
            <router-link :to="{ path: '/k8s/pods/' + row.name, query: { namespace: row.namespace } }" class="pod-link">{{ row.name }}</router-link>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'Running' ? 'success' : 'warning'" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="node_name" label="节点" width="140" />
        <el-table-column prop="pod_ip" label="Pod IP" width="150" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" link type="primary" @click="openTerminal(row)">终端</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!podsLoading && pods.length === 0" description="无关联 Pod" />
    </template>
  </ResourceDetailLayout>
  <el-skeleton v-else :rows="8" animated style="padding: 16px" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import ResourceDetailLayout from '@/components/ResourceDetailLayout.vue'
import { getDeploymentDetail, getPods } from '@/apis/k8s'

const route = useRoute()
const router = useRouter()

const name = computed(() => String(route.params.name || ''))
const namespace = computed(() => String(route.query.namespace || 'default'))
const deploy = ref<any>(null)
const pods = ref<any[]>([])
const podsLoading = ref(false)

const fetchDeploy = async () => {
  try {
    const res: any = await getDeploymentDetail(name.value, namespace.value)
    deploy.value = res.data?.data
  } catch (e) {
    ElMessage.error('获取 Deployment 详情失败')
  }
}

// 关联 Pods：Deployment 管理的 Pod 命名为 <deploy>-<rs>-<pod>
const fetchPods = async () => {
  podsLoading.value = true
  try {
    const res: any = await getPods(namespace.value)
    pods.value = (res.data?.data || []).filter((p: any) => p.name.startsWith(name.value + '-'))
  } catch (e) {
    pods.value = []
  } finally {
    podsLoading.value = false
  }
}

const openTerminal = (pod: any) => {
  const container = pod.containers?.[0]?.name || ''
  router.push({ path: `/k8s/pods/${pod.name}/terminal`, query: { namespace: pod.namespace, container } })
}

onMounted(() => {
  fetchDeploy()
  fetchPods()
})
</script>

<style scoped>
.section-title {
  font-size: 14px;
  margin: 16px 0 8px;
}
.pod-link {
  color: var(--el-color-primary);
}
</style>
