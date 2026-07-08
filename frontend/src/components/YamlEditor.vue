<template>
  <div class="yaml-editor" :style="{ height }">
    <VueMonacoEditor
      v-model:value="code"
      language="yaml"
      :theme="theme"
      :options="editorOptions"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { VueMonacoEditor, loader } from '@guolao/vue-monaco-editor'
import * as monaco from 'monaco-editor'
import { useDark } from '@vueuse/core'

// 初始化 monaco loader（模块级，仅执行一次）
loader.config({ monaco })

const props = withDefaults(defineProps<{
  modelValue: string
  readonly?: boolean
  height?: string
}>(), {
  readonly: false,
  height: '400px'
})

const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const code = ref(props.modelValue)

// 外部值更新同步到内部
watch(() => props.modelValue, (v) => {
  if (v !== code.value) code.value = v
})
// 内部编辑同步到外部
watch(code, (v) => emit('update:modelValue', v))

const isDark = useDark()
const theme = computed(() => (isDark.value ? 'vs-dark' : 'vs'))

// 使用字面量断言以匹配 monaco 的字面量联合类型
const editorOptions = computed(() => ({
  readOnly: props.readonly,
  minimap: { enabled: false },
  fontSize: 13,
  tabSize: 2,
  automaticLayout: true,
  scrollBeyondLastLine: false,
  wordWrap: 'on' as const,
  renderWhitespace: 'selection' as const
}))
</script>

<style scoped>
.yaml-editor {
  width: 100%;
  border: 1px solid var(--el-border-color, #dcdfe6);
  border-radius: 4px;
  overflow: hidden;
}
</style>
