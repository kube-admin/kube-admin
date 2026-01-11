package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// SecretAPI Secret API
type SecretAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	secretService *service.SecretService
}

// NewSecretAPI 创建Secret API
func NewSecretAPI(secretService *service.SecretService) *SecretAPI {
	return &SecretAPI{secretService: secretService}
}

// ListSecrets 获取Secret列表
func (a *SecretAPI) ListSecrets(c *gin.Context) {
	// 从上下文中获取服务实例
	secretService, exists := c.Get("secret_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.DefaultQuery("namespace", "default")

	secrets, err := secretService.(*service.SecretService).ListSecrets(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(secrets))
}

// GetSecret 获取Secret详情
func (a *SecretAPI) GetSecret(c *gin.Context) {
	// 从上下文中获取服务实例
	secretService, exists := c.Get("secret_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")
	decode := c.DefaultQuery("decode", "false") == "true"

	secret, err := secretService.(*service.SecretService).GetSecret(namespace, name, decode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(secret))
}

// CreateSecret 创建Secret
func (a *SecretAPI) CreateSecret(c *gin.Context) {
	// 从上下文中获取服务实例
	secretService, exists := c.Get("secret_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	var req struct {
		Namespace string            `json:"namespace" binding:"required"`
		Name      string            `json:"name" binding:"required"`
		Type      string            `json:"type"`
		Data      map[string]string `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	if req.Type == "" {
		req.Type = "Opaque"
	}

	err := secretService.(*service.SecretService).CreateSecret(req.Namespace, req.Name, req.Type, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// UpdateSecret 更新Secret
func (a *SecretAPI) UpdateSecret(c *gin.Context) {
	// 从上下文中获取服务实例
	secretService, exists := c.Get("secret_service")
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

	err := secretService.(*service.SecretService).UpdateSecret(namespace, name, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// DeleteSecret 删除Secret
func (a *SecretAPI) DeleteSecret(c *gin.Context) {
	// 从上下文中获取服务实例
	secretService, exists := c.Get("secret_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	err := secretService.(*service.SecretService).DeleteSecret(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
