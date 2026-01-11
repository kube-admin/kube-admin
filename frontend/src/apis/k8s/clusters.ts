import request from '@/apis/client/request'
import { Cluster, ClusterRequest, TestConnectionRequest, TestConnectionResponse } from './types'

// 获取集群列表
export const listClusters = () => {
  return request.get<Cluster[]>('/api/v1/clusters')
}

// 获取集群详情
export const getCluster = (id: number) => {
  return request.get<Cluster>(`/api/v1/clusters/${id}`)
}

// 创建集群
export const createCluster = (data: ClusterRequest) => {
  return request.post<Cluster>('/api/v1/clusters', data)
}

// 更新集群
export const updateCluster = (id: number, data: ClusterRequest) => {
  return request.put<Cluster>(`/api/v1/clusters/${id}`, data)
}

// 删除集群
export const deleteCluster = (id: number) => {
  return request.delete(`/api/v1/clusters/${id}`)
}

// 测试集群连接
export const testConnection = (data: TestConnectionRequest) => {
  return request.post<TestConnectionResponse>('/api/v1/clusters/test-connection', data)
}