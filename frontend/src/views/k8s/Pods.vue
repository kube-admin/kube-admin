<template>
  <div class="pods-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>Pod 列表</span>
          <div>
            <el-select v-model="namespaceStore.currentNamespace" placeholder="选择命名空间" style="width: 200px; margin-right: 10px">
              <el-option
                v-for="ns in namespaceStore.namespaces"
                :key="ns.name"
                :label="ns.name"
                :value="ns.name"
              />
            </el-select>
            <el-button type="primary" @click="showCreateDialog">创建 Pod</el-button>
          </div>
        </div>
      </template>

      <el-table :data="pods" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" width="250" />
        <el-table-column prop="namespace" label="命名空间" width="150" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ scope.row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pod_ip" label="Pod IP" width="150" />
        <el-table-column prop="node_name" label="节点" width="180" />
        <el-table-column prop="creation_timestamp" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="250">
          <template #default="scope">
            <el-button size="small" @click="showLogs(scope.row)">日志</el-button>
            <el-button size="small" @click="showTerminal(scope.row)">终端</el-button>
            <el-popconfirm
              title="确定删除这个Pod吗?"
              @confirm="handleDelete(scope.row)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="Pod 日志" width="70%">
      <el-select v-model="selectedContainer" placeholder="选择容器" style="margin-bottom: 10px">
        <el-option
          v-for="container in currentPodContainers"
          :key="container.name"
          :label="container.name"
          :value="container.name"
        />
      </el-select>
      <el-input
        v-model="logs"
        type="textarea"
        :rows="20"
        readonly
        style="font-family: monospace"
      />
    </el-dialog>

    <!-- 终端对话框 -->
    <el-dialog v-model="terminalDialogVisible" title="Pod 终端" width="70%" @close="closeTerminal">
      <div class="terminal-container">
        <div ref="terminalRef" class="terminal-output"></div>
        <div class="terminal-input-container">
          <span class="prompt">$</span>
          <input 
            ref="terminalInputRef" 
            v-model="terminalInput" 
            @keyup.enter="sendCommand" 
            @keydown.up="handleUpArrow"
            @keydown.down="handleDownArrow"
            class="terminal-input"
            placeholder="输入命令..."
            autocomplete="off"
          />
        </div>
      </div>
    </el-dialog>

    <!-- 创建 Pod 对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建 Pod (YAML)" width="70%">
      <el-input
        v-model="podYaml"
        type="textarea"
        :rows="20"
        placeholder="请输入 Pod 的 YAML 定义..."
        style="font-family: monospace"
      />
      <template #footer>
        <el-button @click="loadPodTemplate">加载模板</el-button>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { getPods, deletePod, getPodLogs, getNamespaces, createPodFromYaml, getPodTerminalUrl } from '@/apis/k8s'
import { useNamespaceStore } from '@/stores/namespace'
import type { AxiosResponse } from 'axios'
import type { Result } from '@/apis/client/request'

const namespaceStore = useNamespaceStore()

const loading = ref(false)
const pods = ref<any[]>([])
const logDialogVisible = ref(false)
const terminalDialogVisible = ref(false)
const logs = ref('')
const selectedContainer = ref('')
const currentPodContainers = ref<any[]>([])
const currentPod = ref<any>(null)
const createDialogVisible = ref(false)
const creating = ref(false)
const podYaml = ref('')

// 终端相关
const terminalRef = ref<HTMLDivElement | null>(null)
const terminalInputRef = ref<HTMLInputElement | null>(null)
const terminalInput = ref('')
const commandHistory = ref<string[]>([])
const historyIndex = ref(-1)
let ws: WebSocket | null = null
let currentContainer = ''

// 获取命名空间列表
const fetchNamespaces = async () => {
  try {
    const res: AxiosResponse<Result<any[]>> = await getNamespaces()
    namespaceStore.setNamespaces(res.data.data || [])
  } catch (error) {
    ElMessage.error('获取命名空间失败')
  }
}

