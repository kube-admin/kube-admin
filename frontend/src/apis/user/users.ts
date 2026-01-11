import request from '@/apis/client/request'
import type { User } from './login'

export interface CreateUserRequest {
  username: string
  email: string
  role: string
  password: string
}

export interface UpdateUserRequest {
  username?: string
  email?: string
  role?: string
  password?: string
}

// 获取用户列表
export const getUsers = () => {
  return request.get<User[]>('/api/v1/users')
}

// 创建用户
export const createUser = (data: CreateUserRequest) => {
  return request.post<User>('/api/v1/users', data)
}

// 更新用户
export const updateUser = (id: number, data: UpdateUserRequest) => {
  return request.put<User>(`/api/v1/users/${id}`, data)
}

// 删除用户
export const deleteUser = (id: number) => {
  return request.delete(`/api/v1/users/${id}`)
}
