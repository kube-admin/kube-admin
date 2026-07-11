<template>
  <router-view />
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useNamespaceStore } from '@/stores/namespace'

// 集群/命名空间 URL 状态化：
// - 打开分享链接时从 URL query 恢复集群/命名空间（进对应上下文）
// - Header 切换时同步到 URL（生成的链接可还原）
const router = useRouter()
const route = useRoute()
const nsStore = useNamespaceStore()

// URL → store 恢复（响应式）：main.ts 未 await router.isReady() 即挂载，App.vue setup
// 执行时初始路由尚未解析、route.query 为空；一次性读取会读到空值导致恢复失效。
// 改为 watch route.query：路由解析完成（及后续 SPA 导航到分享链接）时触发恢复。
watch(() => route.query, (q) => {
  if (q.cluster_id !== undefined && q.cluster_id !== '') {
    const cid = Number(q.cluster_id)
    if (!isNaN(cid)) nsStore.setCurrentClusterId(cid)
  }
  if (typeof q.namespace === 'string' && q.namespace) {
    nsStore.setCurrentNamespace(q.namespace)
  }
}, { immediate: true })

// store → URL 同步：Header 切换集群/命名空间时写回 URL，使生成的链接可还原
watch(() => nsStore.currentNamespace, (v) => {
  router.replace({ query: { ...route.query, namespace: v } })
})
watch(() => nsStore.currentClusterId, (v) => {
  const q: Record<string, any> = { ...route.query }
  if (v) q.cluster_id = String(v)
  else delete q.cluster_id
  router.replace({ query: q })
})
</script>

<style lang="scss"></style>