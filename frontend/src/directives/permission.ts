import type { Directive } from 'vue'

// 从 localStorage 获取当前登录用户角色
export function currentUserRole(): string {
  try {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    return user.role || ''
  } catch {
    return ''
  }
}

// v-permission 指令：按角色控制元素显隐，无权限则从 DOM 移除。
// 用法：v-permission="'admin'" 或 v-permission="['admin','operator']"
export const permission: Directive<HTMLElement, string | string[]> = {
  mounted(el, binding) {
    const role = currentUserRole()
    const allowed = Array.isArray(binding.value) ? binding.value : [binding.value]
    if (!allowed.includes(role)) {
      el.parentNode?.removeChild(el)
    }
  }
}
