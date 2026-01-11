#!/bin/bash

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "Docker 未安装，请先安装Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose 未安装，请先安装Docker Compose"
    exit 1
fi

# 启动服务
echo "正在启动Kubernetes管理系统..."
docker-compose up -d

# 检查服务状态
echo "正在检查服务状态..."
docker-compose ps

# 输出访问地址
echo ""
echo "服务已启动，访问地址：http://localhost"
echo "后端API地址：http://localhost:8080"
echo ""
echo "使用以下命令查看日志："
echo "docker-compose logs -f frontend"
echo "docker-compose logs -f backend"
echo ""
echo "使用以下命令停止服务："
echo "docker-compose down"
echo ""
echo "使用以下命令重启服务："
echo "docker-compose restart"
echo ""
