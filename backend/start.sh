#!/bin/bash

# Kubernetes 管理系统后端启动脚本

echo "🚀 启动 Kubernetes 管理系统后端服务..."
echo ""

# 进入后端目录
cd "$(dirname "$0")"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到 Go 环境，请先安装 Go"
    exit 1
fi

# 检查依赖
echo "📦 检查依赖..."
go mod tidy

# 设置环境变量（可选）
export PORT="${PORT:-8080}"
export JWT_SECRET="${JWT_SECRET:-your-secret-key-change-in-production}"
export KUBECONFIG="${KUBECONFIG:-$HOME/.kube/config}"

echo "⚙️  配置信息:"
echo "  - 端口: $PORT"
echo "  - Kubeconfig: $KUBECONFIG"
echo ""

# 运行服务
echo "🎯 启动服务..."
go run cmd/main.go
