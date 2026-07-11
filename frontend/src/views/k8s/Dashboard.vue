<template>
  <div class="dashboard-container">
    <!-- 资源统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover" @click="goTo('/k8s/nodes')">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#409EFF"><Platform /></el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.node_count }}</div>
              <div class="stats-label">集群节点</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover" @click="goTo('/k8s/namespaces')">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#67C23A"><Box /></el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.namespace_count }}</div>
              <div class="stats-label">命名空间</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover" @click="goTo('/k8s/pods')">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#E6A23C"><Box /></el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.pod_count }}</div>
              <div class="stats-label">Pod 总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover" @click="goTo('/k8s/deployments')">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#F56C6C"><Setting /></el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.deployment_count }}</div>
              <div class="stats-label">Deployment</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 集群资源使用率 + Pod 状态分布 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card>
          <template #header><div class="card-header">CPU 使用率</div></template>
          <v-chart class="chart" :option="cpuGaugeOption" autoresize />
          <div class="usage-detail">{{ fmtCpu(stats.cluster_usage?.cpu_used) }} / {{ fmtCpu(stats.cluster_usage?.cpu_capacity) }} 核</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <template #header><div class="card-header">内存使用率</div></template>
          <v-chart class="chart" :option="memGaugeOption" autoresize />
          <div class="usage-detail">{{ fmtMem(stats.cluster_usage?.memory_used) }} / {{ fmtMem(stats.cluster_usage?.memory_capacity) }}</div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header><div class="card-header">Pod 状态分布</div></template>
          <v-chart class="chart" :option="pieOption" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <el-alert
      v-if="!hasMetrics"
      title="未检测到实时资源使用率，请确认集群已安装 metrics-server"
      type="warning"
      :closable="false"
      style="margin-top: 15px"
    />

    <!-- 资源统计 + 快速操作 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header><div class="card-header">资源统计</div></template>
          <el-row :gutter="10">
            <el-col :span="8">
              <div class="resource-item">
                <div class="resource-label">Service</div>
                <div class="resource-value">{{ stats.service_count }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="resource-item">
                <div class="resource-label">ConfigMap</div>
                <div class="resource-value">{{ stats.configmap_count }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="resource-item">
                <div class="resource-label">Secret</div>
                <div class="resource-value">{{ stats.secret_count }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header><div class="card-header">快速操作</div></template>
          <div class="quick-actions">
            <el-button type="primary" @click="goTo('/k8s/pods')">管理 Pods</el-button>
            <el-button type="success" @click="goTo('/k8s/deployments')">管理 Deployments</el-button>
            <el-button type="warning" @click="goTo('/k8s/configmaps')">管理 ConfigMaps</el-button>
            <el-button type="danger" @click="goTo('/k8s/secrets')">管理 Secrets</el-button>
            <el-button type="info" @click="goTo('/k8s/nodes')">查看 Nodes</el-button>
            <el-button @click="goTo('/k8s/events')">查看 Events</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Platform, Box, Setting } from '@element-plus/icons-vue'
import { getDashboardStats } from '@/apis/k8s'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart, GaugeChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'

use([CanvasRenderer, PieChart, GaugeChart, TooltipComponent, LegendComponent])

const router = useRouter()
const stats = ref<any>({
  node_count: 0,
  namespace_count: 0,
  pod_count: 0,
  deployment_count: 0,
  service_count: 0,
  configmap_count: 0,
  secret_count: 0,
  pod_status_stats: { running: 0, pending: 0, failed: 0, succeeded: 0, unknown: 0 },
  cluster_usage: { cpu_percent: 0, memory_percent: 0, cpu_used: '', cpu_capacity: '', memory_used: '', memory_capacity: '' }
})

let timer: ReturnType<typeof setInterval> | null = null

// 是否已获取到 metrics 数据
const hasMetrics = computed(() => !!(stats.value.cluster_usage?.cpu_used || stats.value.cluster_usage?.memory_used))

const fetchStats = async () => {
  try {
    const res = await getDashboardStats()
    stats.value = res.data.data || stats.value
  } catch (error) {
    ElMessage.error('获取集群统计信息失败')
  }
}

const goTo = (path: string) => {
  router.push(path)
}

// CPU Quantity（k8s resource.Quantity 字符串）→ cores，如 "764075769n"(纳核) → 0.76，"500m" → 0.5，"2" → 2.00
const fmtCpu = (q?: string | number) => {
  if (q === undefined || q === null || q === '') return '-'
  const s = String(q)
  if (s.endsWith('n')) return (parseInt(s) / 1e9).toFixed(2)
  if (s.endsWith('u')) return (parseInt(s) / 1e6).toFixed(3)
  if (s.endsWith('m')) return (parseInt(s) / 1e3).toFixed(2)
  return parseFloat(s).toFixed(2)
}
// 内存 Quantity → Gi，如 "15096864Ki" → 14.40 Gi
const UNIT_BYTES: Record<string, number> = { Ki: 1024, Mi: 1024 ** 2, Gi: 1024 ** 3, Ti: 1024 ** 4 }
const fmtMem = (q?: string | number) => {
  if (q === undefined || q === null || q === '') return '-'
  const s = String(q)
  const m = s.match(/^(\d+(?:\.\d+)?)(Ki|Mi|Gi|Ti)?$/)
  if (!m) return s
  const bytes = parseFloat(m[1]) * (m[2] ? UNIT_BYTES[m[2]] : 1)
  return (bytes / 1024 ** 3).toFixed(2) + ' Gi'
}

// Pod 状态饼图
const pieOption = computed(() => ({
  tooltip: { trigger: 'item' },
  legend: { bottom: 0 },
  series: [
    {
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      label: { show: true, formatter: '{b}: {c}' },
      data: [
        { value: stats.value.pod_status_stats?.running || 0, name: 'Running', itemStyle: { color: '#67C23A' } },
        { value: stats.value.pod_status_stats?.pending || 0, name: 'Pending', itemStyle: { color: '#E6A23C' } },
        { value: stats.value.pod_status_stats?.failed || 0, name: 'Failed', itemStyle: { color: '#F56C6C' } },
        { value: stats.value.pod_status_stats?.succeeded || 0, name: 'Succeeded', itemStyle: { color: '#909399' } },
        { value: stats.value.pod_status_stats?.unknown || 0, name: 'Unknown', itemStyle: { color: '#C0C4CC' } }
      ]
    }
  ]
}))

// 仪表盘构造
const buildGauge = (title: string, percent: number, color: string) => ({
  series: [
    {
      type: 'gauge',
      startAngle: 200,
      endAngle: -20,
      min: 0,
      max: 100,
      progress: { show: true, width: 14 },
      axisLine: { lineStyle: { width: 14 } },
      axisTick: { show: false },
      splitLine: { length: 10, lineStyle: { width: 2 } },
      pointer: { width: 4 },
      detail: { valueAnimation: true, formatter: '{value}%', fontSize: 18, offsetCenter: [0, '40%'] },
      title: { show: true, offsetCenter: [0, '70%'], fontSize: 13, color: '#909399' },
      data: [{ value: Math.round(percent), name: title }],
      itemStyle: { color }
    }
  ]
})

const cpuGaugeOption = computed(() => buildGauge('CPU', stats.value.cluster_usage?.cpu_percent || 0, '#409EFF'))
const memGaugeOption = computed(() => buildGauge('内存', stats.value.cluster_usage?.memory_percent || 0, '#67C23A'))

onMounted(() => {
  fetchStats()
  timer = setInterval(fetchStats, 30000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}
.stats-card {
  cursor: pointer;
  transition: transform 0.2s;
}
.stats-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
.stats-content {
  display: flex;
  align-items: center;
  padding: 10px 0;
}
.stats-icon {
  margin-right: 20px;
}
.stats-info {
  flex: 1;
}
.stats-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}
.stats-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}
.card-header {
  font-weight: bold;
}
.chart {
  height: 240px;
}
.usage-detail {
  text-align: center;
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}
.resource-item {
  text-align: center;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}
.resource-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 10px;
}
.resource-value {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
}
.quick-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
