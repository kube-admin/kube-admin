import request from '@/apis/client/request'

export interface User {
  id: number
  username: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: number
  user: User
}

// 登录
export const login = (data: LoginRequest) => {
  return request.post<LoginResponse>('/api/v1/auth/login', data)
}

// 获取用户信息
export const getUserInfo = () => {
  return request.get<User>('/api/v1/auth/user')
}
