package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/internal/api"
	"github.com/kube-admin/kube-admin/backend/internal/middleware"
	"github.com/kube-admin/kube-admin/backend/internal/service"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
)

// SetupRouter 设置路由
func SetupRouter(defaultK8sClient *k8s.Client, k8sManager *k8s.Manager) *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	// 健康检查端点（无需认证，供 k8s liveness/readiness 探针使用）
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 创建服务层
	clusterService := service.NewClusterService()
	userService := service.NewUserService()
	auditService := service.NewAuditService()

	// 创建API层
	authAPI := api.NewAuthAPI(userService)
	clusterAPI := api.NewClusterAPI(clusterService)
	userAPI := api.NewUserAPI(userService)
	auditAPI := api.NewAuditAPI(auditService)
	eventAPI := api.NewEventAPI()
	resourceAPI := api.NewResourceAPI()

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authAPI.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.AuditMiddleware())
	{
		// 用户信息（所有登录用户可访问）
		protected.GET("/auth/user", authAPI.GetUserInfo)

		// 用户管理（仅 admin）
		adminGroup := protected.Group("")
		adminGroup.Use(middleware.RequireRole("admin"))
		{
			adminGroup.GET("/users", userAPI.ListUsers)
			adminGroup.GET("/users/:id", userAPI.GetUser)
			adminGroup.POST("/users", userAPI.CreateUser)
			adminGroup.PUT("/users/:id", userAPI.UpdateUser)
			adminGroup.DELETE("/users/:id", userAPI.DeleteUser)

			// 集群管理（仅 admin）
			adminGroup.GET("/clusters", clusterAPI.ListClusters)
			adminGroup.GET("/clusters/:id", clusterAPI.GetCluster)
			adminGroup.POST("/clusters", clusterAPI.CreateCluster)
			adminGroup.PUT("/clusters/:id", clusterAPI.UpdateCluster)
			adminGroup.DELETE("/clusters/:id", clusterAPI.DeleteCluster)
			adminGroup.POST("/clusters/test-connection", clusterAPI.TestConnection)
			adminGroup.POST("/clusters/:id/test-connection", clusterAPI.TestConnectionByID)

			// 审计日志查询（仅 admin）
			adminGroup.GET("/audit/logs", auditAPI.ListAuditLogs)
		}

		// 创建需要集群参数的API组
		k8sGroup := protected.Group("")
		k8sGroup.Use(middleware.WriteAuth()) // 写操作需 admin/operator/user 角色
		k8sGroup.Use(middleware.ClusterMiddleware(defaultK8sClient, k8sManager))
		{
			// Dashboard
			dashboardAPI := api.NewDashboardAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/dashboard/stats", dashboardAPI.GetStats)

			// Event
			k8sGroup.GET("/events", eventAPI.ListEvents)

			// 通用资源管理（任意 GVR：list/get/delete/apply/patch）
			k8sGroup.GET("/resources", resourceAPI.List)
			k8sGroup.GET("/resources/:name", resourceAPI.Get)
			k8sGroup.DELETE("/resources/:name", resourceAPI.Delete)
			k8sGroup.POST("/resources/apply", resourceAPI.Apply)
			k8sGroup.PATCH("/resources/:name", resourceAPI.Patch)
			// 通用 workload 扩缩容/滚动重启（Deployment/StatefulSet/DaemonSet/ReplicaSet）
			k8sGroup.PUT("/resources/:name/scale", resourceAPI.ScaleResource)
			k8sGroup.PUT("/resources/:name/restart", resourceAPI.RestartResource)

			// Namespace
			namespaceAPI := api.NewNamespaceAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/namespaces", namespaceAPI.ListNamespaces)
			k8sGroup.POST("/namespaces", namespaceAPI.CreateNamespace)
			k8sGroup.DELETE("/namespaces/:name", namespaceAPI.DeleteNamespace)
			// 强制清理 finalizers（解除 namespace Terminating 卡死）
			k8sGroup.POST("/namespaces/:name/finalize", namespaceAPI.FinalizeNamespace)
			// APIService 诊断（namespace DiscoveryFailed 卡死排查）
			k8sGroup.GET("/apiservices/unavailable", namespaceAPI.ListUnavailableAPIServices)
			k8sGroup.DELETE("/apiservices/:name", namespaceAPI.DeleteAPIService)

			// Node
			nodeAPI := api.NewNodeAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/nodes", nodeAPI.ListNodes)
			k8sGroup.GET("/nodes/:name", nodeAPI.GetNode)

			// Pod
			podAPI := api.NewPodAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/pods", podAPI.ListPods)
			k8sGroup.GET("/pods/:name", podAPI.GetPod)
			k8sGroup.DELETE("/pods/:name", podAPI.DeletePod)
			k8sGroup.GET("/pods/:name/logs", podAPI.GetPodLogs)
			k8sGroup.GET("/pods/:name/logs/stream", podAPI.LogsStream)
			k8sGroup.POST("/pods/:name/exec", podAPI.ExecCommand)
			k8sGroup.GET("/pods/:name/terminal", podAPI.ExecTerminal)
			k8sGroup.POST("/pods/yaml", podAPI.CreatePodFromYaml)

			// Deployment
			deploymentAPI := api.NewDeploymentAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/deployments", deploymentAPI.ListDeployments)
			k8sGroup.GET("/deployments/:name", deploymentAPI.GetDeployment)
			k8sGroup.DELETE("/deployments/:name", deploymentAPI.DeleteDeployment)
			k8sGroup.PUT("/deployments/:name/scale", deploymentAPI.ScaleDeployment)
			k8sGroup.PUT("/deployments/:name/restart", deploymentAPI.RestartDeployment)
			k8sGroup.POST("/deployments/yaml", deploymentAPI.CreateDeploymentFromYaml)

			// Service
			serviceAPI := api.NewServiceAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/services", serviceAPI.ListServices)
			k8sGroup.GET("/services/:name", serviceAPI.GetService)
			k8sGroup.DELETE("/services/:name", serviceAPI.DeleteService)
			k8sGroup.POST("/services/yaml", serviceAPI.CreateServiceFromYaml)

			// ConfigMap
			configMapAPI := api.NewConfigMapAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/configmaps", configMapAPI.ListConfigMaps)
			k8sGroup.GET("/configmaps/:name", configMapAPI.GetConfigMap)
			k8sGroup.POST("/configmaps", configMapAPI.CreateConfigMap)
			k8sGroup.PUT("/configmaps/:name", configMapAPI.UpdateConfigMap)
			k8sGroup.DELETE("/configmaps/:name", configMapAPI.DeleteConfigMap)

			// Secret
			secretAPI := api.NewSecretAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/secrets", secretAPI.ListSecrets)
			k8sGroup.GET("/secrets/:name", secretAPI.GetSecret)
			k8sGroup.POST("/secrets", secretAPI.CreateSecret)
			k8sGroup.PUT("/secrets/:name", secretAPI.UpdateSecret)
			k8sGroup.DELETE("/secrets/:name", secretAPI.DeleteSecret)
		}
	}

	return r
}
