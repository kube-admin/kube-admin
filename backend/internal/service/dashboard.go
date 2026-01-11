package service

import (
	"context"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DashboardService Dashboard统计服务
type DashboardService struct {
	k8sClient *k8s.Client
}

// NewDashboardService 创建Dashboard服务
func NewDashboardService(k8sClient *k8s.Client) *DashboardService {
	return &DashboardService{k8sClient: k8sClient}
}

// DashboardStats Dashboard统计数据
type DashboardStats struct {
	NodeCount       int            `json:"node_count"`
	NamespaceCount  int            `json:"namespace_count"`
	PodCount        int            `json:"pod_count"`
	DeploymentCount int            `json:"deployment_count"`
	ServiceCount    int            `json:"service_count"`
	ConfigMapCount  int            `json:"configmap_count"`
	SecretCount     int            `json:"secret_count"`
	PodStatusStats  PodStatusStats `json:"pod_status_stats"`
}

// PodStatusStats Pod状态统计
type PodStatusStats struct {
	Running   int `json:"running"`
	Pending   int `json:"pending"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
	Unknown   int `json:"unknown"`
}

// GetDashboardStats 获取Dashboard统计数据
func (s *DashboardService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 统计节点数
	nodeList, err := s.k8sClient.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.NodeCount = len(nodeList.Items)
	}

	// 统计命名空间数
	nsList, err := s.k8sClient.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.NamespaceCount = len(nsList.Items)
	}

	// 统计所有命名空间的资源
	podList, err := s.k8sClient.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.PodCount = len(podList.Items)
		// 统计Pod状态
		for _, pod := range podList.Items {
			switch pod.Status.Phase {
			case "Running":
				stats.PodStatusStats.Running++
			case "Pending":
				stats.PodStatusStats.Pending++
			case "Failed":
				stats.PodStatusStats.Failed++
			case "Succeeded":
				stats.PodStatusStats.Succeeded++
			default:
				stats.PodStatusStats.Unknown++
			}
		}
	}

	deployList, err := s.k8sClient.ClientSet.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.DeploymentCount = len(deployList.Items)
	}

	svcList, err := s.k8sClient.ClientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.ServiceCount = len(svcList.Items)
	}

	cmList, err := s.k8sClient.ClientSet.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.ConfigMapCount = len(cmList.Items)
	}

	secretList, err := s.k8sClient.ClientSet.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.SecretCount = len(secretList.Items)
	}

	return stats, nil
}
