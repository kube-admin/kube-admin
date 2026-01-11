package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// ClusterAPI 集群API控制器
type ClusterAPI struct {
	clusterService *service.ClusterService
}

// NewClusterAPI 创建集群API实例
func NewClusterAPI(clusterService *service.ClusterService) *ClusterAPI {
	return &ClusterAPI{clusterService: clusterService}
}

// ListClusters 获取集群列表
func (a *ClusterAPI) ListClusters(c *gin.Context) {
	clusters, err := a.clusterService.ListClusters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(clusters))
}

// GetCluster 获取集群详情
func (a *ClusterAPI) GetCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "Invalid cluster ID"))
		return
	}

	cluster, err := a.clusterService.GetCluster(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(404, "Cluster not found"))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(cluster))
}

// CreateCluster 创建集群
func (a *ClusterAPI) CreateCluster(c *gin.Context) {
	var req model.ClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	// 验证必填字段
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "Cluster name is required"))
		return
	}

	// 验证至少提供了一种连接方式
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token"))
		return
	}

	cluster, err := a.clusterService.CreateCluster(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse(cluster))
}

// UpdateCluster 更新集群
func (a *ClusterAPI) UpdateCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "Invalid cluster ID"))
		return
	}

	var req model.ClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	// 验证必填字段
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "Cluster name is required"))
		return
	}

	// 验证至少提供了一种连接方式
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token"))
		return
	}

	cluster, err := a.clusterService.UpdateCluster(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(404, "Cluster not found"))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(cluster))
}

// DeleteCluster 删除集群
func (a *ClusterAPI) DeleteCluster(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "Invalid cluster ID"))
		return
	}

	if err := a.clusterService.DeleteCluster(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"message": "Cluster deleted successfully"}))
}

// TestConnection 测试集群连接
func (a *ClusterAPI) TestConnection(c *gin.Context) {
	var req model.TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, err.Error()))
		return
	}

	// 验证至少提供了一种连接方式
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token"))
		return
	}

	result, err := a.clusterService.TestConnection(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(result))
}
