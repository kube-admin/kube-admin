<template>
  <ResourceDetailLayout
    v-if="obj"
    :title="name"
    :kind="obj.kind || resource"
    :name="name"
    :namespace="namespace"
    :gvr="gvr"
    @refresh="fetchObj"
  >
    <template #overview>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="Kind">{{ obj.kind || '-' }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ obj.metadata?.name }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ namespace || '—（集群级资源）' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ obj.metadata?.creationTimestamp }}</el-descriptions-item>
        <el-descriptions-item label="UID" :span="2">{{ obj.metadata?.uid }}</el-descriptions-item>
        <el-descriptions-item v-if="obj.metadata?.labels && Object.keys(obj.metadata.labels).length" label="Labels" :span="2">
          <el-tag
            v-for="(v, k) in obj.metadata.labels"
            :key="k"
            size="small"
            class="label-tag"
          >{{ k }}={{ v }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item v-if="obj.metadata?.annotations && Object.keys(obj.metadata.annotations).length" label="Annotations" :span="2">
          <div v-for="(v, k) in obj.metadata.annotations" :key="k" class="anno-line" :title="`${k}=${v}`">
            <b>{{ k }}</b>={{ v }}
          </div>
        </el-descriptions-item>
      </el-descriptions>
      <p class="hint">完整定义见「YAML」标签；资源事件见「事件」标签。</p>
    </template>
  </ResourceDetailLayout>
  <el-skeleton v-else :rows="8" animated style="padding: 16px" />
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import ResourceDetailLayout from '@/components/ResourceDetailLayout.vue'
import { getResource, type GVR } from '@/apis/k8s'

// 通用资源详情页：适用于任意 GVR（无专门详情页的资源走此页）。
// GVR + name 从路由读取：path = /k8s/resource-detail/:resource/:name，
// group/version/namespace 经 query 传递（core 资源 group 为空，集群级资源 namespace 为空）。
const route = useRoute()
const resource = computed(() => String(route.params.resource || ''))
const name = computed(() => String(route.params.name || ''))
const group = computed(() => String(route.query.group || ''))
const version = computed(() => String(route.query.version || 'v1'))
const namespace = computed(() => String(route.query.namespace || ''))
const gvr = computed<GVR>(() => ({
  group: group.value,
  version: version.value,
  resource: resource.value
}))

const obj = ref<any>(null)
const fetchObj = async () => {
  try {
    const res: any = await getResource(gvr.value, namespace.value, name.value)
    obj.value = res.data?.data
  } catch (e) {
    ElMessage.error('获取详情失败')
  }
}
onMounted(fetchObj)
</script>

<style scoped>
.label-tag {
  margin: 2px 4px 2px 0;
}
.anno-line {
  font-size: 12px;
  color: var(--el-text-color-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.hint {
  margin: 12px 0 0;
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}
</style>
