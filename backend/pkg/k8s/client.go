package k8s

import (
	"os"
	"path/filepath"

	"github.com/kube-admin/kube-admin/backend/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// Client K8s客户端封装
type Client struct {
	ClientSet        *kubernetes.Clientset
	MetricsClientSet *versioned.Clientset
	Config           *rest.Config
}

// applyConfigDefaults 统一为 rest.Config 注入全局默认值（请求超时）。
// 集群不可达时让请求按 K8S_REQUEST_TIMEOUT 快速失败，而非 client-go 默认挂起 30s。
func applyConfigDefaults(cfg *rest.Config) *rest.Config {
	if cfg == nil {
		return cfg
	}
	cfg.Timeout = config.App.K8sTimeout
	return cfg
}

// NewClient 创建K8s客户端
func NewClient(kubeconfigPath string) (*Client, error) {
	var config *rest.Config
	var err error

	if kubeconfigPath == "" {
		// 尝试使用默认路径
		if home := os.Getenv("HOME"); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	// 首先尝试使用kubeconfig文件
	if _, err := os.Stat(kubeconfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果kubeconfig不存在,尝试使用集群内配置
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	applyConfigDefaults(config)
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	metricsClientSet, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		ClientSet:        clientSet,
		MetricsClientSet: metricsClientSet,
		Config:           config,
	}, nil
}
