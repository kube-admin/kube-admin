package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// PodAPI Pod API
type PodAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	podService *service.PodService
}

// NewPodAPI 创建Pod API
func NewPodAPI(podService *service.PodService) *PodAPI {
	return &PodAPI{podService: podService}
}

// ListPods 获取Pod列表
func (a *PodAPI) ListPods(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.DefaultQuery("namespace", "default")

	pods, err := podService.(*service.PodService).ListPods(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(pods))
}

// GetPod 获取Pod详情
func (a *PodAPI) GetPod(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	pod, err := podService.(*service.PodService).GetPod(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(pod))
}

// DeletePod 删除Pod
func (a *PodAPI) DeletePod(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	err := podService.(*service.PodService).DeletePod(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// GetPodLogs 获取Pod日志
func (a *PodAPI) GetPodLogs(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")
	container := c.Query("container")
	tailLines := c.DefaultQuery("tail_lines", "100")

	lines, _ := strconv.ParseInt(tailLines, 10, 64)
	logs, err := podService.(*service.PodService).GetPodLogs(namespace, name, container, lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(map[string]string{"logs": logs}))
}

// ExecCommand 在Pod中执行命令
func (a *PodAPI) ExecCommand(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	// 获取请求参数
	namespace := c.Query("namespace")
	podName := c.Param("name")
	container := c.Query("container")

	// 获取命令参数
	var req struct {
		Command []string `json:"command"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的请求参数: "+err.Error()))
		return
	}

	if len(req.Command) == 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "命令不能为空"))
		return
	}

	// 执行命令
	stdout, stderr, err := podService.(*service.PodService).ExecCommand(namespace, podName, container, req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "执行命令失败: "+err.Error()))
		return
	}

	// 返回结果
	result := map[string]interface{}{
		"stdout":  stdout,
		"stderr":  stderr,
		"success": stderr == "",
	}

	c.JSON(http.StatusOK, model.SuccessResponse(result))
}

// ExecTerminal WebSocket终端接口
func (a *PodAPI) ExecTerminal(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	// 升级到WebSocket连接
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许跨域
		},
	}

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "升级到WebSocket失败: "+err.Error()))
		return
	}
	defer wsConn.Close()

	// 获取请求参数
	namespace := c.Query("namespace")
	podName := c.Param("name")
	container := c.Query("container")

	// 执行终端连接
	err = podService.(*service.PodService).ExecTerminal(namespace, podName, container, wsConn)
	if err != nil {
		// 只有在连接仍然活跃时才发送错误信息
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("WebSocket error: %v", err)
			wsConn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		}
		return
	}
}

// CreatePodFromYaml 通过YAML创建Pod
func (a *PodAPI) CreatePodFromYaml(c *gin.Context) {
	// 从上下文中获取服务实例
	podService, exists := c.Get("pod_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	var req struct {
		YAML string `json:"yaml"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的请求参数"))
		return
	}

	if req.YAML == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "YAML内容不能为空"))
		return
	}

	ps := podService.(*service.PodService)

	// 创建dynamic client
	dynamicClient, err := dynamic.NewForConfig(ps.GetK8sClient().Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "创建Dynamic Client失败: "+err.Error()))
		return
	}

	// 创建discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(ps.GetK8sClient().Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "创建Discovery Client失败: "+err.Error()))
		return
	}

	// 创建cached discovery client
	cachedDiscoveryClient := memory.NewMemCacheClient(discoveryClient)

	// 创建rest mapper
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)

	// 解析YAML
	decode := yamlutil.NewYAMLOrJSONDecoder(strings.NewReader(req.YAML), 4096)
	for {
		obj := &unstructured.Unstructured{}
		err := decode.Decode(obj)
		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "解析YAML失败: "+err.Error()))
			return
		}

		if len(obj.Object) == 0 {
			continue
		}

		// 获取GVK
		gvk := obj.GroupVersionKind()

		// 获取mapping
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "获取REST Mapping失败: "+err.Error()))
			return
		}

		// 获取namespace
		namespace, _, err := unstructured.NestedString(obj.Object, "metadata", "namespace")
		if err != nil || namespace == "" {
			namespace = "default"
		}

		// 创建资源
		var dr dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			dr = dynamicClient.Resource(mapping.Resource).Namespace(namespace)
		} else {
			dr = dynamicClient.Resource(mapping.Resource)
		}

		_, err = dr.Create(context.TODO(), obj, metav1.CreateOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "创建资源失败: "+err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}
