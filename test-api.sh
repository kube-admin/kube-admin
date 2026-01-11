#!/bin/bash

# Kubernetes 管理系统 API 测试脚本

API_HOST="http://localhost:8080"
TOKEN=""

echo "🧪 Kubernetes 管理系统 API 测试"
echo "================================"

# 1. 测试登录
echo ""
echo "1️⃣  测试登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_HOST/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

echo "$LOGIN_RESPONSE" | jq '.'

# 提取 Token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ 登录失败，无法获取 Token"
    exit 1
fi

echo "✅ 登录成功，Token: ${TOKEN:0:20}..."

# 2. 测试 Dashboard 统计
echo ""
echo "2️⃣  测试 Dashboard 统计..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/dashboard/stats" | jq '.'

# 3. 测试 Namespace 列表
echo ""
echo "3️⃣  测试 Namespace 列表..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/namespaces" | jq '.data[] | {name, status, age}'

# 4. 测试 Node 列表
echo ""
echo "4️⃣  测试 Node 列表..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/nodes" | jq '.data[] | {name, status, internal_ip}'

# 5. 测试 Pod 列表
echo ""
echo "5️⃣  测试 Pod 列表 (default namespace)..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/pods?namespace=default" | jq '.data[] | {name, status, pod_ip, node_name}'

# 6. 测试 Deployment 列表
echo ""
echo "6️⃣  测试 Deployment 列表 (default namespace)..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/deployments?namespace=default" | jq '.data[] | {name, replicas, ready_replicas, available_replicas}'

# 7. 测试 ConfigMap 列表
echo ""
echo "7️⃣  测试 ConfigMap 列表 (default namespace)..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/configmaps?namespace=default" | jq '.data[] | {name, namespace}'

# 8. 测试 Secret 列表
echo ""
echo "8️⃣  测试 Secret 列表 (default namespace)..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/secrets?namespace=default" | jq '.data[] | {name, type}'

# 9. 测试创建 ConfigMap
echo ""
echo "9️⃣  测试创建 ConfigMap..."
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "namespace": "default",
    "name": "test-config",
    "data": {
      "test.key": "test.value",
      "app.name": "kube-admin-test"
    }
  }' \
  "$API_HOST/api/v1/configmaps" | jq '.'

# 10. 测试获取 ConfigMap 详情
echo ""
echo "🔟 测试获取 ConfigMap 详情..."
curl -s -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/configmaps/test-config?namespace=default" | jq '.data'

# 11. 测试删除 ConfigMap
echo ""
echo "1️⃣1️⃣  测试删除 ConfigMap..."
curl -s -X DELETE -H "Authorization: Bearer $TOKEN" \
  "$API_HOST/api/v1/configmaps/test-config?namespace=default" | jq '.'

echo ""
echo "================================"
echo "✅ API 测试完成！"
echo ""
echo "💡 提示:"
echo "  - 确保后端服务已启动 (http://localhost:8080)"
echo "  - 确保有可访问的 Kubernetes 集群"
echo "  - 可以修改脚本测试其他 API 接口"
