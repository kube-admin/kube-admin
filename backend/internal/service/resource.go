package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

// ResourceService 通用资源服务，基于 dynamic client + RESTMapper，
// 支持任意 K8s 资源的列表/获取/删除/apply/patch，避免为每种资源重复实现。
type ResourceService struct {
	k8sClient *k8s.Client
}

// NewResourceService 创建通用资源服务
func NewResourceService(c *k8s.Client) *ResourceService {
	return &ResourceService{k8sClient: c}
}

func (s *ResourceService) dynamicInterface() (dynamic.Interface, error) {
	return dynamic.NewForConfig(s.k8sClient.Config)
}

// gvrFor 通过 RESTMapper 将 GVK 映射为 GVR，并返回是否命名空间级资源
func (s *ResourceService) gvrFor(gvk schema.GroupVersionKind) (schema.GroupVersionResource, bool, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(s.k8sClient.Config)
	if err != nil {
		return schema.GroupVersionResource{}, false, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
	mapping, err := mapper.RESTMapping(gk, gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, false, err
	}
	return mapping.Resource, mapping.Scope.Name() == meta.RESTScopeNameNamespace, nil
}

// namespacedResource 根据 GVR 构造 dynamic 资源接口，namespace 为空或 "all" 时为集群级/所有命名空间
func (s *ResourceService) namespacedResource(gvr schema.GroupVersionResource, namespace string) (dynamic.ResourceInterface, error) {
	dyn, err := s.dynamicInterface()
	if err != nil {
		return nil, err
	}
	// "all" 是前端「所有命名空间」哨兵，视为空（k8s all-namespaces）
	if namespace != "" && namespace != "all" {
		return dyn.Resource(gvr).Namespace(namespace), nil
	}
	return dyn.Resource(gvr), nil
}

// List 通用列表
func (s *ResourceService) List(gvr schema.GroupVersionResource, namespace string) (*unstructured.UnstructuredList, error) {
	iface, err := s.namespacedResource(gvr, namespace)
	if err != nil {
		return nil, err
	}
	return iface.List(context.TODO(), metav1.ListOptions{})
}

// Get 通用获取
func (s *ResourceService) Get(gvr schema.GroupVersionResource, namespace, name string) (*unstructured.Unstructured, error) {
	iface, err := s.namespacedResource(gvr, namespace)
	if err != nil {
		return nil, err
	}
	return iface.Get(context.TODO(), name, metav1.GetOptions{})
}

// Delete 通用删除
func (s *ResourceService) Delete(gvr schema.GroupVersionResource, namespace, name string) error {
	iface, err := s.namespacedResource(gvr, namespace)
	if err != nil {
		return err
	}
	return iface.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// ApplyFromYAML 解析 YAML 并创建或更新资源（存在则更新，不存在则创建）
func (s *ResourceService) ApplyFromYAML(yamlStr string) (*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader([]byte(yamlStr)), 4096)
	var obj unstructured.Unstructured
	if err := decoder.Decode(&obj); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %w", err)
	}
	if len(obj.Object) == 0 {
		return nil, fmt.Errorf("YAML 解析后为空对象")
	}

	gvk := obj.GroupVersionKind()
	gvr, namespaced, err := s.gvrFor(gvk)
	if err != nil {
		return nil, fmt.Errorf("无法识别资源类型 %s: %w", gvk.String(), err)
	}

	dyn, err := s.dynamicInterface()
	if err != nil {
		return nil, err
	}
	var iface dynamic.ResourceInterface
	if namespaced {
		iface = dyn.Resource(gvr).Namespace(obj.GetNamespace())
	} else {
		iface = dyn.Resource(gvr)
	}

	// 存在则更新（保留 resourceVersion），不存在则创建
	existing, err := iface.Get(context.TODO(), obj.GetName(), metav1.GetOptions{})
	if err == nil {
		obj.SetResourceVersion(existing.GetResourceVersion())
		return iface.Update(context.TODO(), &obj, metav1.UpdateOptions{})
	}
	return iface.Create(context.TODO(), &obj, metav1.CreateOptions{})
}

// Patch 通用补丁
func (s *ResourceService) Patch(gvr schema.GroupVersionResource, namespace, name string, patchType types.PatchType, patchData []byte) (*unstructured.Unstructured, error) {
	iface, err := s.namespacedResource(gvr, namespace)
	if err != nil {
		return nil, err
	}
	return iface.Patch(context.TODO(), name, patchType, patchData, metav1.PatchOptions{})
}

// Scale 通用扩缩容（适用于含 spec.replicas 的 workload：Deployment/StatefulSet/DaemonSet/ReplicaSet）
func (s *ResourceService) Scale(gvr schema.GroupVersionResource, namespace, name string, replicas int32) error {
	iface, err := s.namespacedResource(gvr, namespace)
	if err != nil {
		return err
	}
	u, err := iface.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	spec, _ := u.Object["spec"].(map[string]interface{})
	if spec == nil {
		spec = map[string]interface{}{}
	}
	spec["replicas"] = replicas
	u.Object["spec"] = spec
	_, err = iface.Update(context.TODO(), u, metav1.UpdateOptions{})
	return err
}

// Restart 通用滚动重启（向 spec.template.metadata.annotations 注入 restartedAt 触发滚动更新）
func (s *ResourceService) Restart(gvr schema.GroupVersionResource, namespace, name string) error {
	patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, metav1.Now().Format("2006-01-02T15:04:05Z"))
	_, err := s.Patch(gvr, namespace, name, types.MergePatchType, []byte(patchData))
	return err
}
