<template>
  <div class="resource-detail">
    <!-- 标题卡片：返回 + 类型标签 + 资源名 + 状态 + 命名空间 + 操作 -->
    <div class="detail-header">
      <div class="header-left">
        <el-button class="back-btn" circle @click="goBack" title="返回">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <div class="title-block">
          <div class="title-row">
            <el-tag class="kind-tag" size="small" effect="plain" round>{{ kind }}</el-tag>
            <h2 class="detail-title" :title="title || name">{{ title || name }}</h2>
            <StatusTag v-if="status" :status="status" :type="statusType" show-icon />
          </div>
          <div class="subtitle">
            <span class="ns-label">命名空间</span>
            <span>{{ namespace || '—（集群级资源）' }}</span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <slot name="actions" />
        <el-button @click="refresh">
          <el-icon><Refresh /></el-icon><span style="margin-left:4px">刷新</span>
        </el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="detail-tabs">
      <el-tab-pane label="概览" name="overview">
        <div class="tab-content"><slot name="overview" /></div>
      </el-tab-pane>
      <el-tab-pane label="YAML" name="yaml">
        <div class="yaml-wrap">
          <YamlEditor v-if="yamlContent" v-model="yamlContent" readonly height="62vh" />
          <el-skeleton v-else :rows="8" animated />
        </div>
      </el-tab-pane>
      <el-tab-pane :label="events.length ? `事件 (${events.length})` : '事件'" name="events">
        <div class="tab-content">
          <el-table :data="events" size="default" v-loading="eventsLoading" max-height="560">
            <el-table-column label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="row.type === 'Warning' ? 'warning' : 'info'" size="small" effect="light">{{ row.type || '-' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="原因" width="170">
              <template #default="{ row }">
                <span :class="{ 'reason-warn': row.type === 'Warning' }">{{ row.reason || '-' }}</span>
              </template>
            </el-table-column>
            <el-table-column label="最后发生" width="150">
              <template #default="{ row }">
                <span :title="row.last_timestamp || row.lastTimestamp">{{ relTime(row.last_timestamp || row.lastTimestamp) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="次数" width="80" align="center">
              <template #default="{ row }">{{ row.count ?? '-' }}</template>
            </el-table-column>
            <el-table-column prop="message" label="消息" show-overflow-tooltip />
          </el-table>
          <el-empty v-if="!eventsLoading && events.length === 0" description="无事件" />
        </div>
      </el-tab-pane>
      <slot name="extra-tabs" />
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import yaml from 'js-yaml'
import StatusTag from './StatusTag.vue'
import YamlEditor from './YamlEditor.vue'
import { getResource, getEvents, type GVR } from '@/apis/k8s'

// 通用资源详情页布局：标题卡片 + 概览/YAML/Events 标签页。
// 各资源详情页通过 overview slot 定制概览，YAML 与 Events 由本组件统一加载。
// YAML 用 monaco（YamlEditor readonly）提供语法高亮。
const props = defineProps<{
  title?: string
  kind: string
  name: string
  namespace: string
  status?: string
  statusType?: 'pod' | 'deployment' | 'node' | 'service' | 'namespace'
  gvr: GVR
}>()

const emit = defineEmits<{ refresh: [] }>()
const router = useRouter()

const activeTab = ref('overview')
const yamlContent = ref('')
const events = ref<any[]>([])
const eventsLoading = ref(false)

const goBack = () => {
  if (window.history.length > 1) router.back()
}

const fetchYaml = async () => {
  try {
    const res: any = await getResource(props.gvr, props.namespace, props.name)
    const obj = res.data?.data ?? res.data
    yamlContent.value = yaml.dump(obj, { noRefs: true, lineWidth: 120 })
  } catch (e) {
    yamlContent.value = '# 加载 YAML 失败'
  }
}

const fetchEvents = async () => {
  eventsLoading.value = true
  try {
    const res: any = await getEvents(props.namespace, props.kind, props.name)
    events.value = res.data?.data || []
  } catch (e) {
    events.value = []
  } finally {
    eventsLoading.value = false
  }
}

const refresh = () => {
  fetchYaml()
  fetchEvents()
  emit('refresh')
}

// 相对时间：秒/分/时/天前，hover 显示原始时间
const relTime = (ts?: string) => {
  if (!ts) return '-'
  const t = new Date(ts).getTime()
  if (isNaN(t)) return String(ts)
  const s = Math.floor((Date.now() - t) / 1000)
  if (s < 0) return '刚刚'
  if (s < 60) return `${s} 秒前`
  if (s < 3600) return `${Math.floor(s / 60)} 分钟前`
  if (s < 86400) return `${Math.floor(s / 3600)} 小时前`
  return `${Math.floor(s / 86400)} 天前`
}

onMounted(() => {
  fetchYaml()
  fetchEvents()
})
</script>

<style scoped>
.resource-detail {
  padding: 16px 20px 24px;
}
/* 标题卡片 */
.detail-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  background: var(--el-bg-color, #fff);
  border: 1px solid var(--el-border-color-light, #ebeef5);
  border-radius: 8px;
  padding: 16px 20px;
  margin-bottom: 16px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
}
.header-left {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  min-width: 0;
  flex: 1;
}
.back-btn {
  margin-top: 2px;
  flex-shrink: 0;
}
.title-block {
  min-width: 0;
  flex: 1;
}
.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: nowrap;
  min-width: 0;
}
.title-row > * {
  flex-shrink: 0;
}
.kind-tag {
  font-weight: 600;
  letter-spacing: 0.3px;
}
.detail-title {
  margin: 0;
  font-size: 19px;
  font-weight: 600;
  line-height: 1.3;
  color: var(--el-text-color-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 0 1 auto;
  min-width: 0;
  max-width: 480px;
}
.subtitle {
  margin-top: 6px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  display: flex;
  gap: 6px;
  align-items: center;
}
.ns-label {
  color: var(--el-text-color-placeholder);
}
.header-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  align-items: center;
}
/* 标签页卡片 */
.detail-tabs {
  background: var(--el-bg-color, #fff);
  border: 1px solid var(--el-border-color-light, #ebeef5);
  border-radius: 8px;
  padding: 0 20px 20px;
}
.detail-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
}
.detail-tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
}
.tab-content {
  padding-top: 4px;
}
.yaml-wrap {
  padding-top: 4px;
}
.reason-warn {
  color: var(--el-color-warning);
  font-weight: 600;
}
</style>
