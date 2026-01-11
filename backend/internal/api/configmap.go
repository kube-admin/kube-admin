package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// ConfigMapAPI ConfigMap API
type ConfigMapAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	configMapService *service.ConfigMapService
}

// NewConfigMapAPI 创建ConfigMap API
func NewConfigMapAPI(configMapService *service.ConfigMapService) *ConfigMapAPI {
	return &ConfigMapAPI{configMapService: configMapService}
}

// ListConfigMaps 获取ConfigMap列表
func (a *ConfigMapAPI) ListConfigMaps(c *gin.Context) {
	// 从上下文中获取服务实例
	configMapService, exists := c.Get("configmap_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.DefaultQuery("namespace", "default")

	configMaps, err := configMapService.(*service.ConfigMapService).ListConfigMaps(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(configMaps))
}

// GetConfigMap 获取ConfigMap详情
func (a *ConfigMapAPI) GetConfigMap(c *gin.Context) {
	// 从上下文中获取服务实例
	configMapService, exists := c.Get("configmap_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	configMap, err := configMapService.(*service.ConfigMapService).GetConfigMap(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(configMap))
}

// CreateConfigMap 创建ConfigMap
func (a *ConfigMapAPI) CreateConfigMap(c *gin.Context) {
	// 从上下文中获取服务实例
	configMapService, exists := c.Get("configmap_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	var req struct {
		Namespace string            `json:"namespace" binding:"required"`
		Name      string            `json:"name" binding:"required"`
		Data      map[string]string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	err := configMapService.(*service.ConfigMapService).CreateConfigMap(req.Namespace, req.Name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// UpdateConfigMap 更新ConfigMap
func (a *ConfigMapAPI) UpdateConfigMap(c *gin.Context) {
	// 从上下文中获取服务实例
	configMapService, exists := c.Get("configmap_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	var req struct {
		Data map[string]string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	err := configMapService.(*service.ConfigMapService).UpdateConfigMap(namespace, name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteConfigMap 删除ConfigMap
func (a *ConfigMapAPI) DeleteConfigMap(c *gin.Context) {
	// 从上下文中获取服务实例
	configMapService, exists := c.Get("configmap_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	err := configMapService.(*service.ConfigMapService).DeleteConfigMap(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
