package model

import (
	"time"
)

// Cluster 集群信息模型
type Cluster struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"uniqueIndex;not null"`
	Description   string    `json:"description"`
	ServerURL     string    `json:"server_url"`
	Token         string    `json:"token"`
	ConfigPath    string    `json:"config_path"`
	ConfigContent string    `json:"config_content"`                 // 新增：配置文件内容
	Status        string    `json:"status" gorm:"default:'active'"` // active, inactive, error
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ClusterRequest 创建/更新集群请求
type ClusterRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	ServerURL     string `json:"server_url"` // 不再强制要求
	Token         string `json:"token"`      // 不再强制要求
	ConfigPath    string `json:"config_path"`
	ConfigContent string `json:"config_content"` // 新增：配置文件内容
}

// ClusterResponse 集群响应
type ClusterResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	ServerURL     string    `json:"server_url"`
	ConfigPath    string    `json:"config_path"`
	ConfigContent string    `json:"config_content"` // 新增：配置文件内容
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TestConnectionRequest 测试连接请求
type TestConnectionRequest struct {
	ServerURL     string `json:"server_url"` // 不再强制要求
	Token         string `json:"token"`      // 不再强制要求
	ConfigPath    string `json:"config_path"`
	ConfigContent string `json:"config_content"` // 新增：配置文件内容
}

// TestConnectionResponse 测试连接响应
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}
