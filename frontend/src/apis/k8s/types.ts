// Cluster 集群信息
export interface Cluster {
  id: number
  name: string
  description: string
  server_url: string
  status: string
  created_at: string
  updated_at: string
}

// ClusterRequest 创建/更新集群请求
export interface ClusterRequest {
  name: string
  description: string
  server_url: string
  token: string
  config_path: string
}

// TestConnectionRequest 测试连接请求
export interface TestConnectionRequest {
  server_url: string
  token: string
  config_path: string
}

// TestConnectionResponse 测试连接响应
export interface TestConnectionResponse {
  success: boolean
  message: string
  version: string
}