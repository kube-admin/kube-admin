package main

import (
	"fmt"
	"log"

	"github.com/kube-admin/kube-admin/backend/config"
	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/router"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
)

func main() {
	// 初始化数据库
	database.InitDB()

	// 加载配置
	cfg := config.LoadConfig()

	// 创建K8s客户端管理器
	k8sManager := k8s.NewManager()

	// 创建默认的K8s客户端（用于初始连接）
	defaultK8sClient, err := k8s.NewClient(cfg.KubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to create default k8s client: %v", err)
	}

	log.Println("Successfully connected to default Kubernetes cluster")

	// 设置路由
	r := router.SetupRouter(defaultK8sClient, k8sManager)

	// 启动服务器
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
