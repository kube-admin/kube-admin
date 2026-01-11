<template>
  <div class="auto-refresh">
    <el-switch
      v-model="enabled"
      active-text="自动刷新"
      @change="handleToggle"
    />
    <el-select 
      v-if="enabled" 
      v-model="selectedInterval" 
      size="small" 
      style="width: 100px; margin-left: 10px"
      @change="handleIntervalChange"
    >
      <el-option label="5秒" :value="5" />
      <el-option label="10秒" :value="10" />
      <el-option label="30秒" :value="30" />
      <el-option label="60秒" :value="60" />
    </el-select>
    <span v-if="enabled" style="margin-left: 10px; color: #909399; font-size: 12px">
      下次刷新: {{ countdown }}s
    </span>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'

interface Props {
  interval?: number
  autoStart?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  interval: 30,
  autoStart: false
})

const emit = defineEmits<{
  refresh: []
}>()

const enabled = ref(props.autoStart)
const selectedInterval = ref(props.interval)
const countdown = ref(selectedInterval.value)
let timer: ReturnType<typeof setInterval> | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const startRefresh = () => {
  if (timer) clearInterval(timer)
  if (countdownTimer) clearInterval(countdownTimer)
  
  countdown.value = selectedInterval.value
  
  // 倒计时
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      countdown.value = selectedInterval.value
    }
  }, 1000)
  
  // 刷新定时器
  timer = setInterval(() => {
    emit('refresh')
  }, selectedInterval.value * 1000)
}

const stopRefresh = () => {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
}

const handleToggle = (value: boolean) => {
  if (value) {
    startRefresh()
  } else {
    stopRefresh()
  }
}

const handleIntervalChange = () => {
  if (enabled.value) {
    stopRefresh()
    startRefresh()
  }
}

watch(() => props.autoStart, (value) => {
  enabled.value = value
  if (value) {
    startRefresh()
  } else {
    stopRefresh()
  }
})

onMounted(() => {
  if (props.autoStart) {
    startRefresh()
  }
})

onUnmounted(() => {
  stopRefresh()
})
</script>

<style scoped>
.auto-refresh {
  display: inline-flex;
  align-items: center;
}
</style>
