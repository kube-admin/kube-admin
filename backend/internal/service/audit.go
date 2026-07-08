package service

import (
	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/model"
)

// AuditService 审计日志服务
type AuditService struct{}

// NewAuditService 创建审计服务实例
func NewAuditService() *AuditService { return &AuditService{} }

// ListAuditLogs 分页查询审计日志（按时间倒序）
func (s *AuditService) ListAuditLogs(page, pageSize int) ([]model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	var total int64
	if err := database.DB.Model(&model.AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []model.AuditLog
	err := database.DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&logs).Error
	return logs, total, err
}
