package service

import (
	"context"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeService Node服务
type NodeService struct {
	k8sClient *k8s.Client
}

// NewNodeService 创建Node服务
func NewNodeService(k8sClient *k8s.Client) *NodeService {
	return &NodeService{k8sClient: k8sClient}
}

// ListNodes 获取Node列表
func (s *NodeService) ListNodes() ([]model.NodeInfo, error) {
	nodeList, err := s.k8sClient.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodes []model.NodeInfo
	for _, node := range nodeList.Items {
		nodes = append(nodes, s.convertNode(&node))
	}

	return nodes, nil
}

// GetNode 获取Node详情
func (s *NodeService) GetNode(name string) (*model.NodeInfo, error) {
	node, err := s.k8sClient.ClientSet.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	nodeInfo := s.convertNode(node)
	return &nodeInfo, nil
}

// convertNode 转换Node对象
func (s *NodeService) convertNode(node *corev1.Node) model.NodeInfo {
	nodeInfo := model.NodeInfo{
		K8sResource: model.K8sResource{
			Name:              node.Name,
			Labels:            node.Labels,
			Annotations:       node.Annotations,
			CreationTimestamp: node.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   node.ResourceVersion,
		},
		Status:           getNodeStatus(node),
		OSImage:          node.Status.NodeInfo.OSImage,
		KubeletVersion:   node.Status.NodeInfo.KubeletVersion,
		ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
	}

	// 获取内部IP
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			nodeInfo.InternalIP = addr.Address
			break
		}
	}

	// 资源容量
	if cpu, ok := node.Status.Capacity[corev1.ResourceCPU]; ok {
		nodeInfo.Capacity.CPU = cpu.String()
	}
	if memory, ok := node.Status.Capacity[corev1.ResourceMemory]; ok {
		nodeInfo.Capacity.Memory = memory.String()
	}
	if pods, ok := node.Status.Capacity[corev1.ResourcePods]; ok {
		nodeInfo.Capacity.Pods = pods.String()
	}

	// 可分配资源
	if cpu, ok := node.Status.Allocatable[corev1.ResourceCPU]; ok {
		nodeInfo.Allocatable.CPU = cpu.String()
	}
	if memory, ok := node.Status.Allocatable[corev1.ResourceMemory]; ok {
		nodeInfo.Allocatable.Memory = memory.String()
	}
	if pods, ok := node.Status.Allocatable[corev1.ResourcePods]; ok {
		nodeInfo.Allocatable.Pods = pods.String()
	}

	// 节点条件
	for _, condition := range node.Status.Conditions {
		nodeInfo.Conditions = append(nodeInfo.Conditions, model.NodeCondition{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}

	return nodeInfo
}

// getNodeStatus 获取节点状态
func getNodeStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}
