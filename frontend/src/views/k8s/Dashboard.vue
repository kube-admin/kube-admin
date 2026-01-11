<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <!-- 资源统计卡片 -->
      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#409EFF">
              <Platform />
            </el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.node_count }}</div>
              <div class="stats-label">集群节点</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#67C23A">
              <Box />
            </el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.namespace_count }}</div>
              <div class="stats-label">命名空间</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#E6A23C">
              <Box />
            </el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.pod_count }}</div>
              <div class="stats-label">Pod 总数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="6">
        <el-card class="stats-card">
          <div class="stats-content">
            <el-icon class="stats-icon" :size="40" color="#F56C6C">
              <Setting />
            </el-icon>
            <div class="stats-info">
              <div class="stats-value">{{ stats.deployment_count }}</div>
              <div class="stats-label">Deployment</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Pod 状态分布 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>Pod 状态分布</span>
            </div>
          </template>
          <div class="pod-status-stats">
            <div class="status-item">
              <el-tag type="success" size="large">Running</el-tag>
              <span class="status-count">{{ stats.pod_status_stats?.running || 0 }}</span>
            </div>
            <div class="status-item">
              <el-tag type="warning" size="large">Pending</el-tag>
              <span class="status-count">{{ stats.pod_status_stats?.pending || 0 }}</span>
            </div>
            <div class="status-item">
              <el-tag type="danger" size="large">Failed</el-tag>
              <span class="status-count">{{ stats.pod_status_stats?.failed || 0 }}</span>
            </div>
            <div class="status-item">
              <el-tag type="info" size="large">Succeeded</el-tag>
              <span class="status-count">{{ stats.pod_status_stats?.succeeded || 0 }}</span>
            </div>
            <div class="status-item">
              <el-tag size="large">Unknown</el-tag>
              <span class="status-count">{{ stats.pod_status_stats?.unknown || 0 }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>资源统计</span>
            </div>
          </template>
          <div class="resource-stats">
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
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快速操作 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快速操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" @click="goTo('/k8s/pods')">
              <el-icon><Box /></el-icon>
              管理 Pods
            </el-button>
            <el-button type="success" @click="goTo('/k8s/deployments')">
              <el-icon><Setting /></el-icon>
              管理 Deployments
            </el-button>
            <el-button type="warning" @click="goTo('/k8s/configmaps')">
              <el-icon><Document /></el-icon>
              管理 ConfigMaps
            </el-button>
            <el-button type="danger" @click="goTo('/k8s/secrets')">
              <el-icon><Lock /></el-icon>
              管理 Secrets
            </el-button>
            <el-button type="info" @click="goTo('/k8s/nodes')">
              <el-icon><Platform /></el-icon>
              查看 Nodes
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Platform, Box, Setting, Document, Lock } from '@element-plus/icons-vue'
import { getDashboardStats } from '@/apis/k8s'

const router = useRouter()
const stats = ref<any>({
  node_count: 0,
  namespace_count: 0,
  pod_count: 0,
  deployment_count: 0,
  service_count: 0,
  configmap_count: 0,
  secret_count: 0,
  pod_status_stats: {
    running: 0,
    pending: 0,
    failed: 0,
    succeeded: 0,
    unknown: 0
  }
})

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

onMounted(() => {
  fetchStats()
  // 每30秒刷新一次
  setInterval(fetchStats, 30000)
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}

.pod-status-stats {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
}

.status-count {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.resource-stats {
  padding: 10px 0;
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
  color: #409EFF;
}

.quick-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
