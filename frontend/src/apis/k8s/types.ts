// Cluster 集群信息（脱敏响应，不含 Token 与 ConfigContent 明文）
export interface Cluster {
  id: number
  name: string
  description: string
  server_url: string
  config_path: string
  has_config_content: boolean
  has_token: boolean
  status: string
  created_at: string
  updated_at: string
}

// ClusterRequest 创建/更新集群请求。更新时 token/config_content 留空表示不修改。
export interface ClusterRequest {
  name: string
  description: string
  server_url: string
  token: string
  config_path: string
  config_content: string
}

// TestConnectionRequest 测试连接请求（明文，用于未保存集群的预测试）
export interface TestConnectionRequest {
  server_url: string
  token: string
  config_path: string
  config_content: string
}

// TestConnectionResponse 测试连接响应
export interface TestConnectionResponse {
  success: boolean
  message: string
  version: string
}
