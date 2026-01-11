package api

import (
	"context"
	"io"
	"net/http"
	"strconv"
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

// DeploymentAPI Deployment API
type DeploymentAPI struct {
	// 注意：在多集群环境中，服务实例将在中间件中动态注入
	deploymentService *service.DeploymentService
}

// NewDeploymentAPI 创建Deployment API
func NewDeploymentAPI(deploymentService *service.DeploymentService) *DeploymentAPI {
	return &DeploymentAPI{deploymentService: deploymentService}
}

// ListDeployments 获取Deployment列表
func (a *DeploymentAPI) ListDeployments(c *gin.Context) {
	// 从上下文中获取服务实例
	deploymentService, exists := c.Get("deployment_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.DefaultQuery("namespace", "default")

	deployments, err := deploymentService.(*service.DeploymentService).ListDeployments(namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(deployments))
}

// GetDeployment 获取Deployment详情
func (a *DeploymentAPI) GetDeployment(c *gin.Context) {
	// 从上下文中获取服务实例
	deploymentService, exists := c.Get("deployment_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	deployment, err := deploymentService.(*service.DeploymentService).GetDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(deployment))
}

// DeleteDeployment 删除Deployment
func (a *DeploymentAPI) DeleteDeployment(c *gin.Context) {
	// 从上下文中获取服务实例
	deploymentService, exists := c.Get("deployment_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	err := deploymentService.(*service.DeploymentService).DeleteDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// ScaleDeployment 扩缩容Deployment
func (a *DeploymentAPI) ScaleDeployment(c *gin.Context) {
	// 从上下文中获取服务实例
	deploymentService, exists := c.Get("deployment_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")
	replicasStr := c.Query("replicas")

	replicas, err := strconv.ParseInt(replicasStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "replicas参数错误"))
		return
	}

	err = deploymentService.(*service.DeploymentService).ScaleDeployment(namespace, name, int32(replicas))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// RestartDeployment 重启Deployment
func (a *DeploymentAPI) RestartDeployment(c *gin.Context) {
	// 从上下文中获取服务实例
	deploymentService, exists := c.Get("deployment_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "服务未初始化"))
		return
	}

	namespace := c.Query("namespace")
	name := c.Param("name")

	err := deploymentService.(*service.DeploymentService).RestartDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(nil))
}

// CreateDeploymentFromYaml 通过YAML创建Deployment
func (a *DeploymentAPI) CreateDeploymentFromYaml(c *gin.Context) {
	// 从上下文中获取服务实例
	deploymentService, exists := c.Get("deployment_service")
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

	ds := deploymentService.(*service.DeploymentService)

	// 创建dynamic client
	dynamicClient, err := dynamic.NewForConfig(ds.GetK8sClient().Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, "创建Dynamic Client失败: "+err.Error()))
		return
	}

	// 创建discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(ds.GetK8sClient().Config)
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
