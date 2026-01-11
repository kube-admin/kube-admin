package service

import (
	"fmt"

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

// ListClusters 获取集群列表
func (s *ClusterService) ListClusters() ([]model.ClusterResponse, error) {
	var clusters []model.Cluster
	if err := database.DB.Find(&clusters).Error; err != nil {
		return nil, err
	}

	var responses []model.ClusterResponse
	for _, cluster := range clusters {
		responses = append(responses, model.ClusterResponse{
			ID:            cluster.ID,
			Name:          cluster.Name,
			Description:   cluster.Description,
			ServerURL:     cluster.ServerURL,
			ConfigPath:    cluster.ConfigPath,
			ConfigContent: cluster.ConfigContent, // 新增：配置文件内容
			Status:        cluster.Status,
			CreatedAt:     cluster.CreatedAt,
			UpdatedAt:     cluster.UpdatedAt,
		})
	}

	return responses, nil
}

// GetCluster 获取集群详情
func (s *ClusterService) GetCluster(id uint) (*model.Cluster, error) {
	var cluster model.Cluster
	if err := database.DB.First(&cluster, id).Error; err != nil {
		return nil, err
	}
	return &cluster, nil
}

// CreateCluster 创建集群
func (s *ClusterService) CreateCluster(req model.ClusterRequest) (*model.ClusterResponse, error) {
	// 验证至少提供了一种连接方式
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		return nil, fmt.Errorf("必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token")
	}

	cluster := model.Cluster{
		Name:          req.Name,
		Description:   req.Description,
		ServerURL:     req.ServerURL,
		Token:         req.Token,
		ConfigPath:    req.ConfigPath,
		ConfigContent: req.ConfigContent, // 新增：配置文件内容
		Status:        "active",
	}

	if err := database.DB.Create(&cluster).Error; err != nil {
		return nil, err
	}

	response := &model.ClusterResponse{
		ID:            cluster.ID,
		Name:          cluster.Name,
		Description:   cluster.Description,
		ServerURL:     cluster.ServerURL,
		ConfigPath:    cluster.ConfigPath,
		ConfigContent: cluster.ConfigContent, // 新增：配置文件内容
		Status:        cluster.Status,
		CreatedAt:     cluster.CreatedAt,
		UpdatedAt:     cluster.UpdatedAt,
	}

	return response, nil
}

// UpdateCluster 更新集群
func (s *ClusterService) UpdateCluster(id uint, req model.ClusterRequest) (*model.ClusterResponse, error) {
	var cluster model.Cluster
	if err := database.DB.First(&cluster, id).Error; err != nil {
		return nil, err
	}

	// 验证至少提供了一种连接方式
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		return nil, fmt.Errorf("必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token")
	}

	cluster.Name = req.Name
	cluster.Description = req.Description
	cluster.ServerURL = req.ServerURL
	cluster.Token = req.Token
	cluster.ConfigPath = req.ConfigPath
	cluster.ConfigContent = req.ConfigContent // 新增：配置文件内容

	if err := database.DB.Save(&cluster).Error; err != nil {
		return nil, err
	}

	response := &model.ClusterResponse{
		ID:            cluster.ID,
		Name:          cluster.Name,
		Description:   cluster.Description,
		ServerURL:     cluster.ServerURL,
		ConfigPath:    cluster.ConfigPath,
		ConfigContent: cluster.ConfigContent, // 新增：配置文件内容
		Status:        cluster.Status,
		CreatedAt:     cluster.CreatedAt,
		UpdatedAt:     cluster.UpdatedAt,
	}

	return response, nil
}

// DeleteCluster 删除集群
func (s *ClusterService) DeleteCluster(id uint) error {
	return database.DB.Delete(&model.Cluster{}, id).Error
}

// TestConnection 测试集群连接
func (s *ClusterService) TestConnection(req model.TestConnectionRequest) (*model.TestConnectionResponse, error) {
	var config *rest.Config
	var err error

	// 验证至少提供了一种连接方式
	if req.ConfigContent == "" && req.ConfigPath == "" && (req.ServerURL == "" || req.Token == "") {
		return &model.TestConnectionResponse{
			Success: false,
			Message: "必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token",
		}, nil
	}

	// 如果提供了配置内容，优先使用配置内容
	if req.ConfigContent != "" {
		clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(req.ConfigContent))
		if err != nil {
			return &model.TestConnectionResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to build config from content: %v", err),
			}, nil
		}

		config, err = clientConfig.ClientConfig()
		if err != nil {
			return &model.TestConnectionResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to get client config: %v", err),
			}, nil
		}
	} else if req.ConfigPath != "" {
		// 如果提供了config path，使用config文件
		config, err = clientcmd.BuildConfigFromFlags("", req.ConfigPath)
		if err != nil {
			return &model.TestConnectionResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to build config from file: %v", err),
			}, nil
		}
	} else {
		// 否则使用Token方式
		config = &rest.Config{
			Host:        req.ServerURL,
			BearerToken: req.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 在生产环境中应该设置为false并提供证书
			},
		}
	}

	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &model.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create kubernetes client: %v", err),
		}, nil
	}

	// 测试连接
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return &model.TestConnectionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to connect to kubernetes cluster: %v", err),
		}, nil
	}

	return &model.TestConnectionResponse{
		Success: true,
		Message: "Connection successful",
		Version: version.GitVersion,
	}, nil
}

// GetK8sClient 获取K8s客户端
func (s *ClusterService) GetK8sClient(clusterID uint) (*kubernetes.Clientset, error) {
	cluster, err := s.GetCluster(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %v", err)
	}

	// 验证至少提供了一种连接方式
	if cluster.ConfigContent == "" && cluster.ConfigPath == "" && (cluster.ServerURL == "" || cluster.Token == "") {
		return nil, fmt.Errorf("集群配置不完整，必须提供至少一种连接方式：1. kubeconfig内容 2. kubeconfig文件路径 3. 服务器地址和Token")
	}

	var config *rest.Config
	// 如果提供了配置内容，优先使用配置内容
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
		// 使用config文件
		config, err = clientcmd.BuildConfigFromFlags("", cluster.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from file: %v", err)
		}
	} else {
		// 使用Token方式
		config = &rest.Config{
			Host:        cluster.ServerURL,
			BearerToken: cluster.Token,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: true, // 在生产环境中应该设置为false并提供证书
			},
		}
	}

	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return clientset, nil
}
