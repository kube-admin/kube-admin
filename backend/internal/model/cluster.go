package model

import (
	"time"

	"github.com/kube-admin/kube-admin/backend/pkg/crypto"
	"gorm.io/gorm"
)

// Cluster 集群信息模型。
// Token/ConfigContent 在写入数据库前由 BeforeSave 钩子加密，
// 读取时由 AfterFind 钩子解密，业务层始终操作明文。
type Cluster struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"uniqueIndex;not null"`
	Description   string    `json:"description"`
	ServerURL     string    `json:"server_url"`
	Token         string    `json:"-" gorm:"column:token"`           // 加密存储，不序列化输出
	ConfigPath    string    `json:"config_path"`
	ConfigContent string    `json:"-" gorm:"column:config_content"`  // 加密存储，不序列化输出
	Status        string    `json:"status" gorm:"default:'active'"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BeforeSave 写入前加密敏感字段
func (c *Cluster) BeforeSave(tx *gorm.DB) error {
	if c.Token != "" {
		enc, err := crypto.Encrypt(c.Token)
		if err != nil {
			return err
		}
		c.Token = enc
	}
	if c.ConfigContent != "" {
		enc, err := crypto.Encrypt(c.ConfigContent)
		if err != nil {
			return err
		}
		c.ConfigContent = enc
	}
	return nil
}

// AfterFind 读取后解密敏感字段；解密失败则保留原值，兼容历史明文数据平滑迁移
func (c *Cluster) AfterFind(tx *gorm.DB) error {
	if c.Token != "" {
		if dec, err := crypto.Decrypt(c.Token); err == nil {
			c.Token = dec
		}
	}
	if c.ConfigContent != "" {
		if dec, err := crypto.Decrypt(c.ConfigContent); err == nil {
			c.ConfigContent = dec
		}
	}
	return nil
}

// ToResponse 将 Cluster 转为脱敏响应
func (c *Cluster) ToResponse() ClusterResponse {
	return ClusterResponse{
		ID:               c.ID,
		Name:             c.Name,
		Description:      c.Description,
		ServerURL:        c.ServerURL,
		ConfigPath:       c.ConfigPath,
		HasConfigContent: c.ConfigContent != "",
		HasToken:         c.Token != "",
		Status:           c.Status,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

// ClusterRequest 创建/更新集群请求。更新时 Token/ConfigContent 留空表示不修改。
type ClusterRequest struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	ServerURL     string `json:"server_url"`
	Token         string `json:"token"`
	ConfigPath    string `json:"config_path"`
	ConfigContent string `json:"config_content"`
}

// ClusterResponse 集群响应（脱敏，不含 Token 与 ConfigContent 明文）
type ClusterResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ServerURL        string    `json:"server_url"`
	ConfigPath       string    `json:"config_path"`
	HasConfigContent bool      `json:"has_config_content"`
	HasToken         bool      `json:"has_token"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TestConnectionRequest 测试连接请求（明文，用于未保存的连接测试）
type TestConnectionRequest struct {
	ServerURL     string `json:"server_url"`
	Token         string `json:"token"`
	ConfigPath    string `json:"config_path"`
	ConfigContent string `json:"config_content"`
}

// TestConnectionResponse 测试连接响应
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}
