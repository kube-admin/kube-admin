package service

import (
	"context"

	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// nodeMetricsMap 返回 nodeName -> 资源使用量映射。
// metrics-server 不可用时返回 nil，调用方需优雅降级。
func nodeMetricsMap(client *k8s.Client) map[string]corev1.ResourceList {
	if client == nil || client.MetricsClientSet == nil {
		return nil
	}
	list, err := client.MetricsClientSet.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil
	}
	m := make(map[string]corev1.ResourceList, len(list.Items))
	for i := range list.Items {
		m[list.Items[i].Name] = list.Items[i].Usage
	}
	return m
}

// podMetricsMap 返回 podName -> (containerName -> 资源使用量) 映射。
func podMetricsMap(client *k8s.Client, namespace string) map[string]map[string]corev1.ResourceList {
	if client == nil || client.MetricsClientSet == nil {
		return nil
	}
	list, err := client.MetricsClientSet.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil
	}
	m := make(map[string]map[string]corev1.ResourceList, len(list.Items))
	for i := range list.Items {
		containers := make(map[string]corev1.ResourceList, len(list.Items[i].Containers))
		for _, c := range list.Items[i].Containers {
			containers[c.Name] = c.Usage
		}
		m[list.Items[i].Name] = containers
	}
	return m
}

// calcPercent 计算使用率百分比。CPU 与内存均用 MilliValue 取比值，
// 单位一致保证比值正确，且 int64 范围足够容纳内存字节量级。
func calcPercent(used, total resource.Quantity) float64 {
	if total.IsZero() {
		return 0
	}
	return float64(used.MilliValue()) / float64(total.MilliValue()) * 100
}
