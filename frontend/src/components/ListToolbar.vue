<template>
  <div class="list-toolbar">
    <span class="list-title">{{ title }}</span>
    <div class="toolbar-actions">
      <slot name="filters" />
      <slot />
      <el-button :icon="Refresh" :loading="loading" @click="emit('refresh')">刷新</el-button>
      <AutoRefresh :interval="interval" @refresh="emit('refresh')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Refresh } from '@element-plus/icons-vue'
import AutoRefresh from './AutoRefresh.vue'

// 通用列表页工具栏：标题（左）+ 具名默认 slot 放创建等按钮 + 刷新 + 自动刷新（右）。
// 所有列表页复用，统一 header 风格，消除"有的页有刷新/自动刷新、有的没有"的不一致。
defineProps<{
  title: string
  loading?: boolean
  interval?: number
}>()
const emit = defineEmits<{ refresh: [] }>()
</script>

<style scoped>
.list-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
}
.list-title {
  color: var(--el-text-color-primary);
}
.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
