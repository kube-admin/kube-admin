package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kube-admin/kube-admin/backend/config"
	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/internal/router"
	"github.com/kube-admin/kube-admin/backend/pkg/crypto"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
)

func main() {
	// 1. 加载配置（最先执行，其余组件依赖配置）
	cfg := config.LoadConfig()
	gin.SetMode(cfg.GinMode)

	// 2. 注入全局密钥：JWT 签名密钥与凭据加密密钥
	model.InitJWTSecret(cfg.JWTSecret)
	if err := crypto.Init(cfg.EncryptKey); err != nil {
		log.Fatalf("Failed to init crypto: %v", err)
	}

	// 3. 初始化数据库（文件持久化）
	database.InitDB(cfg.DBPath)

	// 4. 创建 K8s 客户端管理器与默认客户端
	k8sManager := k8s.NewManager()
	defaultK8sClient, err := k8s.NewClient(cfg.KubeconfigPath)
	if err != nil {
		// 默认集群连接失败不应直接退出：用户可能通过界面添加集群
		log.Printf("[WARN] Failed to create default k8s client: %v (可忽略，通过界面添加集群)", err)
		defaultK8sClient = nil
	} else {
		log.Println("Successfully connected to default Kubernetes cluster")
	}

	// 5. 设置路由（含健康检查）
	r := router.SetupRouter(defaultK8sClient, k8sManager)

	// 6. 启动 HTTP 服务，支持优雅关闭
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 7. 等待中断信号，优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
