package service

import (
	"context"
	"encoding/json"
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
		info := model.NamespaceInfo{
			Name:   ns.Name,
			Status: string(ns.Status.Phase),
			Age:    formatAge(ns.CreationTimestamp.Time),
		}
		// 转换 FinalizerName 切片为 string 切片
		for _, f := range ns.Spec.Finalizers {
			info.Finalizers = append(info.Finalizers, string(f))
		}
		if ns.DeletionTimestamp != nil {
			t := ns.DeletionTimestamp.Time
			info.DeletionTimestamp = &t
		}
		for _, c := range ns.Status.Conditions {
			info.Conditions = append(info.Conditions, model.Condition{
				Type:               string(c.Type),
				Status:             string(c.Status),
				Reason:             c.Reason,
				Message:            c.Message,
				LastTransitionTime: c.LastTransitionTime.Time,
			})
		}
		namespaces = append(namespaces, info)
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

// FinalizeNamespace 强制清理 namespace 的 finalizers，解除 Terminating 卡死。
//
// 场景：namespace 因残留 finalizer（controller 已不存在等）卡在 Terminating，
// 普通删除无效。通过 /finalize 子资源提交 finalizers=[] 的 namespace 对象，
// 绕过等待 controller，使 namespace 立即被真正删除。
// 危险操作：会立即真正删除该 namespace 及其下所有资源，调用方必须二次确认。
func (s *NamespaceService) FinalizeNamespace(name string) error {
	ns, err := s.k8sClient.ClientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	ns.Finalizers = []string{}
	body, err := json.Marshal(ns)
	if err != nil {
		return err
	}
	_, err = s.k8sClient.ClientSet.CoreV1().RESTClient().Put().
		Resource("namespaces").Name(name).SubResource("finalize").
		Body(body).DoRaw(context.TODO())
	return err
}

// ListUnavailableAPIServices 列出集群中 Available 状态非 True 的 APIService。
// 用于诊断 namespace 因 DiscoveryFailed（失效 APIService）卡在 Terminating 的场景。
func (s *NamespaceService) ListUnavailableAPIServices() ([]model.APIServiceInfo, error) {
	list, err := s.k8sClient.AggregatorClient.ApiregistrationV1().APIServices().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []model.APIServiceInfo
	for _, a := range list.Items {
		available := false
		found := false
		var reason, msg string
		for _, c := range a.Status.Conditions {
			if c.Type == "Available" {
				found = true
				available = string(c.Status) == "True"
				reason = string(c.Reason)
				msg = c.Message
				break
			}
		}
		if !found || available {
			continue
		}
		svc := ""
		if a.Spec.Service != nil {
			svc = a.Spec.Service.Namespace + "/" + a.Spec.Service.Name
		}
		result = append(result, model.APIServiceInfo{
			Name:    a.Name,
			Service: svc,
			Status:  reason,
			Message: msg,
			Age:     formatAge(a.CreationTimestamp.Time),
		})
	}
	return result, nil
}

// DeleteAPIService 删除指定 APIService（集群级危险操作，调用方需二次确认）。
// 用于移除后端服务已失效的 APIService，使被卡住的 namespace 继续完成删除。
func (s *NamespaceService) DeleteAPIService(name string) error {
	return s.k8sClient.AggregatorClient.ApiregistrationV1().APIServices().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
