package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// ResourceAPI 通用资源API，支持任意 GVR 资源的增删改查与 YAML apply
type ResourceAPI struct{}

// NewResourceAPI 创建通用资源API实例
func NewResourceAPI() *ResourceAPI { return &ResourceAPI{} }

// applyRequest apply 请求体
type applyRequest struct {
	YAML string `json:"yaml" binding:"required"`
}

// parseGVR 从查询参数解析 GVR 与 namespace
func parseGVR(c *gin.Context) (schema.GroupVersionResource, string) {
	return schema.GroupVersionResource{
		Group:    c.Query("group"),
		Version:  c.Query("version"),
		Resource: c.Query("resource"),
	}, c.Query("namespace")
}

// validateGVR 校验 GVR 必填
func validateGVR(gvr schema.GroupVersionResource) bool {
	return gvr.Resource != "" && gvr.Version != ""
}

// List 通用列表
func (a *ResourceAPI) List(c *gin.Context) {
	rs, exists := c.Get("resource_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	gvr, ns := parseGVR(c)
	if !validateGVR(gvr) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "version 和 resource 参数必填"))
		return
	}
	list, err := rs.(*service.ResourceService).List(gvr, ns)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(list))
}

// Get 通用获取
func (a *ResourceAPI) Get(c *gin.Context) {
	rs, exists := c.Get("resource_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	gvr, ns := parseGVR(c)
	if !validateGVR(gvr) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "version 和 resource 参数必填"))
		return
	}
	obj, err := rs.(*service.ResourceService).Get(gvr, ns, c.Param("name"))
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(404, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(obj))
}

// Delete 通用删除
func (a *ResourceAPI) Delete(c *gin.Context) {
	rs, exists := c.Get("resource_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	gvr, ns := parseGVR(c)
	if !validateGVR(gvr) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "version 和 resource 参数必填"))
		return
	}
	if err := rs.(*service.ResourceService).Delete(gvr, ns, c.Param("name")); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"message": "deleted"}))
}

// Apply 从 YAML 创建或更新任意资源
func (a *ResourceAPI) Apply(c *gin.Context) {
	rs, exists := c.Get("resource_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	var req applyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}
	obj, err := rs.(*service.ResourceService).ApplyFromYAML(req.YAML)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(obj))
}

// Patch 通用补丁
func (a *ResourceAPI) Patch(c *gin.Context) {
	rs, exists := c.Get("resource_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}
	gvr, ns := parseGVR(c)
	if !validateGVR(gvr) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "version 和 resource 参数必填"))
		return
	}
	var req struct {
		PatchType string `json:"patch_type"`
		Data      string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}
	pt := types.PatchType(req.PatchType)
	if pt == "" {
		pt = types.StrategicMergePatchType
	}
	obj, err := rs.(*service.ResourceService).Patch(gvr, ns, c.Param("name"), pt, []byte(req.Data))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.SuccessResponse(obj))
}
