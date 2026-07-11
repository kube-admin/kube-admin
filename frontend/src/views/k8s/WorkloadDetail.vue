<template>
  <ResourceDetailLayout
    v-if="obj"
    :title="name"
    :kind="kind"
    :name="name"
    :namespace="namespace"
    :gvr="gvr"
    @refresh="fetchObj"
  >
    <template #overview>
      <el-descriptions :column="3" border>
        <el-descriptions-item label="副本">{{ obj.status?.readyReplicas ?? '-' }} / {{ obj.spec?.replicas ?? '-' }}（ready / desired）</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ obj.metadata?.creationTimestamp }}</el-descriptions-item>
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
      </el-table>
      <el-empty v-if="!podsLoading && pods.length === 0" description="无关联 Pod" />
    </template>
  </ResourceDetailLayout>
  <el-skeleton v-else :rows="8" animated style="padding: 16px" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import ResourceDetailLayout from '@/components/ResourceDetailLayout.vue'
import { getResource, getPods, type GVR } from '@/apis/k8s'

// StatefulSet/DaemonSet 共用详情页：按路由 path 推断 kind 与 gvr
const route = useRoute()
const name = computed(() => String(route.params.name || ''))
const namespace = computed(() => String(route.query.namespace || 'default'))

const kind = computed(() => {
  if (route.path.includes('/statefulsets/')) return 'StatefulSet'
  if (route.path.includes('/daemonsets/')) return 'DaemonSet'
  if (route.path.includes('/replicasets/')) return 'ReplicaSet'
  return ''
})
const gvr = computed<GVR>(() => {
  if (kind.value === 'StatefulSet') return { group: 'apps', version: 'v1', resource: 'statefulsets' }
  if (kind.value === 'DaemonSet') return { group: 'apps', version: 'v1', resource: 'daemonsets' }
  if (kind.value === 'ReplicaSet') return { group: 'apps', version: 'v1', resource: 'replicasets' }
  return { version: 'v1', resource: '' }
})

const obj = ref<any>(null)
const pods = ref<any[]>([])
const podsLoading = ref(false)

const fetchObj = async () => {
  try {
    const res: any = await getResource(gvr.value, namespace.value, name.value)
    obj.value = res.data?.data
  } catch (e) {
    ElMessage.error('获取详情失败')
  }
}

// 关联 Pods：用后端解析的 owner 字段（PodInfo.Owner.kind/name）过滤
const fetchPods = async () => {
  podsLoading.value = true
  try {
    const res: any = await getPods(namespace.value)
    pods.value = (res.data?.data || []).filter((p: any) => p.owner?.kind === kind.value && p.owner?.name === name.value)
  } catch (e) {
    pods.value = []
  } finally {
    podsLoading.value = false
  }
}

onMounted(() => {
  fetchObj()
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
