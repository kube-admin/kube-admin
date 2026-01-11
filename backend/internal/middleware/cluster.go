package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/service"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
)

// ClusterMiddleware 集群中间件，根据请求参数获取对应的K8s客户端
func ClusterMiddleware(defaultK8sClient *k8s.Client, k8sManager *k8s.Manager) gin.HandlerFunc {
	// 初始化数据库中的集群服务
	clusterService := service.NewClusterService()

	return func(c *gin.Context) {
		// 从查询参数或表单参数中获取集群ID
		clusterIDStr := c.Query("cluster_id")
		if clusterIDStr == "" {
			clusterIDStr = c.PostForm("cluster_id")
		}

		// 如果没有指定集群ID，使用默认客户端
		if clusterIDStr == "" {
			// 将默认客户端注入到上下文中
			c.Set("k8s_client", defaultK8sClient)

			// 创建服务实例并注入到上下文中
			podService := service.NewPodService(defaultK8sClient)
			deploymentService := service.NewDeploymentService(defaultK8sClient)
			serviceService := service.NewServiceService(defaultK8sClient)
			namespaceService := service.NewNamespaceService(defaultK8sClient)
			nodeService := service.NewNodeService(defaultK8sClient)
			configMapService := service.NewConfigMapService(defaultK8sClient)
			secretService := service.NewSecretService(defaultK8sClient)
			dashboardService := service.NewDashboardService(defaultK8sClient)

			c.Set("pod_service", podService)
			c.Set("deployment_service", deploymentService)
			c.Set("service_service", serviceService)
			c.Set("namespace_service", namespaceService)
			c.Set("node_service", nodeService)
			c.Set("configmap_service", configMapService)
			c.Set("secret_service", secretService)
			c.Set("dashboard_service", dashboardService)

			c.Next()
			return
		}

		// 解析集群ID
		clusterID, err := strconv.ParseUint(clusterIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse(400, "无效的集群ID"))
			c.Abort()
			return
		}

		// 从数据库获取集群信息
		cluster, err := clusterService.GetCluster(uint(clusterID))
		if err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse(404, "集群不存在"))
			c.Abort()
			return
		}

		// 获取集群对应的K8s客户端
		k8sClient, err := k8sManager.GetClient(uint(clusterID), cluster)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse(500, fmt.Sprintf("无法连接到集群: %v", err)))
			c.Abort()
			return
		}

		// 将客户端和服务注入到上下文中
		c.Set("k8s_client", k8sClient)

		podService := service.NewPodService(k8sClient)
		deploymentService := service.NewDeploymentService(k8sClient)
		serviceService := service.NewServiceService(k8sClient)
		namespaceService := service.NewNamespaceService(k8sClient)
		nodeService := service.NewNodeService(k8sClient)
		configMapService := service.NewConfigMapService(k8sClient)
		secretService := service.NewSecretService(k8sClient)
		dashboardService := service.NewDashboardService(k8sClient)

		c.Set("pod_service", podService)
		c.Set("deployment_service", deploymentService)
		c.Set("service_service", serviceService)
		c.Set("namespace_service", namespaceService)
		c.Set("node_service", nodeService)
		c.Set("configmap_service", configMapService)
		c.Set("secret_service", secretService)
		c.Set("dashboard_service", dashboardService)

		c.Next()
	}
}
