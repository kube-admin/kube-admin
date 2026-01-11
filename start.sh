#!/bin/bash

echo "🚀 启动 Kubernetes 管理系统"
echo "================================"

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ 未找到 Go 环境，请先安装 Go 1.21+"
    exit 1
fi

# 检查 Node 环境
if ! command -v node &> /dev/null; then
    echo "❌ 未找到 Node.js 环境，请先安装 Node.js 16+"
    exit 1
fi

echo "✅ 环境检查通过"
echo ""

# 启动后端
echo "📦 启动后端服务..."
cd backend
go mod tidy > /dev/null 2>&1
go run cmd/main.go &
BACKEND_PID=$!
echo "✅ 后端服务已启动 (PID: $BACKEND_PID) - http://localhost:8080"
cd ..

# 等待后端启动
sleep 3

# 启动前端
echo "📦 启动前端服务..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "📥 安装前端依赖..."
    pnpm install
fi
pnpm dev &
FRONTEND_PID=$!
echo "✅ 前端服务已启动 (PID: $FRONTEND_PID) - http://localhost:3000"
cd ..

echo ""
echo "================================"
echo "🎉 系统启动成功!"
echo ""
echo "📝 访问地址:"
echo "   前端: http://localhost:3000"
echo "   后端: http://localhost:8080"
echo ""
echo "👤 默认登录信息:"
echo "   用户名: admin"
echo "   密码: admin123"
echo ""
echo "⚠️  按 Ctrl+C 停止服务"
echo "================================"

# 等待用户中断
trap "echo ''; echo '🛑 正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit 0" INT

wait
