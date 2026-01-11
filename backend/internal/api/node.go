package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// NodeAPI Node API
type NodeAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	nodeService *service.NodeService
}

// NewNodeAPI 创建Node API
func NewNodeAPI(nodeService *service.NodeService) *NodeAPI {
	return &NodeAPI{nodeService: nodeService}
}

// ListNodes 获取Node列表
func (a *NodeAPI) ListNodes(c *gin.Context) {
	// 从上下文中获取服务实例
	nodeService, exists := c.Get("node_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	nodes, err := nodeService.(*service.NodeService).ListNodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nodes))
}

// GetNode 获取Node详情
func (a *NodeAPI) GetNode(c *gin.Context) {
	// 从上下文中获取服务实例
	nodeService, exists := c.Get("node_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	name := c.Param("name")

	node, err := nodeService.(*service.NodeService).GetNode(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(node))
}


