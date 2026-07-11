<template>
  <el-drawer v-model="visible" :title="title" size="65%" destroy-on-close>
    <YamlEditor v-model="yamlContent" height="75vh" />
    <div class="drawer-footer">
      <el-button @click="visible = false">关闭</el-button>
      <el-button type="primary" :loading="saving" @click="save">应用保存</el-button>
    </div>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import yaml from 'js-yaml'
import YamlEditor from './YamlEditor.vue'
import { getResource, applyResource, type GVR } from '@/apis/k8s'

// 通用 YAML 查看/编辑抽屉：任意 GVR 资源，查看用 getResource 转 YAML，保存用 applyResource（server-side apply）。
// 各列表页通过 ref 调用 open(gvr, namespace, name)，统一"查看/编辑"操作，告别每页各自实现。
const emit = defineEmits<{ saved: [] }>()

const visible = ref(false)
const title = ref('YAML')
const yamlContent = ref('')
const saving = ref(false)
const current = ref<{ gvr: GVR; namespace: string; name: string } | null>(null)

const open = async (gvr: GVR, namespace: string, name: string) => {
  current.value = { gvr, namespace, name }
  title.value = `YAML · ${name}`
  visible.value = true
  try {
    const res: any = await getResource(gvr, namespace, name)
    const obj = res.data?.data ?? res.data
    yamlContent.value = yaml.dump(obj, { noRefs: true, lineWidth: 120 })
  } catch (e: any) {
    ElMessage.error(e?.message || '加载 YAML 失败')
    yamlContent.value = ''
  }
}

const save = async () => {
  if (!current.value) return
  saving.value = true
  try {
    await applyResource(yamlContent.value)
    ElMessage.success('保存成功')
    visible.value = false
    emit('saved')
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

defineExpose({ open })
</script>

<style scoped>
.drawer-footer {
  margin-top: 12px;
  text-align: right;
}
</style>
