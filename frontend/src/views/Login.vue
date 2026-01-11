<template>
  <el-row class="login" justify="center" align="middle">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>用户登录</span>
          <div class="dark-icon" @click="toggleDark()">
            <el-icon>
              <Moon v-if="isDark" />
              <Sunny v-else />
            </el-icon>
          </div>
        </div>
      </template>
      <div>
        <el-form
          ref="ruleFormRef"
          :model="ruleForm"
          status-icon
          :rules="rules"
          label-width="60px"
          class="demo-ruleForm"
        >
          <el-form-item label="用户名" prop="username">
            <el-input v-model.number="ruleForm.username" />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="ruleForm.password"
              type="password"
              autocomplete="off"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitForm(ruleFormRef)"
              >登录</el-button
            >
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </el-row>
</template>

<style scoped>
.login {
  position: absolute;
  top: 20%;
  bottom: 60%;
  width: 100%;
}

.box-card {
  width: 450px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dark-icon {
  font-size: 20px;
  cursor: pointer;
}

.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}
</style>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toggleDark, isDark } from '@/stores/dark'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { login } from '@/apis/user/login'

const ruleFormRef = ref<FormInstance>()
const router = useRouter()

const validateEmpty = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('字段不能为空'))
  } else {
    callback()
  }
}

const ruleForm = reactive({
  username: 'admin',
  password: 'admin123'
})

const rules = reactive<FormRules>({
  username: [{ validator: validateEmpty, trigger: 'blur' }],
  password: [{ validator: validateEmpty, trigger: 'blur' }]
})

const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.validate(async (valid) => {
    if (valid) {
      try {
        const response = await login(ruleForm)
        const responseData = response.data as any
        const { token, user } = responseData.data
        
        // 保存 token 到 localStorage
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
        
        ElMessage.success(`欢迎回来, ${user.username}!`)
        
        // 跳转到首页
        router.push('/')
      } catch (error: any) {
        console.error('登录失败:', error)
        ElMessage.error(error.response?.data?.message || '登录失败，请检查用户名和密码')
      }
    } else {
      ElMessage.warning('请填写完整的登录信息')
    }
  })
}
</script>
