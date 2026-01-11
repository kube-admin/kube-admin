package api

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

// ServiceAPI Service API处理器
type ServiceAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	serviceService *service.ServiceService
}

// NewServiceAPI 创建Service API
func NewServiceAPI(serviceService *service.ServiceService) *ServiceAPI {
	return &ServiceAPI{serviceService: serviceService}
}

// ListServices 获取Service列表
func (a *ServiceAPI) ListServices(c *gin.Context) {
	// 从上下文中获取服务实例
	serviceService, exists := c.Get("service_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.DefaultQuery("namespace", "default")

	services, err := serviceService.(*service.ServiceService).ListServices(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(services))
}

// GetService 获取Service详情
func (a *ServiceAPI) GetService(c *gin.Context) {
	// 从上下文中获取服务实例
	serviceService, exists := c.Get("service_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	name := c.Param("name")
	namespace := c.DefaultQuery("namespace", "default")

	svc, err := serviceService.(*service.ServiceService).GetService(namespace, name)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(svc))
}

// DeleteService 删除Service
func (a *ServiceAPI) DeleteService(c *gin.Context) {
	// 从上下文中获取服务实例
	serviceService, exists := c.Get("service_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	name := c.Param("name")
	namespace := c.DefaultQuery("namespace", "default")

	err := serviceService.(*service.ServiceService).DeleteService(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(gin.H{"message": "Service deleted successfully"}))
}

// CreateServiceFromYaml 通过YAML创建Service
func (a *ServiceAPI) CreateServiceFromYaml(c *gin.Context) {
	// 从上下文中获取服务实例
	serviceService, exists := c.Get("service_service")
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

	ss := serviceService.(*service.ServiceService)

	// 创建dynamic client
	dynamicClient, err := dynamic.NewForConfig(ss.GetK8sClient().Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "创建Dynamic Client失败: "+err.Error()))
		return
	}

	// 创建discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(ss.GetK8sClient().Config)
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
