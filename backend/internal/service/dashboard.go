package service

import (
	"context"

	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	ClusterUsage    ClusterUsage   `json:"cluster_usage"` // 集群资源实时使用率
}

// PodStatusStats Pod状态统计
type PodStatusStats struct {
	Running   int `json:"running"`
	Pending   int `json:"pending"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
	Unknown   int `json:"unknown"`
}

// ClusterUsage 集群整体资源使用率（基于 node metrics 聚合）
type ClusterUsage struct {
	CPUPercent     float64 `json:"cpu_percent"`
	MemoryPercent  float64 `json:"memory_percent"`
	CPUUsed        string  `json:"cpu_used"`
	CPUCapacity    string  `json:"cpu_capacity"`
	MemoryUsed     string  `json:"memory_used"`
	MemoryCapacity string  `json:"memory_capacity"`
}

// GetDashboardStats 获取Dashboard统计数据
func (s *DashboardService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 节点列表（复用：计数 + 使用率聚合）
	nodeList, err := s.k8sClient.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.NodeCount = len(nodeList.Items)
		s.fillClusterUsage(stats, nodeList)
	}

	// 命名空间数
	nsList, err := s.k8sClient.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.NamespaceCount = len(nsList.Items)
	}

	// 统计所有命名空间的资源
	podList, err := s.k8sClient.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err == nil {
		stats.PodCount = len(podList.Items)
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

// fillClusterUsage 聚合所有节点的实时资源使用率。metrics-server 不可用时使用率为 0。
func (s *DashboardService) fillClusterUsage(stats *DashboardStats, nodeList *corev1.NodeList) {
	metricsMap := nodeMetricsMap(s.k8sClient)

	var usedCPU, allocCPU, usedMem, allocMem resource.Quantity
	for i := range nodeList.Items {
		node := &nodeList.Items[i]
		allocCPU.Add(node.Status.Allocatable[corev1.ResourceCPU])
		allocMem.Add(node.Status.Allocatable[corev1.ResourceMemory])
		if usage, ok := metricsMap[node.Name]; ok {
			usedCPU.Add(usage[corev1.ResourceCPU])
			usedMem.Add(usage[corev1.ResourceMemory])
		}
	}

	stats.ClusterUsage = ClusterUsage{
		CPUPercent:     calcPercent(usedCPU, allocCPU),
		MemoryPercent:  calcPercent(usedMem, allocMem),
		CPUUsed:        usedCPU.String(),
		CPUCapacity:    allocCPU.String(),
		MemoryUsed:     usedMem.String(),
		MemoryCapacity: allocMem.String(),
	}
}