// 获取Pod列表
const fetchPods = async () => {
  loading.value = true
  try {
    const res: AxiosResponse<Result<any[]>> = await getPods(namespaceStore.currentNamespace)
    pods.value = res.data.data || []
  } catch (error) {
    ElMessage.error('获取Pod列表失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    'Running': 'success',
    'Pending': 'warning',
    'Failed': 'danger',
    'Succeeded': 'info'
  }
  return typeMap[status] || 'info'
}

const handleDelete = async (pod: any) => {
  try {
    await deletePod(pod.name, pod.namespace)
    ElMessage.success('删除成功')
    fetchPods()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const showLogs = async (pod: any) => {
  currentPod.value = pod
  currentPodContainers.value = pod.containers || []
  if (currentPodContainers.value.length > 0) {
    selectedContainer.value = currentPodContainers.value[0].name
  }
  logDialogVisible.value = true
  await fetchLogs()
}

const fetchLogs = async () => {
  if (!currentPod.value || !selectedContainer.value) return
  try {
    const res: AxiosResponse<Result<{ logs: string }>> = await getPodLogs(
      currentPod.value.name,
      currentPod.value.namespace,
      selectedContainer.value,
      100
    )
    logs.value = res.data.data?.logs || ''
  } catch (error) {
    ElMessage.error('获取日志失败')
  }
}

// 显示终端对话框
const showTerminal = (pod: any) => {
  currentPod.value = pod
  currentPodContainers.value = pod.containers || []
  if (currentPodContainers.value.length > 0) {
    currentContainer = currentPodContainers.value[0].name
  }
  terminalDialogVisible.value = true
  terminalInput.value = ''
  commandHistory.value = []
  historyIndex.value = -1
  
  // 在下一次DOM更新后连接WebSocket
  nextTick(() => {
    console.log('准备连接到终端:', {
      pod: pod.name,
      namespace: pod.namespace,
      container: currentContainer
    })
    connectTerminal()
    if (terminalInputRef.value) {
      terminalInputRef.value.focus()
    }
  })
}

// 连接终端WebSocket
const connectTerminal = () => {
  if (!currentPod.value) return
  
  const wsUrl = getPodTerminalUrl(
    currentPod.value.name,
    currentPod.value.namespace,
    currentContainer
  )
  
  // 关闭现有连接
  if (ws) {
    ws.close()
  }
  
  try {
    console.log('Connecting to WebSocket:', wsUrl)
    ws = new WebSocket(wsUrl)
    
    ws.onopen = () => {
      console.log('WebSocket connection opened')
      appendToTerminal('连接到终端成功\n')
    }
    
    ws.onmessage = (event) => {
      appendToTerminal(event.data)
    }
    
    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      // 提供更详细的错误信息
      if (error instanceof Event) {
        appendToTerminal(`终端连接错误: 网络错误或连接被拒绝\n`)
      } else {
        appendToTerminal(`终端连接错误: ${JSON.stringify(error)}\n`)
      }
    }
    
    ws.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason)
      // 对于正常的关闭状态码，不显示错误信息给用户
      if (event.code === 1000 || event.code === 1005) {
        // 正常关闭，不显示任何消息给用户
        console.log('WebSocket connection closed normally')
      } else if (event.code === 1006) {
        appendToTerminal('终端连接异常关闭: 可能是网络问题或服务器未响应\n')
      } else {
        appendToTerminal(`终端连接已关闭 (代码: ${event.code})\n`)
      }
    }
  } catch (error: any) {
    console.error('Failed to create WebSocket connection:', error)
    appendToTerminal(`连接终端失败: ${error?.message || JSON.stringify(error)}\n`)
  }
}

// ANSI转义序列到CSS样式的映射
const ansiToCss = (codes: string): string => {
  const codeArray = codes.split(';').map(Number);
  const styles: string[] = [];
  
  for (let i = 0; i < codeArray.length; i++) {
    const code = codeArray[i];
    switch (code) {
      case 0: // 重置
        return '</span>';
      case 1: // 粗体
        styles.push('font-weight: bold');
        break;
      case 30: // 黑色
        styles.push('color: #000000');
        break;
      case 31: // 红色
        styles.push('color: #ff6b6b');
        break;
      case 32: // 绿色
        styles.push('color: #51cf66');
        break;
      case 33: // 黄色
        styles.push('color: #ffd43b');
        break;
      case 34: // 蓝色
        styles.push('color: #4fc1ff');
        break;
      case 35: // 洋红色
        styles.push('color: #cc5de8');
        break;
      case 36: // 青色
        styles.push('color: #22b8cf');
        break;
      case 37: // 白色
        styles.push('color: #f8f9fa');
        break;
      case 40: // 背景黑色
        styles.push('background-color: #000000');
        break;
      case 41: // 背景红色
        styles.push('background-color: #ff6b6b');
        break;
      case 42: // 背景绿色
        styles.push('background-color: #51cf66');
        break;
      case 43: // 背景黄色
        styles.push('background-color: #ffd43b');
        break;
      case 44: // 背景蓝色
        styles.push('background-color: #4fc1ff');
        break;
      case 45: // 背景洋红色
        styles.push('background-color: #cc5de8');
        break;
      case 46: // 背景青色
        styles.push('background-color: #22b8cf');
        break;
      case 47: // 背景白色
        styles.push('background-color: #f8f9fa');
        break;
    }
  }
  
  return `<span style="${styles.join('; ')}">`;
};

