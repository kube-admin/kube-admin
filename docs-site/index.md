---
layout: home

hero:
  name: Kube Admin
  text: 多集群 Kubernetes 管理平台
  image:
    src: /hero.svg
    alt: Kube Admin 集群总览仪表盘
  tagline: Vue 3 + Go + client-go，资源管理、实时监控、Web 终端、YAML 工作台、RBAC 鉴权一站式
  actions:
    - theme: brand
      text: 快速上手
      link: /guide/getting-started
    - theme: alt
      text: 项目介绍
      link: /guide/overview

features:
  - title: 多集群管理
    details: 通过 kubeconfig / 文件路径 / ServerURL+Token 接入多集群，凭据 AES-256-GCM 加密存储
  - title: 实时监控
    details: 接入 metrics-server，节点与 Pod 的 CPU/内存实时使用率，Dashboard 可视化图表
  - title: 资源全覆盖
    details: 内置核心资源 CRUD + 通用资源浏览器（dynamic client）管理任意 K8s 资源
  - title: Web 终端与日志流
    details: Pod 交互式终端（xterm）+ WebSocket 实时日志（follow / previous / 搜索 / 下载）
  - title: YAML 工作台
    details: Monaco 编辑器在线编辑、校验、apply 任意资源，支持创建与更新
  - title: 安全合规
    details: JWT + bcrypt + RBAC 角色鉴权 + 操作审计日志，健康探针 + 优雅关闭
---
