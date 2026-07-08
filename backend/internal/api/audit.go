package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// AuditAPI 审计日志API
type AuditAPI struct {
	auditService *service.AuditService
}

// NewAuditAPI 创建审计API实例
func NewAuditAPI(auditService *service.AuditService) *AuditAPI {
	return &AuditAPI{auditService: auditService}
}

// ListAuditLogs 分页查询审计日志
func (a *AuditAPI) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	logs, total, err := a.auditService.ListAuditLogs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(model.PageResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Items:    logs,
	}))
}
