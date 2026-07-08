import { describe, it, expect, beforeEach } from 'vitest'
import { currentUserRole } from './permission'

describe('currentUserRole', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('未登录时返回空字符串', () => {
    expect(currentUserRole()).toBe('')
  })

  it('正确返回已登录用户的角色', () => {
    localStorage.setItem('user', JSON.stringify({ username: 'admin', role: 'admin' }))
    expect(currentUserRole()).toBe('admin')
  })

  it('用户对象无 role 字段时返回空', () => {
    localStorage.setItem('user', JSON.stringify({ username: 'foo' }))
    expect(currentUserRole()).toBe('')
  })

  it('localStorage 存在非法 JSON 时返回空且不抛异常', () => {
    localStorage.setItem('user', '{not valid json')
    expect(currentUserRole()).toBe('')
  })
})
