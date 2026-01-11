package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// DashboardAPI Dashboard API
type DashboardAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	dashboardService *service.DashboardService
}

// NewDashboardAPI 创建Dashboard API
func NewDashboardAPI(dashboardService *service.DashboardService) *DashboardAPI {
	return &DashboardAPI{dashboardService: dashboardService}
}

// GetStats 获取集群统计信息
func (a *DashboardAPI) GetStats(c *gin.Context) {
	// 从上下文中获取服务实例
	dashboardService, exists := c.Get("dashboard_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	stats, err := dashboardService.(*service.DashboardService).GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(stats))
}
