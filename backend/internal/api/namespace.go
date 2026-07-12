package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// NamespaceAPI Namespace API
type NamespaceAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	namespaceService *service.NamespaceService
}

// NewNamespaceAPI 创建Namespace API
func NewNamespaceAPI(namespaceService *service.NamespaceService) *NamespaceAPI {
	return &NamespaceAPI{namespaceService: namespaceService}
}

// ListNamespaces 获取Namespace列表
func (a *NamespaceAPI) ListNamespaces(c *gin.Context) {
	// 从上下文中获取服务实例
	namespaceService, exists := c.Get("namespace_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespaces, err := namespaceService.(*service.NamespaceService).ListNamespaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(namespaces))
}

// CreateNamespace 创建Namespace
func (a *NamespaceAPI) CreateNamespace(c *gin.Context) {
	// 从上下文中获取服务实例
	namespaceService, exists := c.Get("namespace_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	err := namespaceService.(*service.NamespaceService).CreateNamespace(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteNamespace 删除Namespace
func (a *NamespaceAPI) DeleteNamespace(c *gin.Context) {
	// 从上下文中获取服务实例
	namespaceService, exists := c.Get("namespace_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	name := c.Param("name")

	err := namespaceService.(*service.NamespaceService).DeleteNamespace(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// FinalizeNamespace 强制清理 finalizers（解除 Terminating 卡死）
func (a *NamespaceAPI) FinalizeNamespace(c *gin.Context) {
	namespaceService, exists := c.Get("namespace_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	name := c.Param("name")
	if err := namespaceService.(*service.NamespaceService).FinalizeNamespace(name); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// ListUnavailableAPIServices 列出失效 APIService（诊断 namespace Terminating 用）
func (a *NamespaceAPI) ListUnavailableAPIServices(c *gin.Context) {
	namespaceService, exists := c.Get("namespace_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	list, err := namespaceService.(*service.NamespaceService).ListUnavailableAPIServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(list))
}

// DeleteAPIService 删除指定 APIService（集群级危险操作）
func (a *NamespaceAPI) DeleteAPIService(c *gin.Context) {
	namespaceService, exists := c.Get("namespace_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	name := c.Param("name")
	if err := namespaceService.(*service.NamespaceService).DeleteAPIService(name); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
