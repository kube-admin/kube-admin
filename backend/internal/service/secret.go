package service

import (
	"context"
	"encoding/base64"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SecretService Secret服务
type SecretService struct {
	k8sClient *k8s.Client
}

// NewSecretService 创建Secret服务
func NewSecretService(k8sClient *k8s.Client) *SecretService {
	return &SecretService{k8sClient: k8sClient}
}

// ListSecrets 获取Secret列表
func (s *SecretService) ListSecrets(namespace string) ([]model.SecretInfo, error) {
	secretList, err := s.k8sClient.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var secrets []model.SecretInfo
	for _, secret := range secretList.Items {
		secrets = append(secrets, s.convertSecret(&secret, false))
	}

	return secrets, nil
}

// GetSecret 获取Secret详情
func (s *SecretService) GetSecret(namespace, name string, decode bool) (*model.SecretInfo, error) {
	secret, err := s.k8sClient.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	secretInfo := s.convertSecret(secret, decode)
	return &secretInfo, nil
}

// CreateSecret 创建Secret
func (s *SecretService) CreateSecret(namespace, name, secretType string, data map[string]string) error {
	// 将字符串数据转换为字节数组
	byteData := make(map[string][]byte)
	for k, v := range data {
		byteData[k] = []byte(v)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: corev1.SecretType(secretType),
		Data: byteData,
	}

	_, err := s.k8sClient.ClientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	return err
}

// UpdateSecret 更新Secret
func (s *SecretService) UpdateSecret(namespace, name string, data map[string]string) error {
	secret, err := s.k8sClient.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 将字符串数据转换为字节数组
	byteData := make(map[string][]byte)
	for k, v := range data {
		byteData[k] = []byte(v)
	}

	secret.Data = byteData
	_, err = s.k8sClient.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	return err
}

// DeleteSecret 删除Secret
func (s *SecretService) DeleteSecret(namespace, name string) error {
	return s.k8sClient.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// convertSecret 转换Secret对象
func (s *SecretService) convertSecret(secret *corev1.Secret, decode bool) model.SecretInfo {
	data := make(map[string]string)

	if decode {
		// 解码显示实际值
		for k, v := range secret.Data {
			data[k] = string(v)
		}
	} else {
		// 不解码,显示base64编码后的值
		for k, v := range secret.Data {
			data[k] = base64.StdEncoding.EncodeToString(v)
		}
	}

	return model.SecretInfo{
		K8sResource: model.K8sResource{
			Name:              secret.Name,
			Namespace:         secret.Namespace,
			Labels:            secret.Labels,
			Annotations:       secret.Annotations,
			CreationTimestamp: secret.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   secret.ResourceVersion,
		},
		Type: string(secret.Type),
		Data: data,
	}
}
