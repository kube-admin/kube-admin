package service

import (
	"context"
	"fmt"
	"time"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceService Namespace服务
type NamespaceService struct {
	k8sClient *k8s.Client
}

// NewNamespaceService 创建Namespace服务
func NewNamespaceService(k8sClient *k8s.Client) *NamespaceService {
	return &NamespaceService{k8sClient: k8sClient}
}

// formatAge 格式化年龄显示
func formatAge(creationTime time.Time) string {
	duration := time.Since(creationTime)

	// 如果小于1分钟，显示秒数
	if duration < time.Minute {
		return "< 1m"
	}

	// 如果小于1小时，显示分钟数
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%dm", minutes)
	}

	// 如果小于1天，显示小时数
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh", hours)
	}

	// 如果小于7天，显示天数
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}

	// 如果小于30天，显示周数
	if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / (24 * 7))
		return fmt.Sprintf("%dw", weeks)
	}

	// 如果小于365天，显示月数
	if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (24 * 30))
		return fmt.Sprintf("%dmon", months)
	}

	// 超过一年，显示年数
	years := int(duration.Hours() / (24 * 365))
	return fmt.Sprintf("%dy", years)
}

// ListNamespaces 获取Namespace列表
func (s *NamespaceService) ListNamespaces() ([]model.NamespaceInfo, error) {
	namespaceList, err := s.k8sClient.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var namespaces []model.NamespaceInfo
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, model.NamespaceInfo{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
			Age:    formatAge(ns.CreationTimestamp.Time),
		})
	}

	return namespaces, nil
}

// CreateNamespace 创建Namespace
func (s *NamespaceService) CreateNamespace(name string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := s.k8sClient.ClientSet.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

// DeleteNamespace 删除Namespace
func (s *NamespaceService) DeleteNamespace(name string) error {
	return s.k8sClient.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
