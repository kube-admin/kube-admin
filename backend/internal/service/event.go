package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventService K8s 事件服务
type EventService struct {
	k8sClient *k8s.Client
}

// NewEventService 创建事件服务
func NewEventService(k8sClient *k8s.Client) *EventService {
	return &EventService{k8sClient: k8sClient}
}

// ListEvents 查询事件。fieldSelector 可按 involvedObject.kind/name 过滤。
func (s *EventService) ListEvents(namespace, fieldSelector string) ([]model.EventInfo, error) {
	list, err := s.k8sClient.ClientSet.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fieldSelector,
	})
	if err != nil {
		return nil, err
	}

	events := make([]model.EventInfo, 0, len(list.Items))
	for i := range list.Items {
		events = append(events, convertEvent(&list.Items[i]))
	}
	return events, nil
}

// convertEvent 转换 Event 对象
func convertEvent(e *corev1.Event) model.EventInfo {
	info := model.EventInfo{
		K8sResource: model.K8sResource{
			Name:              e.Name,
			Namespace:         e.Namespace,
			CreationTimestamp: e.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   e.ResourceVersion,
		},
		Type:           e.Type,
		Reason:         e.Reason,
		Message:        e.Message,
		Count:          e.Count,
		FirstTimestamp: formatMetaTime(e.FirstTimestamp),
		LastTimestamp:  formatMetaTime(e.LastTimestamp),
	}

	if e.InvolvedObject.Kind != "" {
		info.InvolvedObject = fmt.Sprintf("%s/%s", e.InvolvedObject.Kind, e.InvolvedObject.Name)
	}

	// Source 含 component/host，按需拼接
	parts := make([]string, 0, 2)
	if e.Source.Component != "" {
		parts = append(parts, e.Source.Component)
	}
	if e.Source.Host != "" {
		parts = append(parts, e.Source.Host)
	}
	info.Source = strings.Join(parts, "/")

	return info
}

// formatMetaTime 格式化 metav1.Time，零值返回空串
func formatMetaTime(t metav1.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
