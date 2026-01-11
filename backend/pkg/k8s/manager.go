package k8s

import (
	"fmt"
	"sync"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	versioned "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Manager 多集群管理器
type Manager struct {
	clusters map[uint]*Client
	mutex    sync.RWMutex
}

// NewManager 创建多集群管理器
func NewManager() *Manager {
	return &Manager{
		clusters: make(map[uint]*Client),
	}
}

// GetClient 获取指定集群的客户端
func (m *Manager) GetClient(clusterID uint, cluster *model.Cluster) (*Client, error) {
	m.mutex.RLock()
	client, exists := m.clusters[clusterID]
	m.mutex.RUnlock()

	if exists {
		return client, nil
	}

	// 创建新的客户端
	newClient, err := m.createClient(cluster)
	if err != nil {
		return nil, err
	}

	// 存储客户端
	m.mutex.Lock()
	m.clusters[clusterID] = newClient
	m.mutex.Unlock()

	return newClient, nil
}

// createClient 根据集群信息创建K8s客户端
func (m *Manager) createClient(cluster *model.Cluster) (*Client, error) {
	var config *rest.Config
	var err error

	// 优先使用配置内容
	if cluster.ConfigContent != "" {
		clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(cluster.ConfigContent))
		if err != nil {
			return nil, fmt.Errorf("failed to build config from content: %v", err)
		}

		config, err = clientConfig.ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get client config: %v", err)
		}
	} else if cluster.ConfigPath != "" {
		// 使用配置文件路径
		config, err = clientcmd.BuildConfigFromFlags("", cluster.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from file: %v", err)
		}
	} else if cluster.ServerURL != "" && cluster.Token != "" {
		// 使用URL和Token
		config = &rest.Config{
			Host:        cluster.ServerURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true,
			},
		}
	} else {
		return nil, fmt.Errorf("no valid configuration found for cluster")
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	metricsClientSet, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	return &Client{
		ClientSet:        clientSet,
		MetricsClientSet: metricsClientSet,
		Config:           config,
	}, nil
}

// RemoveClient 移除指定集群的客户端
func (m *Manager) RemoveClient(clusterID uint) {
	m.mutex.Lock()
	delete(m.clusters, clusterID)
	m.mutex.Unlock()
}