// 向终端输出追加内容
const appendToTerminal = (text: string) => {
  if (terminalRef.value) {
    // 创建一个临时元素来解析ANSI转义序列
    const tempDiv = document.createElement('div');
    tempDiv.style.whiteSpace = 'pre';
    tempDiv.style.fontFamily = 'monospace';
    
    // 处理ANSI转义序列
    let formattedText = text;
    
    // 处理光标定位请求等控制序列（例如 \x1b[6n）
    formattedText = formattedText.replace(/\x1b\[\?(\d+)([a-zA-Z])/g, ''); // 查询设备状态等
    formattedText = formattedText.replace(/\x1b\[(\d+)([A-Za-z])/g, ''); // 光标移动等控制序列
    
    // 处理颜色和格式控制序列
    formattedText = formattedText.replace(/\x1b\[([0-9;]*)m/g, (match, codes) => {
      if (codes === '') {
        // 只有ESC[m的情况，默认为重置
        return '</span>';
      }
      return ansiToCss(codes);
    });
    
    // 确保所有打开的标签都被关闭
    let openTags = (formattedText.match(/<span/g) || []).length;
    let closeTags = (formattedText.match(/<\/span/g) || []).length;
    for (let i = 0; i < openTags - closeTags; i++) {
      formattedText += '</span>';
    }
    
    tempDiv.innerHTML = formattedText;
    
    // 将处理后的内容添加到终端
    terminalRef.value.appendChild(tempDiv);
    
    // 滚动到底部
    terminalRef.value.scrollTop = terminalRef.value.scrollHeight;
  }
}

// 发送命令到终端
const sendCommand = () => {
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    ElMessage.error('终端连接未建立')
    return
  }
  
  const command = terminalInput.value.trim()
  if (!command) return
  
  // 添加到命令历史
  if (commandHistory.value[commandHistory.value.length - 1] !== command) {
    commandHistory.value.push(command)
  }
  historyIndex.value = commandHistory.value.length
  
  // 发送到WebSocket
  ws.send(command + '\n')
  
  // 清空输入框
  terminalInput.value = ''
}

// 关闭终端连接
const closeTerminal = () => {
  if (ws) {
    ws.close()
    ws = null
  }
  if (terminalRef.value) {
    terminalRef.value.innerHTML = ''
  }
}

// 处理上箭头键（命令历史）
const handleUpArrow = (event: KeyboardEvent) => {
  if (commandHistory.value.length === 0) return
  
  if (historyIndex.value === -1 || historyIndex.value > commandHistory.value.length - 1) {
    historyIndex.value = commandHistory.value.length
  }
  
  if (historyIndex.value > 0) {
    historyIndex.value--
    terminalInput.value = commandHistory.value[historyIndex.value]
  }
  
  event.preventDefault()
}

// 处理下箭头键（命令历史）
const handleDownArrow = (event: KeyboardEvent) => {
  if (commandHistory.value.length === 0) return
  
  if (historyIndex.value < commandHistory.value.length - 1) {
    historyIndex.value++
    terminalInput.value = commandHistory.value[historyIndex.value]
  } else {
    historyIndex.value = commandHistory.value.length
    terminalInput.value = ''
  }
  
  event.preventDefault()
}

const showCreateDialog = () => {
  loadPodTemplate()
  createDialogVisible.value = true
}

const loadPodTemplate = () => {
  podYaml.value = `apiVersion: v1
kind: Pod
metadata:
  name: my-pod
  namespace: ${namespaceStore.currentNamespace}
  labels:
    app: my-app
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"`
}

const handleCreate = async () => {
  if (!podYaml.value.trim()) {
    ElMessage.warning('请输入 Pod YAML 定义')
    return
  }
  
  creating.value = true
  try {
    await createPodFromYaml(podYaml.value)
    ElMessage.success('创建成功')
    createDialogVisible.value = false
    fetchPods()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || error.message || '创建失败')
  } finally {
    creating.value = false
  }
}

// 监听命名空间变化
watch(() => namespaceStore.currentNamespace, () => {
  fetchPods()
})

// 监听集群变化
watch(() => namespaceStore.currentClusterId, () => {
  // 集群变化时重新获取命名空间和Pod列表
  fetchNamespaces()
  fetchPods()
})

onMounted(() => {
  fetchNamespaces()
  fetchPods()
  
  // 监听集群变更事件
  window.addEventListener('clusterChanged', (event: any) => {
    const cluster = event.detail
    namespaceStore.setCurrentClusterId(cluster.id)
    // 集群变化时重新获取数据
    fetchNamespaces()
    fetchPods()
  })
})
</script>

<style scoped>
.pods-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.terminal-container {
  background-color: #1e1e1e;
  color: #ffffff;
  border-radius: 4px;
  padding: 10px;
  font-family: 'Courier New', monospace;
  height: 400px;
  display: flex;
  flex-direction: column;
}

.terminal-output {
  flex: 1;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
  margin-bottom: 10px;
}

.terminal-input-container {
  display: flex;
  align-items: center;
  border-top: 1px solid #333;
  padding-top: 5px;
}

.prompt {
  color: #4ec9b0;
  margin-right: 5px;
  font-weight: bold;
}

.terminal-input {
  flex: 1;
  background: transparent;
  border: none;
  color: #ffffff;
  font-family: 'Courier New', monospace;
  outline: none;
}

.terminal-input::placeholder {
  color: #666;
}
</style>