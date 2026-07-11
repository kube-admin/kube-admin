<template>
  <div class="terminal-wrapper">
    <div ref="hostRef" class="terminal-host"></div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import '@xterm/xterm/css/xterm.css'
import { ElMessage } from 'element-plus'
import { getPodTerminalUrl } from '@/apis/k8s'

// 可复用的交互式终端组件：基于 xterm.js + WebSocket 裸字节透传到后端 exec TTY。
// 容器切换由父组件用 :key 重建实例实现，本组件不内部 watch container。
const props = defineProps<{
  podName: string
  namespace: string
  container: string
}>()

const hostRef = ref<HTMLDivElement | null>(null)
let term: Terminal | null = null
let fitAddon: FitAddon | null = null
let ws: WebSocket | null = null

const sendResize = (cols: number, rows: number) => {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'resize', cols, rows }))
  }
}

const fit = () => {
  try {
    fitAddon?.fit()
  } catch (e) {
    // 容器尺寸未就绪时可能抛错，忽略
  }
}

const onWindowResize = () => fit()

onMounted(() => {
  if (!hostRef.value) return
  term = new Terminal({
    cursorBlink: true,
    fontSize: 13,
    fontFamily: "'Courier New', monospace",
    theme: { background: '#1e1e1e', foreground: '#ffffff' }
  })
  fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.open(hostRef.value)
  // 按键直达容器 stdin（支持 vim/top/Ctrl+C/Tab 等交互式程序）
  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) ws.send(data)
  })
  term.onResize(({ cols, rows }) => sendResize(cols, rows))
  fit()

  // 连接后端 exec WebSocket
  const wsUrl = getPodTerminalUrl(props.podName, props.namespace, props.container)
  ws = new WebSocket(wsUrl)
  ws.onopen = () => {
    // 连接建立后把当前终端尺寸同步给后端 TTY
    if (term) sendResize(term.cols, term.rows)
    term?.focus()
  }
  ws.onmessage = (event) => {
    // 后端裸字节透传，直接写入 xterm（完整 ANSI 渲染）
    term?.write(typeof event.data === 'string' ? event.data : '')
  }
  ws.onerror = () => ElMessage.error('终端连接错误: 网络错误或连接被拒绝')
  ws.onclose = (event) => {
    if (event.code !== 1000 && event.code !== 1005) {
      ElMessage.warning('终端连接已关闭，代码: ' + event.code)
    }
  }

  window.addEventListener('resize', onWindowResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onWindowResize)
  if (ws) {
    ws.close()
    ws = null
  }
  term?.dispose()
  term = null
  fitAddon = null
})
</script>

<style scoped>
.terminal-wrapper {
  width: 100%;
  height: 100%;
  background: #1e1e1e;
  padding: 6px;
  box-sizing: border-box;
}
.terminal-host {
  width: 100%;
  height: 100%;
}
</style>
