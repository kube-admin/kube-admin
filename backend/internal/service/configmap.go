package service

import (
	"context"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapService ConfigMap服务
type ConfigMapService struct {
	k8sClient *k8s.Client
}

// NewConfigMapService 创建ConfigMap服务
func NewConfigMapService(k8sClient *k8s.Client) *ConfigMapService {
	return &ConfigMapService{k8sClient: k8sClient}
}

// ListConfigMaps 获取ConfigMap列表
func (s *ConfigMapService) ListConfigMaps(namespace string) ([]model.ConfigMapInfo, error) {
	configMapList, err := s.k8sClient.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var configMaps []model.ConfigMapInfo
	for _, cm := range configMapList.Items {
		configMaps = append(configMaps, s.convertConfigMap(&cm))
	}

	return configMaps, nil
}

// GetConfigMap 获取ConfigMap详情
func (s *ConfigMapService) GetConfigMap(namespace, name string) (*model.ConfigMapInfo, error) {
	cm, err := s.k8sClient.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	configMapInfo := s.convertConfigMap(cm)
	return &configMapInfo, nil
}

// CreateConfigMap 创建ConfigMap
func (s *ConfigMapService) CreateConfigMap(namespace, name string, data map[string]string) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}

	_, err := s.k8sClient.ClientSet.CoreV1().ConfigMaps(namespace).Create(context.TODO(), cm, metav1.CreateOptions{})
	return err
}

// UpdateConfigMap 更新ConfigMap
func (s *ConfigMapService) UpdateConfigMap(namespace, name string, data map[string]string) error {
	cm, err := s.k8sClient.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	cm.Data = data
	_, err = s.k8sClient.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cm, metav1.UpdateOptions{})
	return err
}

// DeleteConfigMap 删除ConfigMap
func (s *ConfigMapService) DeleteConfigMap(namespace, name string) error {
	return s.k8sClient.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// convertConfigMap 转换ConfigMap对象
func (s *ConfigMapService) convertConfigMap(cm *corev1.ConfigMap) model.ConfigMapInfo {
	return model.ConfigMapInfo{
		K8sResource: model.K8sResource{
			Name:              cm.Name,
			Namespace:         cm.Namespace,
			Labels:            cm.Labels,
			Annotations:       cm.Annotations,
			CreationTimestamp: cm.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   cm.ResourceVersion,
		},
		Data: cm.Data,
	}
}
