package model

// K8sResource K8s资源通用模型
type K8sResource struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	CreationTimestamp string            `json:"creation_timestamp"`
	ResourceVersion   string            `json:"resource_version"`
}

// PodInfo Pod信息
type PodInfo struct {
	K8sResource
	Status     string            `json:"status"`
	PodIP      string            `json:"pod_ip"`
	NodeName   string            `json:"node_name"`
	Containers []ContainerInfo   `json:"containers"`
	Conditions []PodCondition    `json:"conditions"`
}

// ContainerInfo 容器信息
type ContainerInfo struct {
	Name         string        `json:"name"`
	Image        string        `json:"image"`
	Ready        bool          `json:"ready"`
	RestartCount int32         `json:"restart_count"`
	State        string        `json:"state"`
	Resources    ResourceUsage `json:"resources"`
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPURequest    string `json:"cpu_request"`
	MemoryRequest string `json:"memory_request"`
	CPULimit      string `json:"cpu_limit"`
	MemoryLimit   string `json:"memory_limit"`
}

// PodCondition Pod状态条件
type PodCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// DeploymentInfo Deployment信息
type DeploymentInfo struct {
	K8sResource
	Replicas          int32  `json:"replicas"`
	ReadyReplicas     int32  `json:"ready_replicas"`
	UpdatedReplicas   int32  `json:"updated_replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
	Strategy          string `json:"strategy"`
}

// ServiceInfo Service信息
type ServiceInfo struct {
	K8sResource
	Type       string              `json:"type"`
	ClusterIP  string              `json:"cluster_ip"`
	ExternalIP []string            `json:"external_ip,omitempty"`
	Ports      []ServicePort       `json:"ports"`
	Selector   map[string]string   `json:"selector"`
}

// ServicePort Service端口
type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort string `json:"target_port"`
	NodePort   int32  `json:"node_port,omitempty"`
}

// NodeInfo Node信息
type NodeInfo struct {
	K8sResource
	Status         string            `json:"status"`
	InternalIP     string            `json:"internal_ip"`
	OSImage        string            `json:"os_image"`
	KubeletVersion string            `json:"kubelet_version"`
	ContainerRuntime string          `json:"container_runtime"`
	Capacity       ResourceCapacity  `json:"capacity"`
	Allocatable    ResourceCapacity  `json:"allocatable"`
	Conditions     []NodeCondition   `json:"conditions"`
}

// ResourceCapacity 资源容量
type ResourceCapacity struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	Pods   string `json:"pods"`
}

// NodeCondition Node状态条件
type NodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// NamespaceInfo Namespace信息
type NamespaceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Age    string `json:"age"`
}

// ConfigMapInfo ConfigMap信息
type ConfigMapInfo struct {
	K8sResource
	Data map[string]string `json:"data"`
}

// SecretInfo Secret信息
type SecretInfo struct {
	K8sResource
	Type string            `json:"type"`
	Data map[string]string `json:"data"` // base64 encoded
}
