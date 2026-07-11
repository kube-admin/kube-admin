<template>
  <div class="pod-terminal-page">
    <div class="terminal-toolbar">
      <el-button size="small" @click="goBack">← 返回</el-button>
      <span class="pod-name">{{ podName }}</span>
      <span class="pod-ns">/{{ namespace }}</span>
      <el-select
        v-if="containers.length > 1"
        v-model="selectedContainer"
        size="small"
        class="container-select"
      >
        <el-option v-for="c in containers" :key="c.name" :label="c.name" :value="c.name" />
      </el-select>
      <span v-if="loading" class="loading-text">加载容器列表...</span>
    </div>
    <div class="terminal-body">
      <Terminal
        v-if="selectedContainer && !loading"
        :key="podName + '/' + selectedContainer"
        :pod-name="podName"
        :namespace="namespace"
        :container="selectedContainer"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Terminal from '@/components/Terminal.vue'
import { getPodDetail } from '@/apis/k8s'

// Pod 终端全屏页：独立顶级路由，无侧边栏，纯全屏交互。
const route = useRoute()
const router = useRouter()

const podName = computed(() => String(route.params.name || ''))
const namespace = computed(() => String(route.query.namespace || 'default'))

const containers = ref<any[]>([])
const selectedContainer = ref('')
const loading = ref(true)

const goBack = () => {
  // 有历史则后退，否则回 Pods 列表
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/k8s/pods')
  }
}

onMounted(async () => {
  try {
    const res: any = await getPodDetail(podName.value, namespace.value)
    const pod = res.data?.data
    containers.value = pod?.containers || []
    if (containers.value.length === 0) {
      ElMessage.error('未获取到容器列表')
      return
    }
    const queryContainer = route.query.container ? String(route.query.container) : ''
    selectedContainer.value = queryContainer || containers.value[0].name
  } catch (e) {
    ElMessage.error('获取 Pod 详情失败')
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.pod-terminal-page {
  position: fixed;
  inset: 0;
  display: flex;
  flex-direction: column;
  background: #fff;
  z-index: 2000;
}
.terminal-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-bottom: 1px solid #ebeef5;
  background: #f5f7fa;
}
.pod-name {
  font-weight: 600;
}
.pod-ns {
  color: #909399;
  font-size: 13px;
}
.container-select {
  width: 200px;
}
.loading-text {
  color: #909399;
  font-size: 13px;
}
.terminal-body {
  flex: 1;
  min-height: 0;
}
</style>
