<template>
  <el-tag :type="tagType" :size="size" :effect="effect">
    <el-icon v-if="showIcon" style="margin-right: 3px">
      <component :is="iconComponent" />
    </el-icon>
    {{ displayText }}
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
