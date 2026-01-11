<template>
  <div class="home-container">
    <!-- 系统标题 -->
    <div class="hero-section">
      <h1 class="system-title">Kubernetes 管理平台</h1>
      <p class="system-subtitle">一站式 Kubernetes 集群管理和资源监控平台</p>
    </div>

    <!-- 系统简介 -->
    <el-card class="intro-card">
      <template #header>
        <div class="card-header">
          <span>系统介绍</span>
        </div>
      </template>
      <div class="intro-content">
        <p>Kubernetes 管理平台是一个功能强大的多集群管理工具，提供直观的 Web 界面来管理您的 Kubernetes 集群和资源。</p>
        <p>通过本平台，您可以轻松地：</p>
        <ul>
          <li>管理多个 Kubernetes 集群</li>
          <li>监控集群状态和资源使用情况</li>
          <li>管理 Namespaces、Pods、Deployments、Services 等核心资源</li>
          <li>查看和编辑 ConfigMaps、Secrets 等配置资源</li>
          <li>实时查看集群节点状态</li>
        </ul>
      </div>
    </el-card>

    <!-- 快速导航 -->
    <el-card class="quick-nav-card">
      <template #header>
        <div class="card-header">
          <span>快速开始</span>
        </div>
      </template>
      <div class="quick-nav-grid">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="nav-item" @click="$router.push('/k8s/dashboard')">
              <el-icon class="nav-icon" :size="32" color="#409EFF"><Platform /></el-icon>
              <div class="nav-text">集群概览</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="nav-item" @click="$router.push('/k8s/clusters')">
              <el-icon class="nav-icon" :size="32" color="#67C23A"><Management /></el-icon>
              <div class="nav-text">集群管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="nav-item" @click="$router.push('/k8s/namespaces')">
              <el-icon class="nav-icon" :size="32" color="#E6A23C"><Collection /></el-icon>
              <div class="nav-text">命名空间</div>
            </div>
          </el-col>
        </el-row>
        <el-row :gutter="20" style="margin-top: 20px">
          <el-col :span="8">
            <div class="nav-item" @click="$router.push('/k8s/pods')">
              <el-icon class="nav-icon" :size="32" color="#F56C6C"><Box /></el-icon>
              <div class="nav-text">Pod 管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="nav-item" @click="$router.push('/k8s/deployments')">
              <el-icon class="nav-icon" :size="32" color="#909399"><Setting /></el-icon>
              <div class="nav-text">部署管理</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="nav-item" @click="$router.push('/k8s/nodes')">
              <el-icon class="nav-icon" :size="32" color="#79bbff"><Monitor /></el-icon>
              <div class="nav-text">节点监控</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- 系统特性 -->
    <el-card class="features-card">
      <template #header>
        <div class="card-header">
          <span>核心特性</span>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :span="8">
          <div class="feature-item">
            <el-icon class="feature-icon" :size="24" color="#409EFF"><Connection /></el-icon>
            <div class="feature-content">
              <h4>多集群支持</h4>
              <p>支持管理多个 Kubernetes 集群，轻松切换不同环境</p>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="feature-item">
            <el-icon class="feature-icon" :size="24" color="#67C23A"><Monitor /></el-icon>
            <div class="feature-content">
              <h4>实时监控</h4>
              <p>实时监控集群状态和资源使用情况</p>
            </div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="feature-item">
            <el-icon class="feature-icon" :size="24" color="#E6A23C"><Edit /></el-icon>
            <div class="feature-content">
              <h4>资源管理</h4>
              <p>全面管理各种 Kubernetes 资源</p>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { Platform, Management, Collection, Box, Setting, Monitor, Connection, Edit } from '@element-plus/icons-vue'
</script>

<style scoped lang="scss">
.home-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.hero-section {
  text-align: center;
  padding: 40px 0;
  margin-bottom: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
}

.system-title {
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 10px;
}

.system-subtitle {
  font-size: 18px;
  opacity: 0.9;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
  display: flex;
  align-items: center;
}

.intro-card {
  margin-bottom: 20px;
}

.intro-content p {
  font-size: 16px;
  line-height: 1.6;
  margin-bottom: 15px;
}

.intro-content ul {
  padding-left: 20px;
  margin-bottom: 15px;
}

.intro-content li {
  margin-bottom: 8px;
}

.quick-nav-card {
  margin-bottom: 20px;
}

.quick-nav-grid {
  padding: 20px 0;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  height: 120px;
}

.nav-item:hover {
  border-color: #409EFF;
  box-shadow: 0 2px 12px 0 rgba(64, 158, 255, 0.2);
  transform: translateY(-2px);
}

.nav-icon {
  margin-bottom: 10px;
}

.nav-text {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.features-card {
  margin-bottom: 20px;
}

.feature-item {
  display: flex;
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  margin-bottom: 15px;
}

.feature-icon {
  margin-right: 15px;
  flex-shrink: 0;
}

.feature-content h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
}

.feature-content p {
  margin: 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
}
</style>