package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
)

// EventAPI K8s 事件API
type EventAPI struct{}

// NewEventAPI 创建事件API实例
func NewEventAPI() *EventAPI { return &EventAPI{} }

// ListEvents 查询事件。支持 namespace 过滤，支持按 kind+name 过滤某资源的事件。
func (a *EventAPI) ListEvents(c *gin.Context) {
	eventService, exists := c.Get("event_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace") // 空表示所有命名空间

	// 支持按关联资源过滤：involvedObject.kind + involvedObject.name
	kind := c.Query("kind")
	name := c.Query("name")
	var fieldSelector string
	if kind != "" || name != "" {
		parts := make([]string, 0, 2)
		if kind != "" {
			parts = append(parts, "involvedObject.kind="+kind)
		}
		if name != "" {
			parts = append(parts, "involvedObject.name="+name)
		}
		fieldSelector = strings.Join(parts, ",")
	}

	events, err := eventService.(*service.EventService).ListEvents(namespace, fieldSelector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(events))
}
