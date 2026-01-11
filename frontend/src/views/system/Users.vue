<template>
  <div class="users-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <div>
            <el-button type="primary" @click="handleAddUser">添加用户</el-button>
          </div>
        </div>
      </template>

      <el-table :data="users" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="180" />
        <el-table-column prop="email" label="邮箱" width="250" />
        <el-table-column prop="role" label="角色" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.role === 'admin' ? 'success' : 'info'">
              {{ scope.row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="200" />
        <el-table-column prop="updated_at" label="更新时间" width="200" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="scope">
            <el-button size="small" type="primary" @click="handleEditUser(scope.row)">
              编辑
            </el-button>
            <el-popconfirm
              title="确定删除这个用户吗?"
              @confirm="handleDeleteUser(scope.row)"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
    >
      <el-form ref="userFormRef" :model="formData" label-position="top" :rules="formRules">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="formData.role" placeholder="请选择角色">
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!formData.id" label="密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item v-if="formData.id" label="新密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="留空则不修改密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>


  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, FormInstance, FormRules } from 'element-plus'
import { User } from '@/apis/user/login'
import { getUsers as fetchUsers, createUser, updateUser, deleteUser } from '@/apis/user/users'

// 表格数据
const users = ref<User[]>([])
const loading = ref(false)

// 对话框状态
const dialogVisible = ref(false)
const dialogTitle = ref('添加用户')

// 表单数据
const formData = reactive({
  id: 0,
  username: '',
  email: '',
  role: 'user',
  password: ''
})

// 表单引用
const userFormRef = ref<FormInstance>()

// 密码验证
const validatePassword = (rule: any, value: string, callback: any) => {
  // 编辑时密码可以为空
  if (!formData.id && !value) {
    callback(new Error('请输入密码'))
  } else if (value && value.length < 6) {
    callback(new Error('密码长度不能少于 6 个字符'))
  } else {
    callback()
  }
}

// 表单验证规则
const formRules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  password: [
    { validator: validatePassword, trigger: 'blur' }
  ]
})

// 获取用户列表
const getUsers = async () => {
  loading.value = true
  try {
    const response = await fetchUsers()
    users.value = response.data.data || []
  } catch (error: any) {
    ElMessage.error(error?.message || '获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 处理添加用户
const handleAddUser = () => {
  dialogTitle.value = '添加用户'
  formData.id = 0
  formData.username = ''
  formData.email = ''
  formData.role = 'user'
  formData.password = ''
  dialogVisible.value = true
}

// 处理编辑用户
const handleEditUser = (row: User) => {
  dialogTitle.value = '编辑用户'
  formData.id = row.id
  formData.username = row.username
  formData.email = row.email
  formData.role = row.role
  formData.password = ''
  dialogVisible.value = true
}

// 处理表单提交
const handleSubmit = async () => {
  if (!userFormRef.value) return
  
  try {
    await userFormRef.value.validate()
    
    if (formData.id) {
      // 更新用户
      await updateUser(formData.id, {
        username: formData.username,
        email: formData.email,
        role: formData.role,
        password: formData.password || undefined
      })
      ElMessage.success('更新用户成功')
    } else {
      // 创建用户
      await createUser({
        username: formData.username,
        email: formData.email,
        role: formData.role,
        password: formData.password
      })
      ElMessage.success('创建用户成功')
    }
    
    dialogVisible.value = false
    getUsers()
  } catch (error: any) {
    ElMessage.error(error?.message || '操作失败')
  }
}

// 处理删除用户
const handleDeleteUser = async (row: User) => {
  try {
    await deleteUser(row.id)
    ElMessage.success('删除用户成功')
    getUsers()
  } catch (error: any) {
    ElMessage.error(error?.message || '删除失败')
  }
}

// 初始化
onMounted(() => {
  getUsers()
})
</script>

<style scoped>
.users-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
