package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
	"path/filepath"
)

// Client K8s客户端封装
type Client struct {
	ClientSet        *kubernetes.Clientset
	MetricsClientSet *versioned.Clientset
	Config           *rest.Config
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
