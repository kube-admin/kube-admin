# 仪表盘与监控

仪表盘提供集群全局概览与实时资源使用率。

## 仪表盘

进入「Kubernetes → Dashboard」，可看到：

- **资源计数卡片**：节点 / 命名空间 / Pod / Deployment 数量
- **CPU 与内存使用率仪表盘**：全集群聚合的实时使用率（基于 metrics-server）
- **Pod 状态分布饼图**：Running / Pending / Failed / Succeeded / Unknown
- **资源统计**：Service / ConfigMap / Secret 数量

页面每 30 秒自动刷新。

## 节点监控

进入「Kubernetes → Nodes」，每个节点显示：

- **CPU 使用率**：实时百分比进度条 + 已用量 / 可分配
- **内存使用率**：实时百分比进度条 + 已用量 / 可分配
- 点击「查看详情」可看节点条件、标签、容量等。

## Pod 资源

Pod 详情中可查看每个容器的：
- 资源 request / limit
- 实时 CPU / 内存使用量（metrics-server）

## 关于 metrics-server

::: tip 必须安装 metrics-server
实时 CPU/内存使用率依赖集群已部署 [metrics-server](https://github.com/kubernetes-sigs/metrics-server)。若未安装，相关使用率会显示为空，并出现黄色提示条。
:::

在 minikube 中启用：

```bash
minikube addons enable metrics-server
```

在标准集群中部署：

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

## 下一步

- [核心资源](./workloads.md)
- [终端与日志](./terminal-logs.md)
