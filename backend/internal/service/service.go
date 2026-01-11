package service

import (
	"context"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceService Service服务
type ServiceService struct {
	k8sClient *k8s.Client
}

// NewServiceService 创建Service服务
func NewServiceService(k8sClient *k8s.Client) *ServiceService {
	return &ServiceService{k8sClient: k8sClient}
}

// GetK8sClient 获取K8s客户端
func (s *ServiceService) GetK8sClient() *k8s.Client {
	return s.k8sClient
}

// ListServices 获取Service列表
func (s *ServiceService) ListServices(namespace string) ([]model.ServiceInfo, error) {
	serviceList, err := s.k8sClient.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var services []model.ServiceInfo
	for _, svc := range serviceList.Items {
		services = append(services, s.convertService(&svc))
	}

	return services, nil
}

// GetService 获取Service详情
func (s *ServiceService) GetService(namespace, name string) (*model.ServiceInfo, error) {
	svc, err := s.k8sClient.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	serviceInfo := s.convertService(svc)
	return &serviceInfo, nil
}

// DeleteService 删除Service
func (s *ServiceService) DeleteService(namespace, name string) error {
	return s.k8sClient.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// convertService 转换Service对象
func (s *ServiceService) convertService(svc *corev1.Service) model.ServiceInfo {
	serviceInfo := model.ServiceInfo{
		K8sResource: model.K8sResource{
			Name:              svc.Name,
			Namespace:         svc.Namespace,
			Labels:            svc.Labels,
			Annotations:       svc.Annotations,
			CreationTimestamp: svc.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   svc.ResourceVersion,
		},
		Type:      string(svc.Spec.Type),
		ClusterIP: svc.Spec.ClusterIP,
		Selector:  svc.Spec.Selector,
	}

	// 外部IP
	if len(svc.Spec.ExternalIPs) > 0 {
		serviceInfo.ExternalIP = svc.Spec.ExternalIPs
	}

	// 端口信息
	for _, port := range svc.Spec.Ports {
		serviceInfo.Ports = append(serviceInfo.Ports, model.ServicePort{
			Name:       port.Name,
			Protocol:   string(port.Protocol),
			Port:       port.Port,
			TargetPort: port.TargetPort.String(),
			NodePort:   port.NodePort,
		})
	}

	return serviceInfo
}
