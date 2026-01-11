package service

import (
	"context"

	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeploymentService Deployment服务
type DeploymentService struct {
	k8sClient *k8s.Client
}

// NewDeploymentService 创建Deployment服务
func NewDeploymentService(k8sClient *k8s.Client) *DeploymentService {
	return &DeploymentService{k8sClient: k8sClient}
}

// GetK8sClient 获取K8s客户端
func (s *DeploymentService) GetK8sClient() *k8s.Client {
	return s.k8sClient
}

// ListDeployments 获取Deployment列表
func (s *DeploymentService) ListDeployments(namespace string) ([]model.DeploymentInfo, error) {
	deploymentList, err := s.k8sClient.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deployments []model.DeploymentInfo
	for _, deploy := range deploymentList.Items {
		deployments = append(deployments, s.convertDeployment(&deploy))
	}

	return deployments, nil
}

// GetDeployment 获取Deployment详情
func (s *DeploymentService) GetDeployment(namespace, name string) (*model.DeploymentInfo, error) {
	deploy, err := s.k8sClient.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	deployInfo := s.convertDeployment(deploy)
	return &deployInfo, nil
}

// DeleteDeployment 删除Deployment
func (s *DeploymentService) DeleteDeployment(namespace, name string) error {
	return s.k8sClient.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// ScaleDeployment 扩缩容Deployment
func (s *DeploymentService) ScaleDeployment(namespace, name string, replicas int32) error {
	deploy, err := s.k8sClient.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	deploy.Spec.Replicas = &replicas
	_, err = s.k8sClient.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	return err
}

// RestartDeployment 重启Deployment
func (s *DeploymentService) RestartDeployment(namespace, name string) error {
	deploy, err := s.k8sClient.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deploy.Spec.Template.Annotations == nil {
		deploy.Spec.Template.Annotations = make(map[string]string)
	}
	deploy.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = metav1.Now().Format("2006-01-02T15:04:05Z")

	_, err = s.k8sClient.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	return err
}

// convertDeployment 转换Deployment对象
func (s *DeploymentService) convertDeployment(deploy *appsv1.Deployment) model.DeploymentInfo {
	replicas := int32(0)
	if deploy.Spec.Replicas != nil {
		replicas = *deploy.Spec.Replicas
	}

	return model.DeploymentInfo{
		K8sResource: model.K8sResource{
			Name:              deploy.Name,
			Namespace:         deploy.Namespace,
			Labels:            deploy.Labels,
			Annotations:       deploy.Annotations,
			CreationTimestamp: deploy.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   deploy.ResourceVersion,
		},
		Replicas:          replicas,
		ReadyReplicas:     deploy.Status.ReadyReplicas,
		UpdatedReplicas:   deploy.Status.UpdatedReplicas,
		AvailableReplicas: deploy.Status.AvailableReplicas,
		Strategy:          string(deploy.Spec.Strategy.Type),
	}
}
