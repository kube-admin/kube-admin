<template>
  <div class="resource-actions">
    <el-button 
      v-if="showView" 
      :size="size" 
      @click="$emit('view', resource)"
    >
      查看
    </el-button>
    
    <el-button 
      v-if="showEdit" 
      :size="size" 
      type="primary" 
      @click="$emit('edit', resource)"
    >
      编辑
    </el-button>
    
    <el-button 
      v-if="showLogs && resource.type === 'pod'" 
      :size="size" 
      type="info" 
      @click="$emit('logs', resource)"
    >
      日志
    </el-button>
    
    <el-button 
      v-if="showScale && resource.type === 'deployment'" 
      :size="size" 
      type="warning" 
      @click="$emit('scale', resource)"
    >
      扩缩容
    </el-button>
    
    <el-popconfirm
      v-if="showDelete && !disabled"
      :title="`确定删除 ${resource.name} 吗?`"
      @confirm="$emit('delete', resource)"
    >
      <template #reference>
        <el-button :size="size" type="danger">删除</el-button>
      </template>
    </el-popconfirm>
    
    <el-tooltip v-if="showDelete && disabled" content="此资源不可删除" placement="top">
      <el-button :size="size" type="danger" disabled>删除</el-button>
    </el-tooltip>
    
    <slot name="extra" :resource="resource"></slot>
  </div>
</template>

<script setup lang="ts">
interface Props {
  resource: any
  size?: 'small' | 'default' | 'large'
  showView?: boolean
  showEdit?: boolean
  showDelete?: boolean
  showLogs?: boolean
  showScale?: boolean
  disabled?: boolean
}

withDefaults(defineProps<Props>(), {
  size: 'small',
  showView: true,
  showEdit: false,
  showDelete: true,
  showLogs: false,
  showScale: false,
  disabled: false
})

defineEmits<{
  view: [resource: any]
  edit: [resource: any]
  delete: [resource: any]
  logs: [resource: any]
  scale: [resource: any]
}>()
</script>

<style scoped>
.resource-actions {
  display: inline-flex;
  gap: 5px;
}
</style>
