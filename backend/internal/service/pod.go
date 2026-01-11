package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/kube-admin/kube-admin/backend/internal/model"
	"github.com/kube-admin/kube-admin/backend/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
)

// PodService Pod服务
type PodService struct {
	k8sClient *k8s.Client
}

// NewPodService 创建Pod服务
func NewPodService(k8sClient *k8s.Client) *PodService {
	return &PodService{k8sClient: k8sClient}
}

// GetK8sClient 获取K8s客户端
func (s *PodService) GetK8sClient() *k8s.Client {
	return s.k8sClient
}

// ListPods 获取Pod列表
func (s *PodService) ListPods(namespace string) ([]model.PodInfo, error) {
	podList, err := s.k8sClient.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var pods []model.PodInfo
	for _, pod := range podList.Items {
		pods = append(pods, s.convertPod(&pod))
	}

	return pods, nil
}

// GetPod 获取Pod详情
func (s *PodService) GetPod(namespace, name string) (*model.PodInfo, error) {
	pod, err := s.k8sClient.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podInfo := s.convertPod(pod)
	return &podInfo, nil
}

// DeletePod 删除Pod
func (s *PodService) DeletePod(namespace, name string) error {
	return s.k8sClient.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// GetPodLogs 获取Pod日志
func (s *PodService) GetPodLogs(namespace, name, container string, tailLines int64) (string, error) {
	req := s.k8sClient.ClientSet.CoreV1().Pods(namespace).GetLogs(name, &corev1.PodLogOptions{
		Container: container,
		TailLines: &tailLines,
	})

	logs, err := req.DoRaw(context.TODO())
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

// ExecCommand 在Pod中执行命令
func (s *PodService) ExecCommand(namespace, podName, containerName string, command []string) (string, string, error) {
	// 创建REST client
	client := s.k8sClient.ClientSet.CoreV1().RESTClient()

	// 创建exec请求
	req := client.Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName).
		Param("stdout", "true").
		Param("stderr", "true")

	// 添加命令参数
	for _, cmd := range command {
		req = req.Param("command", cmd)
	}

	// 创建executor
	executor, err := remotecommand.NewSPDYExecutor(s.k8sClient.Config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("创建执行器失败: %v", err)
	}

	// 执行命令
	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("执行命令失败: %v", err)
	}

	return stdout.String(), stderr.String(), nil
}

// ExecTerminal 在Pod中创建交互式终端
func (s *PodService) ExecTerminal(namespace, podName, containerName string, wsConn *websocket.Conn) error {
	// 创建REST client
	client := s.k8sClient.ClientSet.CoreV1().RESTClient()

	// 创建exec请求
	req := client.Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("command", "/bin/sh")

	// 创建executor
	executor, err := remotecommand.NewSPDYExecutor(s.k8sClient.Config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("创建执行器失败: %v", err)
	}

	// 创建WebSocket到TTY的桥接器
	ptyHandler := &wsStreamHandler{
		conn:     wsConn,
		sizeChan: make(chan remotecommand.TerminalSize),
	}

	// 执行命令
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})

	if err != nil {
		return fmt.Errorf("执行终端失败: %v", err)
	}

	return nil
}

// wsStreamHandler WebSocket流处理程序
type wsStreamHandler struct {
	conn     *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
}

// Read 从WebSocket读取数据
func (w *wsStreamHandler) Read(p []byte) (int, error) {
	_, message, err := w.conn.ReadMessage()
	if err != nil {
		return copy(p, "\x04"), err // 发送EOF字符
	}
	return copy(p, message), nil
}

// Write 向WebSocket写入数据
func (w *wsStreamHandler) Write(p []byte) (int, error) {
	err := w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Next 返回下一个终端大小
func (w *wsStreamHandler) Next() *remotecommand.TerminalSize {
	select {
	case size := <-w.sizeChan:
		return &size
	default:
		return nil
	}
}

// Close 关闭连接
func (w *wsStreamHandler) Close() error {
	return w.conn.Close()
}

// convertPod 转换Pod对象
func (s *PodService) convertPod(pod *corev1.Pod) model.PodInfo {
	podInfo := model.PodInfo{
		K8sResource: model.K8sResource{
			Name:              pod.Name,
			Namespace:         pod.Namespace,
			Labels:            pod.Labels,
			Annotations:       pod.Annotations,
			CreationTimestamp: pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
			ResourceVersion:   pod.ResourceVersion,
		},
		Status:   string(pod.Status.Phase),
		PodIP:    pod.Status.PodIP,
		NodeName: pod.Spec.NodeName,
	}

	// 转换容器信息
	for _, container := range pod.Spec.Containers {
		var containerStatus *corev1.ContainerStatus
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Name == container.Name {
				containerStatus = &cs
				break
			}
		}

		containerInfo := model.ContainerInfo{
			Name:  container.Name,
			Image: container.Image,
		}

		if containerStatus != nil {
			containerInfo.Ready = containerStatus.Ready
			containerInfo.RestartCount = containerStatus.RestartCount

			if containerStatus.State.Running != nil {
				containerInfo.State = "Running"
			} else if containerStatus.State.Waiting != nil {
				containerInfo.State = fmt.Sprintf("Waiting: %s", containerStatus.State.Waiting.Reason)
			} else if containerStatus.State.Terminated != nil {
				containerInfo.State = fmt.Sprintf("Terminated: %s", containerStatus.State.Terminated.Reason)
			}
		}

		// 资源请求和限制
		if container.Resources.Requests != nil {
			if cpu, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				containerInfo.Resources.CPURequest = cpu.String()
			}
			if memory, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				containerInfo.Resources.MemoryRequest = memory.String()
			}
		}
		if container.Resources.Limits != nil {
			if cpu, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				containerInfo.Resources.CPULimit = cpu.String()
			}
			if memory, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				containerInfo.Resources.MemoryLimit = memory.String()
			}
		}

		podInfo.Containers = append(podInfo.Containers, containerInfo)
	}

	// 转换条件信息
	for _, condition := range pod.Status.Conditions {
		podInfo.Conditions = append(podInfo.Conditions, model.PodCondition{
			Type:    string(condition.Type),
			Status:  string(condition.Status),
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}

	return podInfo
}
