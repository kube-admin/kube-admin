package service

import (
	"fmt"

	"github.com/kube-admin/kube-admin/backend/config"
	"github.com/kube-admin/kube-admin/backend/database"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ClusterService 集群服务
type ClusterService struct{}

// NewClusterService 创建集群服务实例
func NewClusterService() *ClusterService {
	return &ClusterService{}
}

// buildRestConfig 根据集群连接信息构建 rest.Config。
// 优先级：ConfigContent > ConfigPath > ServerURL+Token。DRY：供 Test/GetK8sClient 复用。
func buildRestConfig(configContent, configPath, serverURL, token string) (*rest.Config, error) {
	if configContent != "" {
		clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(configContent))
		if err != nil {
			return nil, fmt.Errorf("failed to build config from content: %v", err)
		}
		return clientConfig.ClientConfig()
	}
	if configPath != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from file: %v", err)
		}
		return cfg, nil
	}
	if serverURL != "" && token != "" {
		return &rest.Config{
			Host:        serverURL,
			BearerToken: token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: config.App.TLSSkipVerify,
			},
		}, nil
	}
	return nil, fmt.Errorf("必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token")
}

// ListClusters 获取集群列表（脱敏）
func (s *ClusterService) ListClusters() ([]model.ClusterResponse, error) {
	var clusters []model.Cluster
	if err := database.DB.Find(&clusters).Error; err != nil {
		return nil, err
	}

	responses := make([]model.ClusterResponse, 0, len(clusters))
	for i := range clusters {
		responses = append(responses, clusters[i].ToResponse())
	}
	return responses, nil
}

// GetCluster 获取集群详情（明文，供内部业务使用）
func (s *ClusterService) GetCluster(id uint) (*model.Cluster, error) {
	var cluster model.Cluster
	if err := database.DB.First(&cluster, id).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

// GetClusterResponse 获取集群脱敏响应（供 API 返回）
func (s *ClusterService) GetClusterResponse(id uint) (*model.ClusterResponse, error) {
	cluster, err := s.GetCluster(id)
	if err != nil {
		return nil, err
	}
	resp := cluster.ToResponse()
	return &resp, nil
}

// CreateCluster 创建集群
func (s *ClusterService) CreateCluster(req model.ClusterRequest) (*model.ClusterResponse, error) {
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		return nil, fmt.Errorf("必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token")
	}

	cluster := model.Cluster{
		Name:          req.Name,
		Description:   req.Description,
		ServerURL:     req.ServerURL,
		Token:         req.Token,
		ConfigPath:    req.ConfigPath,
		ConfigContent: req.ConfigContent,
		Status:        "active",
	}

	if err := database.DB.Create(&cluster).Error; err != nil {
		return nil, err
	}

	// 重新查询以触发 AfterFind 解密，再生成脱敏响应
	var created model.Cluster
	if err := database.DB.First(&created, cluster.ID).Error; err != nil {
		return nil, err
	}
	resp := created.ToResponse()
	return &resp, nil
}

// UpdateCluster 更新集群。Token/ConfigContent 留空表示保留原值。
func (s *ClusterService) UpdateCluster(id uint, req model.ClusterRequest) (*model.ClusterResponse, error) {
	cluster, err := s.GetCluster(id)
	if err != nil {
		return nil, err
	}

	cluster.Name = req.Name
	cluster.Description = req.Description
	cluster.ServerURL = req.ServerURL
	cluster.ConfigPath = req.ConfigPath

	// 仅在提供新值时更新敏感字段，留空表示保留
	if req.Token != "" {
		cluster.Token = req.Token
	}
	if req.ConfigContent != "" {
		cluster.ConfigContent = req.ConfigContent
	}

	// 校验：更新后仍需至少一种可用连接方式
	if cluster.ConfigContent == "" && cluster.ConfigPath == "" && (cluster.ServerURL == "" || cluster.Token == "") {
		return nil, fmt.Errorf("更新后集群无可用连接方式，请保留或重新提供凭据")
	}

	if err := database.DB.Save(cluster).Error; err != nil {
		return nil, err
	}

	// 重新查询获取干净的明文，生成脱敏响应
	var updated model.Cluster
	if err := database.DB.First(&updated, id).Error; err != nil {
		return nil, err
	}
	resp := updated.ToResponse()
	return &resp, nil
}

// DeleteCluster 删除集群
func (s *ClusterService) DeleteCluster(id uint) error {
	return database.DB.Delete(&model.Cluster{}, id).Error
}

// TestConnection 测试连接（基于请求中的明文凭据，用于未保存集群的预测试）
func (s *ClusterService) TestConnection(req model.TestConnectionRequest) (*model.TestConnectionResponse, error) {
	cfg, err := buildRestConfig(req.ConfigContent, req.ConfigPath, req.ServerURL, req.Token)
	if err != nil {
		return &model.TestConnectionResponse{Success: false, Message: err.Error()}, nil
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return &model.TestConnectionResponse{Success: false, Message: fmt.Sprintf("failed to create kubernetes client: %v", err)}, nil
	}

	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return &model.TestConnectionResponse{Success: false, Message: fmt.Sprintf("failed to connect to cluster: %v", err)}, nil
	}

	return &model.TestConnectionResponse{Success: true, Message: "Connection successful", Version: version.GitVersion}, nil
}

// TestConnectionByID 基于已保存集群ID测试连接（用解密后的凭据）
func (s *ClusterService) TestConnectionByID(id uint) (*model.TestConnectionResponse, error) {
	cluster, err := s.GetCluster(id)
	if err != nil {
		return nil, fmt.Errorf("集群不存在: %v", err)
	}
	req := model.TestConnectionRequest{
		ServerURL:     cluster.ServerURL,
		Token:         cluster.Token,
		ConfigPath:    cluster.ConfigPath,
		ConfigContent: cluster.ConfigContent,
	}
	return s.TestConnection(req)
}

// GetK8sClient 获取K8s客户端（明文集群）
func (s *ClusterService) GetK8sClient(clusterID uint) (*kubernetes.Clientset, error) {
	cluster, err := s.GetCluster(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %v", err)
	}

	cfg, err := buildRestConfig(cluster.ConfigContent, cluster.ConfigPath, cluster.ServerURL, cluster.Token)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}
	return clientset, nil
}
