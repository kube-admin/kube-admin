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

	// 创建服务层
	clusterService := service.NewClusterService()
	userService := service.NewUserService()

	// 创建API层
	authAPI := api.NewAuthAPI(userService)
	clusterAPI := api.NewClusterAPI(clusterService)
	userAPI := api.NewUserAPI(userService)

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authAPI.Login)
	}

	// 需要认证的路由
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户信息
		protected.GET("/auth/user", authAPI.GetUserInfo)

		// 用户管理
		protected.GET("/users", userAPI.ListUsers)
		protected.GET("/users/:id", userAPI.GetUser)
		protected.POST("/users", userAPI.CreateUser)
		protected.PUT("/users/:id", userAPI.UpdateUser)
		protected.DELETE("/users/:id", userAPI.DeleteUser)

		// 集群管理
		protected.GET("/clusters", clusterAPI.ListClusters)
		protected.GET("/clusters/:id", clusterAPI.GetCluster)
		protected.POST("/clusters", clusterAPI.CreateCluster)
		protected.PUT("/clusters/:id", clusterAPI.UpdateCluster)
		protected.DELETE("/clusters/:id", clusterAPI.DeleteCluster)
		protected.POST("/clusters/test-connection", clusterAPI.TestConnection)

		// 创建需要集群参数的API组
		k8sGroup := protected.Group("")
		k8sGroup.Use(middleware.ClusterMiddleware(defaultK8sClient, k8sManager))
		{
			// Dashboard
			dashboardAPI := api.NewDashboardAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/dashboard/stats", dashboardAPI.GetStats)

			// Namespace
			namespaceAPI := api.NewNamespaceAPI(nil) // 将在中间件中注入正确的客户端
			k8sGroup.GET("/namespaces", namespaceAPI.ListNamespaces)
			k8sGroup.POST("/namespaces", namespaceAPI.CreateNamespace)
			k8sGroup.DELETE("/namespaces/:name", namespaceAPI.DeleteNamespace)

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
