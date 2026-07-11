<template>
  <el-tag :type="tagType" :size="size" :effect="effect" class="status-tag">
    <el-icon v-if="showIcon" class="status-tag__icon">
      <component :is="iconComponent" />
    </el-icon>
    <span class="status-tag__text">{{ displayText }}</span>
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  CircleCheck,
  CircleClose,
  Warning,
  Loading as LoadingIcon,
  QuestionFilled
} from '@element-plus/icons-vue'

interface Props {
  status: string
  type?: 'pod' | 'deployment' | 'node' | 'service' | 'namespace'
  size?: 'small' | 'default' | 'large'
  effect?: 'dark' | 'light' | 'plain'
  showIcon?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'pod',
  size: 'default',
  effect: 'light',
  showIcon: false
})

const tagType = computed(() => {
  const status = props.status.toLowerCase()

  switch (props.type) {
    case 'pod':
      if (['running', 'succeeded'].includes(status)) return 'success'
      if (['pending', 'containercreating'].includes(status)) return 'warning'
      if (['failed', 'error', 'crashloopbackoff'].includes(status)) return 'danger'
      return 'info'

    case 'deployment':
      if (status === 'ready') return 'success'
      if (status === 'updating') return 'warning'
      return 'danger'

    case 'node':
      return status === 'ready' ? 'success' : 'danger'

    case 'service':
      const typeMap: Record<string, string> = {
        'clusterip': 'success',
        'nodeport': 'warning',
        'loadbalancer': 'primary'
      }
      return typeMap[status] || 'info'

    case 'namespace':
      return status === 'active' ? 'success' : 'danger'

    default:
      return 'info'
  }
})

const iconComponent = computed(() => {
  const type = tagType.value
  switch (type) {
    case 'success':
      return CircleCheck
    case 'danger':
      return CircleClose
    case 'warning':
      return Warning
    case 'info':
      return QuestionFilled
    default:
      return LoadingIcon
  }
})

const displayText = computed(() => {
  return props.status || 'Unknown'
})
</script>

<style scoped>
/* 强制图标与文字同行不换行，避免在标题栏空间紧张时文字被挤换行 */
.status-tag {
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
}
/* 覆盖 element-plus el-tag 内部 __content 的固定 width:50px，否则图标+文字塞不下会换行 */
.status-tag :deep(.el-tag__content) {
  width: auto !important;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
}
.status-tag__icon {
  margin-right: 3px;
  flex-shrink: 0;
}
.status-tag__text {
  line-height: 1;
}
</style>
