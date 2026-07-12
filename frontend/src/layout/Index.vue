<template>
  <el-watermark :font="font" :content="watermarkContent">
    <el-container>
      <Sidebar />
      <el-container>
        <el-header>
          <Header />
        </el-header>
        <TagView />
        <el-main>
          <RouterView />
        </el-main>
        <Footer />
      </el-container>
    </el-container>
  </el-watermark>
</template>

<script setup lang="ts">
import Header from './components/Header/Index.vue'
import TagView from './components/TagsView/Index.vue'
import Sidebar from './components/Sidebar/Index.vue'
import Footer from './components/Footer.vue'
import { reactive, watch } from 'vue'
import { isDark } from '@/stores/dark'

const font = reactive({
  color: 'rgba(0, 0, 0, .05)'
})

watch(
  isDark,
  () => {
    font.color = isDark.value
      ? 'rgba(255, 255, 255, .05)'
      : 'rgba(0, 0, 0, .05)'
  },
  {
    immediate: true
  }
)

// 水印两行：上=当前用户名，下=KubeAdmin（泄露溯源 + 品牌曝光）
const username = (() => {
  try {
    return JSON.parse(localStorage.getItem('user') || '{}')?.username || 'Kube Admin'
  } catch {
    return 'Kube Admin'
  }
})()
const watermarkContent = [username, 'KubeAdmin']
</script>
